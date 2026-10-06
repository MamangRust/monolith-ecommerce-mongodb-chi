package merchantsociallinkapimapper

import (
	pbmerchants_detail "github.com/MamangRust/monolith-ecommerce-pb/merchants_detail"
	pbmerchants_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchants_social_link"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantSocialLinkBaseResponseMapper interface {
	MapMerchantSocialLink(doc *pbmerchants_detail.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse
	ToApiResponseMerchantSocialLink(doc *pbmerchants_social_link.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink
}

type MerchantSocialLinkQueryResponseMapper interface {
	MerchantSocialLinkBaseResponseMapper
}

type MerchantSocialLinkCommandResponseMapper interface {
	MerchantSocialLinkBaseResponseMapper
}
