package merchant_business_test

import (
	"context"
	"testing"

	biz_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/cache"
	biz_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/handler"
	biz_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	biz_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_business "github.com/MamangRust/monolith-ecommerce-pb/merchants_business"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MerchantBusinessGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbmerchants_business.MerchantBusinessQueryServiceClient
	commandClient pbmerchants_business.MerchantBusinessCommandServiceClient
}

func (s *MerchantBusinessGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Business dependencies
	mencache := biz_cache.NewMencache(cacheStore)
	repos := biz_repo.NewRepositories(
		mongoDB,
		pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	svc := biz_service.NewService(&biz_service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := biz_handler.NewHandler(&biz_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbmerchants_business.RegisterMerchantBusinessQueryServiceServer(server, handler.MerchantBusinessQuery)
	pbmerchants_business.RegisterMerchantBusinessCommandServiceServer(server, handler.MerchantBusinessCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbmerchants_business.NewMerchantBusinessQueryServiceClient(conn)
	s.commandClient = pbmerchants_business.NewMerchantBusinessCommandServiceClient(conn)
}

func (s *MerchantBusinessGapiTestSuite) TestMerchantBusinessGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := int32(s.SeedMerchant(ctx, userID))

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pbmerchants_business.CreateMerchantBusinessRequest{
		MerchantId:        merchID,
		BusinessType:      "GAPI Corp",
		TaxId:             "GAPI-TAX",
		EstablishedYear:   2020,
		NumberOfEmployees: 50,
		WebsiteUrl:        "http://gapi.corp",
	})
	s.NoError(err)
	s.NotNil(createRes)
	bizID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbmerchants_business.FindByIdMerchantBusinessRequest{Id: bizID})
	s.NoError(err)
	s.Equal("GAPI Corp", getRes.Data.BusinessType)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbmerchants.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbmerchants.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pbmerchants_business.UpdateMerchantBusinessRequest{
		MerchantBusinessInfoId: bizID,
		BusinessType:           "GAPI Corp Updated",
		TaxId:                  "GAPI-TAX-UPDATED",
		EstablishedYear:        2021,
		NumberOfEmployees:      60,
		WebsiteUrl:             "http://gapi.corp/updated",
	})
	s.NoError(err)
	s.Equal("GAPI Corp Updated", updateRes.Data.BusinessType)

	// 7. Trash
	_, err = s.commandClient.TrashedMerchantBusiness(ctx, &pbmerchants_business.FindByIdMerchantBusinessRequest{Id: bizID})
	s.NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbmerchants.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreMerchantBusiness(ctx, &pbmerchants_business.FindByIdMerchantBusinessRequest{Id: bizID})
	s.NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedMerchantBusiness(ctx, &pbmerchants_business.FindByIdMerchantBusinessRequest{Id: bizID})
	_, err = s.commandClient.DeleteMerchantBusinessPermanent(ctx, &pbmerchants_business.FindByIdMerchantBusinessRequest{Id: bizID})
	s.NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllMerchantBusiness(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllMerchantBusinessPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestMerchantBusinessGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantBusinessGapiTestSuite))
}
