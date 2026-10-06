package apps

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
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

	roleAddr := viper.GetString("GRPC_ROLE_ADDR")

	roleConn, err := grpc.NewClient(
		roleAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to role service: %w", err)
	}

	repos := repository.NewRepositories(&repository.Deps{
		Db:              srv.MongoDB,
		RoleQueryClient: pbroles.NewRoleQueryServiceClient(roleConn),
		UserRoleClient:  pbuserroles.NewUserRoleServiceClient(roleConn),
		Guard: repository.GuardOptions{
			Role: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			UserRole: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
		},
	})
	hashing := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("user-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Repositories:  repos,
		Hash:          hashing,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbusers.RegisterUserQueryServiceServer(gs, h.UserQuery)
		pbusers.RegisterUserCommandServiceServer(gs, h.UserCommand)
	}

	return srv, nil
}
