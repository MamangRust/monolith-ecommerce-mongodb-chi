package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
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
	obs, _ := observability.NewObservability("category-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbcategories.RegisterCategoryQueryServiceServer(gs, h.CategoryQuery)
		pbcategories.RegisterCategoryCommandServiceServer(gs, h.CategoryCommand)
		pbcategories.RegisterCategoryStatsServiceServer(gs, h.CategoryStats)
		pbcategories.RegisterCategoryStatsByIdServiceServer(gs, h.CategoryStatsById)
		pbcategories.RegisterCategoryStatsByMerchantServiceServer(gs, h.CategoryStatsByMerchant)
	}

	return srv, nil
}
