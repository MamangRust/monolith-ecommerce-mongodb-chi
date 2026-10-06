package bannerapimapper

import (
	pbbanners "github.com/MamangRust/monolith-ecommerce-pb/banners"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerBaseResponseMapper interface {
	ToResponseBanner(banner *pbbanners.BannerResponse) *response.BannerResponse
	ToResponsesBanner(banners []*pbbanners.BannerResponse) []*response.BannerResponse
}

type BannerQueryResponseMapper interface {
	BannerBaseResponseMapper
	ToApiResponseBanner(pbResponse *pbbanners.ApiResponseBanner) *response.ApiResponseBanner
	ToApiResponsesBanner(pbResponse *pbbanners.ApiResponsesBanner) *response.ApiResponsesBanner
	ToApiResponsePaginationBanner(pbResponse *pbbanners.ApiResponsePaginationBanner) *response.ApiResponsePaginationBanner
	ToApiResponsePaginationBannerDeleteAt(pbResponse *pbbanners.ApiResponsePaginationBannerDeleteAt) *response.ApiResponsePaginationBannerDeleteAt
}

type BannerCommandResponseMapper interface {
	BannerBaseResponseMapper
	ToApiResponseBanner(pbResponse *pbbanners.ApiResponseBanner) *response.ApiResponseBanner
	ToResponseBannerDeleteAt(banner *pbbanners.BannerResponseDeleteAt) *response.BannerResponseDeleteAt
	ToResponsesBannerDeleteAt(banners []*pbbanners.BannerResponseDeleteAt) []*response.BannerResponseDeleteAt
	ToApiResponseBannerDeleteAt(pbResponse *pbbanners.ApiResponseBannerDeleteAt) *response.ApiResponseBannerDeleteAt
	ToApiResponseBannerDelete(pbResponse *pbbanners.ApiResponseBannerDelete) *response.ApiResponseBannerDelete
	ToApiResponseBannerAll(pbResponse *pbbanners.ApiResponseBannerAll) *response.ApiResponseBannerAll
}
