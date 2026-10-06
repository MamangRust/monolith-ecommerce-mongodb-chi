package handler

import (
	pbmerchants_detail "github.com/MamangRust/monolith-ecommerce-pb/merchants_detail"
	pbmerchants_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchants_social_link"
)

type MerchantDetailQueryHandler interface {
	pbmerchants_detail.MerchantDetailQueryServiceServer
}

type MerchantDetailCommandHandler interface {
	pbmerchants_detail.MerchantDetailCommandServiceServer
}

type MerchantSocialLinkCommandHandler interface {
	pbmerchants_social_link.MerchantSocialCommandServiceServer
}
