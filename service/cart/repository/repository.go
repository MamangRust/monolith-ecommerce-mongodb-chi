package repository

import (
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// GuardOptions carries the resilience guard options applied to each outbound
// gRPC adapter at construction time.
type GuardOptions struct {
	User    []adapter.GuardOption
	Product []adapter.GuardOption
}

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

// NewRepositories assembles the cart Mongo repositories and wraps the shared
// user/product gRPC clients in adapters; per-dependency guards arrive via the
// optional GuardOptions and are applied on each adapter.
func NewRepositories(db *mongo.Database,
	userQuery pbusers.UserQueryServiceClient,
	productQuery pbproducts.ProductQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		CartQuery:    NewCartQueryRepository(db.Collection("carts")),
		CartCommand:  NewCartCommandRepository(db.Collection("carts")),
		UserQuery:    useradapter.New(userQuery, nil, g.User...),
		ProductQuery: productadapter.New(productQuery, nil, g.Product...),
	}
}
