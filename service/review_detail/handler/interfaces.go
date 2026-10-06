package handler

import (
	"context"

	pbreviews "github.com/MamangRust/monolith-ecommerce-pb/reviews"
	pbreviews_detail "github.com/MamangRust/monolith-ecommerce-pb/reviews_detail"
)

type ReviewDetailQueryHandler interface {
	pbreviews_detail.ReviewDetailQueryServiceServer
}

type ReviewDetailCommandHandler interface {
	pbreviews_detail.ReviewDetailCommandServiceServer
}

type ReviewDetailHandleGrpc interface {
	FindAll(ctx context.Context, request *pbreviews.FindAllReviewRequest) (*pbreviews_detail.ApiResponsePaginationReviewDetails, error)
	FindById(ctx context.Context, request *pbreviews_detail.FindByIdReviewDetailRequest) (*pbreviews_detail.ApiResponseReviewDetail, error)
	FindByActive(ctx context.Context, request *pbreviews.FindAllReviewRequest) (*pbreviews_detail.ApiResponsePaginationReviewDetailsDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pbreviews.FindAllReviewRequest) (*pbreviews_detail.ApiResponsePaginationReviewDetailsDeleteAt, error)
	Create(ctx context.Context, request *pbreviews_detail.CreateReviewDetailRequest) (*pbreviews_detail.ApiResponseReviewDetail, error)
	Update(ctx context.Context, request *pbreviews_detail.UpdateReviewDetailRequest) (*pbreviews_detail.ApiResponseReviewDetail, error)
	TrashedReviewDetail(ctx context.Context, request *pbreviews_detail.FindByIdReviewDetailRequest) (*pbreviews_detail.ApiResponseReviewDetailDeleteAt, error)
	RestoreReviewDetail(ctx context.Context, request *pbreviews_detail.FindByIdReviewDetailRequest) (*pbreviews_detail.ApiResponseReviewDetailDeleteAt, error)
	DeleteReviewDetailPermanent(ctx context.Context, request *pbreviews_detail.FindByIdReviewDetailRequest) (*pbreviews.ApiResponseReviewDelete, error)
}
