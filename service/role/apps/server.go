package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-role/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
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
	obs, _ := observability.NewObservability("role-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbroles.RegisterRoleQueryServiceServer(gs, h.RoleQuery)
		pbroles.RegisterRoleCommandServiceServer(gs, h.RoleCommand)
		pbuserroles.RegisterUserRoleServiceServer(gs, h.UserRole)
	}

	return srv, nil
}
