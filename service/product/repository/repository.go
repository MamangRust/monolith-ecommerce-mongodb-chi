package repository

import (
	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/categories"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	categoryadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	"go.mongodb.org/mongo-driver/mongo"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(db *mongo.Database, categoryClient pbcategory.CategoryQueryServiceClient, merchantClient pbmerchant.MerchantQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  categoryadapter.New(categoryClient, g.Category...),
		MerchantQuery:  merchantadapter.New(merchantClient, g.Merchant...),
	}
}
