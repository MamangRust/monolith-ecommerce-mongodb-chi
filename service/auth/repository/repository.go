package repository

import (
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user_role"
	"go.mongodb.org/mongo-driver/mongo"
)

// GuardOptions collects the dependency guards for each remote dependency. The
// options are applied to the adapters built inside NewRepositories, so every
// outbound gRPC call runs under a per-call deadline, bulkhead and circuit
// breaker.
type GuardOptions struct {
	User     []adapter.GuardOption
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

// Deps carries everything NewRepositories needs: the database, the user
// adapter and the role service clients. The role and user-role adapters are
// constructed here so their guards live next to them.
type Deps struct {
	Db              *mongo.Database
	User            *useradapter.Repository
	RoleQueryClient pbroles.RoleQueryServiceClient
	UserRoleClient  pbuserroles.UserRoleServiceClient
	Guard           GuardOptions
}

// NewRepositories assembles the auth repositories: the Mongo-backed
// refresh/reset token repositories plus the user, role and user-role gRPC
// adapters.
func NewRepositories(deps *Deps) *Repositories {
	for _, opt := range deps.Guard.User {
		opt(deps.User)
	}

	refreshTokenColl := deps.Db.Collection("refresh_tokens")
	resetTokenColl := deps.Db.Collection("reset_tokens")

	return &Repositories{
		User:         deps.User,
		RefreshToken: NewRefreshTokenRepository(refreshTokenColl),
		UserRole:     userroleadapter.New(deps.RoleQueryClient, deps.UserRoleClient, deps.Guard.UserRole...),
		Role:         roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...),
		ResetToken:   NewResetTokenRepository(resetTokenColl),
	}
}
