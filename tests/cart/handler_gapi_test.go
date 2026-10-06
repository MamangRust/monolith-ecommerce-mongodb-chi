package cart_test

import (
	"context"
	"testing"

	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	cart_handler "github.com/MamangRust/monolith-ecommerce-grpc-cart/handler"
	cart_repo "github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	cart_service "github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	pbcarts "github.com/MamangRust/monolith-ecommerce-pb/carts"
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type CartGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbcarts.CartQueryServiceClient
	commandClient pbcarts.CartCommandServiceClient
}

func (s *CartGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Cart dependencies
	mencache := cart_cache.NewMencache(cacheStore)
	repos := cart_repo.NewRepositories(
		mongoDB,
		pbusers.NewUserQueryServiceClient(s.Conns["user"]),
		pbproducts.NewProductQueryServiceClient(s.Conns["product"]),
	)
	svc := cart_service.NewService(&cart_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := cart_handler.NewHandler(&cart_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbcarts.RegisterCartQueryServiceServer(server, handler.CartQuery)
	pbcarts.RegisterCartCommandServiceServer(server, handler.CartCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbcarts.NewCartQueryServiceClient(conn)
	s.commandClient = pbcarts.NewCartCommandServiceClient(conn)
}

func (s *CartGapiTestSuite) TestGapiLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchantID, categoryID)

	// Add to Cart
	createRes, err := s.commandClient.Create(ctx, &pbcarts.CreateCartRequest{
		UserId:    int32(userID),
		ProductId: int32(prodID),
		Quantity:  5,
	})
	s.Require().NoError(err)
	s.NotNil(createRes)

	// Get
	listRes, err := s.queryClient.FindAll(ctx, &pbcarts.FindAllCartRequest{UserId: int32(userID)})
	s.Require().NoError(err)
	s.NotEmpty(listRes.Data)
}

func TestCartGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartGapiTestSuite))
}
