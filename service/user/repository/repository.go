package repository

import (
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user_role"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        roleadapter.QueryRepository
	UserRole    userroleadapter.CommandRepository
}

type GuardOptions struct {
	UserRole []adapter.GuardOption
	Role     []adapter.GuardOption
}

type Deps struct {
	Db              *mongo.Database
	RoleQueryClient pbroles.RoleQueryServiceClient
	UserRoleClient  pbuserroles.UserRoleServiceClient
	Guard           GuardOptions
}

// NewRepositories builds the user repositories: Mongo-backed command/query on
// the users collection plus the role and user-role gRPC adapters, each
// constructed here under its guard.
func NewRepositories(deps *Deps) *Repositories {
	userColl := deps.Db.Collection("users")

	return &Repositories{
		UserCommand: NewUserCommandRepository(deps.Db, userColl),
		UserQuery:   NewUserQueryRepository(userColl),
		UserRole:    userroleadapter.New(deps.RoleQueryClient, deps.UserRoleClient, deps.Guard.UserRole...),
		Role:        roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...),
	}
}
