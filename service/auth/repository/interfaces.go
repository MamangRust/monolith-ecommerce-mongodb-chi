package repository

import (
	"context"

	roleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user_role"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type AuthUser = salemongo.User

type UserRepository = *useradapter.Repository

type ResetTokenRepository interface {
	FindByToken(ctx context.Context, code string) (*salemongo.ResetToken, error)
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*salemongo.ResetToken, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*salemongo.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*salemongo.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*salemongo.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*salemongo.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

// AuthUserRole mirrors the user role fields.
type AuthUserRole = salemongo.UserRole

// UserRoleRepository delegates to the shared user_role adapter (command side).
type UserRoleRepository = userroleadapter.CommandRepository

// AuthRole mirrors the role fields.
type AuthRole = salemongo.Role

// RoleRepository delegates to the shared role adapter (query side).
type RoleRepository = roleadapter.QueryRepository

// Repositories aggregates all auth repositories.
type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}
