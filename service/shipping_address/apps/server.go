package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
	pbshipping_addresses "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.MongoDB)
	mencache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("shipping_address-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbshipping_addresses.RegisterShippingQueryServiceServer(gs, h.ShippingQuery)
		pbshipping_addresses.RegisterShippingCommandServiceServer(gs, h.ShippingCommand)
	}

	return srv, nil
}
