package review_detail_test

import (
	"context"
	"testing"

	detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/cache"
	detail_handler "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/handler"
	detail_repo "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	detail_service "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
	pbreviews_detail "github.com/MamangRust/monolith-ecommerce-pb/reviews_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type ReviewDetailGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbreviews_detail.ReviewDetailQueryServiceClient
	commandClient pbreviews_detail.ReviewDetailCommandServiceClient
}

func (s *ReviewDetailGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()
	s.SetupReviewService()
	s.SetupReviewDetailService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Detail dependencies
	mencache := detail_cache.NewMencache(cacheStore)
	repos := detail_repo.NewRepositories(mongoDB)
	svc := detail_service.NewService(&detail_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := detail_handler.NewHandler(&detail_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbreviews_detail.RegisterReviewDetailQueryServiceServer(server, handler.ReviewDetailQuery)
	pbreviews_detail.RegisterReviewDetailCommandServiceServer(server, handler.ReviewDetailCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbreviews_detail.NewReviewDetailQueryServiceClient(conn)
	s.commandClient = pbreviews_detail.NewReviewDetailCommandServiceClient(conn)
}

func (s *ReviewDetailGapiTestSuite) TestGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	reviewID := s.SeedReview(ctx, userID, productID)

	// Create Detail
	createRes, err := s.commandClient.Create(ctx, &pbreviews_detail.CreateReviewDetailRequest{
		ReviewId: int32(reviewID),
		Type:     "photo",
		Url:      "http://example.com/image.jpg",
		Caption:  "GAPI Detail Comment",
	})
	s.NoError(err)
	s.NotNil(createRes)

	detailID := createRes.Data.Id

	// Get
	getRes, err := s.queryClient.FindById(ctx, &pbreviews_detail.FindByIdReviewDetailRequest{Id: detailID})
	s.NoError(err)
	s.Equal("GAPI Detail Comment", getRes.Data.Caption)
}

func TestReviewDetailGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewDetailGapiTestSuite))
}
