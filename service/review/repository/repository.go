package repository

import (
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

// GuardOptions collects the dependency guards applied to the adapters built
// inside NewRepositories. The user and product sides get separate breakers so
// one failure mode cannot starve the other.
type GuardOptions struct {
	User    []adapter.GuardOption
	Product []adapter.GuardOption
}

// NewRepositories assembles the review Mongo repositories and builds the shared
// user/product gRPC adapters from the raw clients under their guards, so every
// outbound call runs with a per-call deadline, bulkhead and circuit breaker.
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
		ProductQuery:  productadapter.New(productQuery, nil, g.Product...),
		ReviewQuery:   NewReviewQueryRepository(db),
		UserQuery:     useradapter.New(userQuery, nil, g.User...),
		ReviewCommand: NewReviewCommandRepository(db),
	}
}
