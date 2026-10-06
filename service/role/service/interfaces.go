package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, error)
	FindByID(ctx context.Context, role_id int) (*salemongo.Role, error)
	FindByName(ctx context.Context, name string) (*salemongo.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*salemongo.Role, error)
}

type RoleCommandService interface {
	Create(ctx context.Context, request *requests.CreateRoleRequest) (*salemongo.Role, error)
	Update(ctx context.Context, request *requests.UpdateRoleRequest) (*salemongo.Role, error)
	Trash(ctx context.Context, role_id int) (*salemongo.Role, error)
	Restore(ctx context.Context, role_id int) (*salemongo.Role, error)
	DeletePermanent(ctx context.Context, role_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
