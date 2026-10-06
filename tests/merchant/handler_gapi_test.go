package merchant_test

import (
	"context"
	"testing"

	merchant_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"

	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MerchantGapiTestSuite struct {
	tests.BaseTestSuite
	commandClient pbmerchants.MerchantCommandServiceClient
	queryClient   pbmerchants.MerchantQueryServiceClient
	userID        int
	merchantID    int
}

func (s *MerchantGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)
	repos := repository.NewRepositories(mongoDB, pbusers.NewUserQueryServiceClient(s.Conns["user"]))

	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mencache := merchant_cache.NewMencache(cacheStore)

	svc := service.NewService(&service.Deps{
		Kafka:         nil,
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencache,
		Observability: s.Obs,
	})

	merchantHandler := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbmerchants.RegisterMerchantCommandServiceServer(server, merchantHandler.MerchantCommandHandler)
	pbmerchants.RegisterMerchantQueryServiceServer(server, merchantHandler.MerchantQuery)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.commandClient = pbmerchants.NewMerchantCommandServiceClient(conn)
	s.queryClient = pbmerchants.NewMerchantQueryServiceClient(conn)

	// 1. Seed dependencies
	s.userID = s.SeedUser(context.Background())
}

func (s *MerchantGapiTestSuite) TestMerchantGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pbmerchants.CreateMerchantRequest{
		UserId:       int32(s.userID),
		Name:         "Gapi Merchant",
		Description:  "Detailed description of the merchant.",
		Address:      "Merchant Street No. 1",
		ContactEmail: "gapi.merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}
	res, err := s.commandClient.Create(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Name, res.Data.Name)
	merchantID := res.Data.Id

	// 2. FindById
	found, err := s.queryClient.FindById(ctx, &pbmerchants.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)
	s.Equal(merchantID, found.Data.Id)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbmerchants.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbmerchants.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateReq := &pbmerchants.UpdateMerchantRequest{
		MerchantId:   merchantID,
		UserId:       int32(s.userID),
		Name:         "Gapi Merchant Updated",
		Description:  "Updated description.",
		Address:      "New Street 2",
		ContactEmail: "updated@example.com",
		ContactPhone: "08987654321",
		Status:       "waiting",
	}
	updateRes, err := s.commandClient.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updateRes.Data.Name)

	// 6. Update Status
	statusRes, err := s.commandClient.UpdateStatus(ctx, &pbmerchants.UpdateMerchantStatusRequest{
		MerchantId: merchantID,
		Status:     "active",
	})
	s.NoError(err)
	s.Equal("active", statusRes.Data.Status)

	// 7. Trash
	_, err = s.commandClient.TrashedMerchant(ctx, &pbmerchants.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbmerchants.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreMerchant(ctx, &pbmerchants.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedMerchant(ctx, &pbmerchants.FindByIdMerchantRequest{Id: merchantID})
	_, err = s.commandClient.DeleteMerchantPermanent(ctx, &pbmerchants.FindByIdMerchantRequest{Id: merchantID})
	s.NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllMerchant(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestMerchantGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantGapiTestSuite))
}
