package reviewapimapper

import (
	pbreviews "github.com/MamangRust/monolith-ecommerce-pb/reviews"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewBaseResponseMapper interface {
	ToResponseReview(pbResponse *pbreviews.ReviewResponse) *response.ReviewResponse
	ToResponsesReview(pbResponses []*pbreviews.ReviewResponse) []*response.ReviewResponse
	ToResponseReviewsDetail(pbResponse *pbreviews.ReviewsDetailResponse) *response.ReviewsDetailResponse
	ToResponsesReviewsDetail(pbResponses []*pbreviews.ReviewsDetailResponse) []*response.ReviewsDetailResponse
}

type ReviewQueryResponseMapper interface {
	ReviewBaseResponseMapper
	ToApiResponseReview(pbResponse *pbreviews.ApiResponseReview) *response.ApiResponseReview
	ToApiResponsesReview(pbResponse *pbreviews.ApiResponsesReview) *response.ApiResponsesReview
	ToApiResponsePaginationReview(pbResponse *pbreviews.ApiResponsePaginationReview) *response.ApiResponsePaginationReview
	ToApiResponsePaginationReviewsDetail(pbResponse *pbreviews.ApiResponsePaginationReviewDetail) *response.ApiResponsePaginationReviewsDetail
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pbreviews.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
}

type ReviewCommandResponseMapper interface {
	ReviewBaseResponseMapper
	ToResponseReviewDeleteAt(pbResponse *pbreviews.ReviewResponseDeleteAt) *response.ReviewResponseDeleteAt
	ToResponsesReviewDeleteAt(pbResponses []*pbreviews.ReviewResponseDeleteAt) []*response.ReviewResponseDeleteAt
	ToApiResponseReviewDeleteAt(pbResponse *pbreviews.ApiResponseReviewDeleteAt) *response.ApiResponseReviewDeleteAt
	ToApiResponseReviewDelete(pbResponse *pbreviews.ApiResponseReviewDelete) *response.ApiResponseReviewDelete
	ToApiResponseReviewAll(pbResponse *pbreviews.ApiResponseReviewAll) *response.ApiResponseReviewAll
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pbreviews.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
}
