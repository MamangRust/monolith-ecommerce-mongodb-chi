package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// UserQueryService handles query operations related to user data.
type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindByID(ctx context.Context, id int) (*salemongo.User, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*salemongo.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*salemongo.User, error)
	FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
}

// UserCommandService handles command operations related to user management.
type UserCommandService interface {
	Create(ctx context.Context, request *requests.CreateUserRequest) (*salemongo.User, error)
	Update(ctx context.Context, request *requests.UpdateUserRequest) (*salemongo.User, error)
	UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*salemongo.User, error)
	UpdatePassword(ctx context.Context, user_id int, password string) (*salemongo.User, error)
	Trash(ctx context.Context, user_id int) (*salemongo.User, error)
	Restore(ctx context.Context, user_id int) (*salemongo.User, error)
	DeletePermanent(ctx context.Context, user_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
