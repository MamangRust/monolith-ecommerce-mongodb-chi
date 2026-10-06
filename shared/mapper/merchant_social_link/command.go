package merchantsociallinkapimapper

import (
	pbmerchants_detail "github.com/MamangRust/monolith-ecommerce-pb/merchants_detail"
	pbmerchants_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchants_social_link"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantSocialLinkCommandResponseMapper struct{}

func NewMerchantSocialLinkCommandResponseMapper() MerchantSocialLinkCommandResponseMapper {
	return &merchantSocialLinkCommandResponseMapper{}
}

func (m *merchantSocialLinkCommandResponseMapper) MapMerchantSocialLink(doc *pbmerchants_detail.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse {
	if doc == nil {
		return nil
	}
	return &response.MerchantSocialLinkResponse{
		ID:               int(doc.Id),
		MerchantDetailID: int(doc.MerchantDetailId),
		Platform:         doc.Platform,
		URL:              doc.Url,
	}
}

func (m *merchantSocialLinkCommandResponseMapper) ToApiResponseMerchantSocialLink(doc *pbmerchants_social_link.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink {
	return &response.ApiResponseMerchantSocialLink{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantSocialLink(doc.Data),
	}
}
