package handler

import (
	pbmerchants_business "github.com/MamangRust/monolith-ecommerce-pb/merchants_business"
)

type MerchantBusinessQueryHandler interface {
	pbmerchants_business.MerchantBusinessQueryServiceServer
}

type MerchantBusinessCommandHandler interface {
	pbmerchants_business.MerchantBusinessCommandServiceServer
}
