package handler

import (
	pbmerchants_award "github.com/MamangRust/monolith-ecommerce-pb/merchants_award"
)

type MerchantAwardQueryHandler interface {
	pbmerchants_award.MerchantAwardQueryServiceServer
}

type MerchantAwardCommandHandler interface {
	pbmerchants_award.MerchantAwardCommandServiceServer
}
