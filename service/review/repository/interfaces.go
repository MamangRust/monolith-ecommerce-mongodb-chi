package repository

import (
	"context"

	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// UserQueryRepository delegates to the shared user adapter.
type UserQueryRepository = *useradapter.Repository

// ProductQueryRepository delegates to the shared product adapter.
type ProductQueryRepository = productadapter.QueryRepository

type ReviewResult struct {
	ReviewID      int32
	UserID        int32
	ProductID     int32
	Name          string
	Comment       string
	Rating        int32
	CreatedAt     *string
	UpdatedAt     *string
	DeletedAt     *string
	TotalCount    int64
	ReviewDetails interface{}
}

type ReviewQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*ReviewResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*ReviewResult, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error)
	FindByID(ctx context.Context, id int) (*ReviewResult, error)
}

type ReviewCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateReviewRequest) (*salemongo.Review, error)
	Update(ctx context.Context, request *requests.UpdateReviewRequest) (*salemongo.Review, error)
	Trash(ctx context.Context, review_id int) (*salemongo.Review, error)
	Restore(ctx context.Context, review_id int) (*salemongo.Review, error)
	DeletePermanent(ctx context.Context, review_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
