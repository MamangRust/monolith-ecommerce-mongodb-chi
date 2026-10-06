package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/service"
	pbbanners "github.com/MamangRust/monolith-ecommerce-pb/banners"
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

	observability, _ := observability.NewObservability("banner-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: observability,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbbanners.RegisterBannerQueryServiceServer(gs, h.BannerQuery)
		pbbanners.RegisterBannerCommandServiceServer(gs, h.BannerCommand)
	}

	return srv, nil
}
