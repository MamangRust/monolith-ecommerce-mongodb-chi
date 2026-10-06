package handler

import (
	pbmerchants_policy "github.com/MamangRust/monolith-ecommerce-pb/merchants_policy"
)

type MerchantPolicyQueryHandler interface {
	pbmerchants_policy.MerchantPolicyQueryServiceServer
}

type MerchantPolicyCommandHandler interface {
	pbmerchants_policy.MerchantPolicyCommandServiceServer
}
