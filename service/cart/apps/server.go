package apps

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	pbcarts "github.com/MamangRust/monolith-ecommerce-pb/carts"
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := viper.GetString("GRPC_USER_ADDR")

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	productConn, err := grpc.NewClient(
		productAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	repos := repository.NewRepositories(
		srv.MongoDB,
		pbusers.NewUserQueryServiceClient(userConn),
		pbproducts.NewProductQueryServiceClient(productConn),
		repository.GuardOptions{
			User: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Product: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
		},
	)

	obs, _ := observability.NewObservability("cart-service", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbcarts.RegisterCartQueryServiceServer(gs, h.CartQuery)
		pbcarts.RegisterCartCommandServiceServer(gs, h.CartCommand)
	}

	return srv, nil
}
