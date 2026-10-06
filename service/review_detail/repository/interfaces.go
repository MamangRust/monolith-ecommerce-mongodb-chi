package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailListRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailTrashedRow, *int, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.ReviewDetailRow, error)
	FindByIDTrashed(ctx context.Context, user_id int) (*salemongo.ReviewDetail, error)
}

type ReviewDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*salemongo.ReviewDetailRow, error)
	Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*salemongo.ReviewDetailRow, error)
	Trash(ctx context.Context, ReviewDetail_id int) (*salemongo.ReviewDetail, error)
	Restore(ctx context.Context, ReviewDetail_id int) (*salemongo.ReviewDetail, error)

	DeletePermanent(
		ctx context.Context,
		review_detail_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
