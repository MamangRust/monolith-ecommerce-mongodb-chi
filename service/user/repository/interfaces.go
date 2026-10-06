package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserResult struct {
	UserID     int32
	Firstname  string
	Lastname   string
	Email      string
	Password   string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type UserResultTrashed = UserResult

type UserQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.User, error)
	FindByEmail(ctx context.Context, email string) (*salemongo.User, error)
	FindByIDWithPassword(ctx context.Context, user_id int) (*salemongo.User, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*salemongo.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*salemongo.User, error)
}

type UserCommandRepository interface {
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
