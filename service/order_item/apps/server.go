package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/service"
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
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
	obs, _ := observability.NewObservability("order_item-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pborder_items.RegisterOrderItemQueryServiceServer(gs, h.OrderItemQuery)
		pborder_items.RegisterOrderItemCommandServiceServer(gs, h.OrderItemCommand)
	}

	return srv, nil
}
