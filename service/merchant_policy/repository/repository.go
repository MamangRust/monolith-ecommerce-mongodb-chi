package repository

import (
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	"go.mongodb.org/mongo-driver/mongo"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Merchant []adapter.GuardOption
}

type Repositories struct {
	MerchantPoliciesQuery   MerchantPoliciesQueryRepository
	MerchantPoliciesCommand MerchantPoliciesCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(db *mongo.Database, merchantQuery pbmerchants.MerchantQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantPoliciesQuery:   NewMerchantPoliciesQueryRepository(db),
		MerchantPoliciesCommand: NewMerchantPoliciesCommandRepository(db),
		MerchantQuery:           merchantadapter.New(merchantQuery, g.Merchant...),
	}
}
