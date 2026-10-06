package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
	pbreviews_detail "github.com/MamangRust/monolith-ecommerce-pb/reviews_detail"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewDetail        ReviewDetailHandleGrpc
	ReviewDetailQuery   pbreviews_detail.ReviewDetailQueryServiceServer
	ReviewDetailCommand pbreviews_detail.ReviewDetailCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewDetailQuery:   NewReviewDetailQueryHandler(deps.Service.ReviewDetailQuery, deps.Logger),
		ReviewDetailCommand: NewReviewDetailCommandHandler(deps.Service.ReviewDetailCommand, deps.Logger),
	}
}
