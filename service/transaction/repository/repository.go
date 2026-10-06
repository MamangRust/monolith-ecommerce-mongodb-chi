package repository

import (
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	pborder "github.com/MamangRust/monolith-ecommerce-pb/orders"
	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order_item"
	shippingadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	TransactionCommand         TransactionCommandRepository
	TransactionQuery           TransactionQueryRepository
	TransactionStats           TransactionStatsRepository
	TransactionStatsByMerchant TransactionStatsByMerchantRepository
	OrderItem                  OrderItemRepository
	OrderQuery                 OrderQueryRepository
	MerchantQuery              MerchantQueryRepository
	ShippingAddress            ShippingAddressQueryRepository
	UserQuery                  UserQueryRepository
	Outbox                     OutboxRepository
}

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User      []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
	Shipping  []adapter.GuardOption
}

type Deps struct {
	DB             *mongo.Database
	UserQuery      pbuser.UserQueryServiceClient
	MerchantQuery  pbmerchant.MerchantQueryServiceClient
	OrderQuery     pborder.OrderQueryServiceClient
	OrderItemQuery pborder_item.OrderItemQueryServiceClient
	ShippingQuery  pbshipping_address.ShippingQueryServiceClient
	Guards         GuardOptions
}

// NewRepositories assembles the transaction repositories: Mongo-backed
// repositories for the transaction itself plus the shared gRPC adapters built
// here from the raw clients under their guards, so every outbound call runs
// with a per-call deadline, bulkhead and circuit breaker.
func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	return &Repositories{
		TransactionCommand:         NewTransactionCommandRepository(deps.DB),
		TransactionQuery:           NewTransactionQueryRepository(deps.DB),
		OrderItem:                  orderitemadapter.NewQueryAdapter(deps.OrderItemQuery, g.OrderItem...),
		OrderQuery:                 orderadapter.New(deps.OrderQuery, g.Order...),
		MerchantQuery:              merchantadapter.New(deps.MerchantQuery, g.Merchant...),
		ShippingAddress:            shippingadapter.NewQueryAdapter(deps.ShippingQuery, g.Shipping...),
		TransactionStats:           NewTransactionStatsRepository(deps.DB),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantRepository(deps.DB),
		UserQuery:                  useradapter.NewQueryAdapter(deps.UserQuery, g.User...),
		Outbox:                     NewOutboxRepository(deps.DB),
	}
}
