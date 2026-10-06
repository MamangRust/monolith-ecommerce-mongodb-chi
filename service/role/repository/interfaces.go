package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// RoleResult is the response type for role list queries (replaces sqlc row types).
type RoleResult struct {
	RoleID     int32
	RoleName   string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type RoleQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error)
	FindActive(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error)
	FindByID(ctx context.Context, role_id int) (*salemongo.Role, error)
	FindByName(ctx context.Context, name string) (*salemongo.Role, error)
	FindByUserId(ctx context.Context, user_id int) ([]*salemongo.Role, error)
}

type RoleCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateRoleRequest) (*salemongo.Role, error)
	Update(ctx context.Context, request *requests.UpdateRoleRequest) (*salemongo.Role, error)
	Trash(ctx context.Context, role_id int) (*salemongo.Role, error)
	Restore(ctx context.Context, role_id int) (*salemongo.Role, error)
	DeletePermanent(ctx context.Context, role_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
