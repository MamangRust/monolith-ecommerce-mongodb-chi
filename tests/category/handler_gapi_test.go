package category_test

import (
	"context"
	"testing"

	cat_cache "github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	cat_handler "github.com/MamangRust/monolith-ecommerce-grpc-category/handler"
	cat_repo "github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	cat_service "github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient           pbcategories.CategoryQueryServiceClient
	commandClient         pbcategories.CategoryCommandServiceClient
	statsClient           pbcategories.CategoryStatsServiceClient
	statsByIdClient       pbcategories.CategoryStatsByIdServiceClient
	statsByMerchantClient pbcategories.CategoryStatsByMerchantServiceClient
}

func (s *CategoryGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Category dependencies
	mencache := cat_cache.NewMencache(cacheStore)
	repos := cat_repo.NewRepositories(mongoDB)
	svc := cat_service.NewService(&cat_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := cat_handler.NewHandler(&cat_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// GRPC Server
	server := grpc.NewServer()
	pbcategories.RegisterCategoryQueryServiceServer(server, handler.CategoryQuery)
	pbcategories.RegisterCategoryCommandServiceServer(server, handler.CategoryCommand)
	pbcategories.RegisterCategoryStatsServiceServer(server, handler.CategoryStats)
	pbcategories.RegisterCategoryStatsByIdServiceServer(server, handler.CategoryStatsById)
	pbcategories.RegisterCategoryStatsByMerchantServiceServer(server, handler.CategoryStatsByMerchant)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbcategories.NewCategoryQueryServiceClient(conn)
	s.commandClient = pbcategories.NewCategoryCommandServiceClient(conn)
	s.statsClient = pbcategories.NewCategoryStatsServiceClient(conn)
	s.statsByIdClient = pbcategories.NewCategoryStatsByIdServiceClient(conn)
	s.statsByMerchantClient = pbcategories.NewCategoryStatsByMerchantServiceClient(conn)
}

func (s *CategoryGapiTestSuite) TestCategoryGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.commandClient.Create(ctx, &pbcategories.CreateCategoryRequest{
		Name:          "GAPI Category",
		Description:   "Testing via GRPC",
		SlugCategory:  "gapi-cat",
		ImageCategory: "gapi.jpg",
	})
	s.NoError(err)
	s.NotNil(createRes)
	catID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbcategories.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)
	s.Equal("GAPI Category", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbcategories.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbcategories.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.commandClient.Update(ctx, &pbcategories.UpdateCategoryRequest{
		CategoryId:    catID,
		Name:          "GAPI Category Updated",
		Description:   "Updated via GRPC",
		SlugCategory:  "gapi-cat-updated",
		ImageCategory: "gapi-updated.jpg",
	})
	s.NoError(err)
	s.Equal("GAPI Category Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.commandClient.TrashedCategory(ctx, &pbcategories.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbcategories.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.commandClient.RestoreCategory(ctx, &pbcategories.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 9. DeletePermanent
	_, _ = s.commandClient.TrashedCategory(ctx, &pbcategories.FindByIdCategoryRequest{Id: catID})
	_, err = s.commandClient.DeleteCategoryPermanent(ctx, &pbcategories.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 10. RestoreAll
	_, err = s.commandClient.RestoreAllCategory(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 11. DeleteAll
	_, err = s.commandClient.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func (s *CategoryGapiTestSuite) TestGapiStats() {
	ctx := context.Background()

	_, err := s.statsClient.FindMonthlyTotalPrices(ctx, &pbcategories.FindYearMonthTotalPrices{Year: 2024, Month: 4})
	s.NoError(err)

	_, err = s.statsClient.FindYearlyTotalPrices(ctx, &pbcategories.FindYearTotalPrices{Year: 2024})
	s.NoError(err)

	_, err = s.statsByIdClient.FindMonthlyTotalPricesById(ctx, &pbcategories.FindYearMonthTotalPriceById{CategoryId: 1, Year: 2024, Month: 4})
	s.NoError(err)
}

func TestCategoryGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryGapiTestSuite))
}
