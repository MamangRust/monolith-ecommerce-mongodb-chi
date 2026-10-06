package cache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryCache interface {
	GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailListRow, *int, bool)
	SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*salemongo.ReviewDetailListRow, total *int)

	GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailActiveRow, *int, bool)
	SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*salemongo.ReviewDetailActiveRow, total *int)

	GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailTrashedRow, *int, bool)
	SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*salemongo.ReviewDetailTrashedRow, total *int)

	GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*salemongo.ReviewDetailRow, bool)
	SetCachedReviewDetailCache(ctx context.Context, data *salemongo.ReviewDetailRow)

	GetCachedReviewDetailTrashedCache(ctx context.Context, reviewID int) (*salemongo.ReviewDetail, bool)
	SetCachedReviewDetailTrashedCache(ctx context.Context, data *salemongo.ReviewDetail)
}

type ReviewDetailCommandCache interface {
	DeleteReviewDetailCache(ctx context.Context, reviewID int)
}
