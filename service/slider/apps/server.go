package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/service"
	pbsliders "github.com/MamangRust/monolith-ecommerce-pb/sliders"
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
	cache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("slider-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Mencache:      cache,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbsliders.RegisterSliderQueryServiceServer(gs, h.SliderQuery)
		pbsliders.RegisterSliderCommandServiceServer(gs, h.SliderCommand)
	}

	return srv, nil
}
