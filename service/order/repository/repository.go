package repository

import (
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbshipping_addresses "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order_item"
	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	shippingadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/shipping_address"
	transactionadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/transaction"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	MerchantQuery      MerchantQueryRepository
	ProductQuery       ProductQueryRepository
	ProductCommand     ProductCommandRepository
	OrderItemQuery     OrderItemQueryRepository
	OrderItemCommand   OrderItemCommandRepository
	OrderQuery         OrderQueryRepository
	OrderCommand       OrderCommandRepository
	UserQuery          UserQueryRepository
	ShippingAddress    ShippingAddressCommandRepository
	TransactionCommand TransactionCommandRepository
	ShippingQuery      ShippingQueryRepository
	StockReservation   StockReservationRepository
	Stats              OrderStatsRepository
	StatsByMerchant    OrderStatsByMerchantRepository
	Outbox             OutboxRepository
}

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	Merchant    []adapter.GuardOption
	Product     []adapter.GuardOption
	OrderItem   []adapter.GuardOption
	User        []adapter.GuardOption
	Shipping    []adapter.GuardOption
	Transaction []adapter.GuardOption
}

type Deps struct {
	DB                 *mongo.Database
	MerchantQuery      pbmerchants.MerchantQueryServiceClient
	ProductQuery       pbproducts.ProductQueryServiceClient
	ProductCommand     pbproducts.ProductCommandServiceClient
	OrderItemQuery     pborder_items.OrderItemQueryServiceClient
	OrderItemCommand   pborder_items.OrderItemCommandServiceClient
	UserQuery          pbusers.UserQueryServiceClient
	ShippingCommand    pbshipping_addresses.ShippingCommandServiceClient
	ShippingQuery      pbshipping_addresses.ShippingQueryServiceClient
	TransactionCommand pbtransactions.TransactionCommandServiceClient

	Guards GuardOptions
}

// NewRepositories assembles the order repositories: Mongo-backed repositories
// for orders, stock reservations, stats and outbox, plus the shared gRPC
// adapters built here from the raw clients under their guards, so every
// outbound call runs with a per-call deadline, bulkhead and circuit breaker.
func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	// The shipping query client is optional; without it the service-level nil
	// check keeps reporting the missing dependency instead of dialing nil.
	var shippingQuery ShippingQueryRepository
	if deps.ShippingQuery != nil {
		shippingQuery = shippingadapter.NewQueryAdapter(deps.ShippingQuery, g.Shipping...)
	}

	orderItem := orderitemadapter.New(deps.OrderItemQuery, deps.OrderItemCommand, g.OrderItem...)

	return &Repositories{
		MerchantQuery:      merchantadapter.NewQueryAdapter(deps.MerchantQuery, g.Merchant...),
		ProductQuery:       productadapter.NewQueryAdapter(deps.ProductQuery, g.Product...),
		ProductCommand:     productadapter.NewCommandAdapter(deps.ProductCommand, g.Product...),
		OrderItemQuery:     orderItem,
		OrderItemCommand:   orderItem,
		OrderQuery:         NewOrderQueryRepository(deps.DB),
		OrderCommand:       NewOrderCommandRepository(deps.DB),
		UserQuery:          useradapter.NewQueryAdapter(deps.UserQuery, g.User...),
		ShippingAddress:    shippingadapter.NewCommandAdapter(deps.ShippingCommand, g.Shipping...),
		TransactionCommand: transactionadapter.New(deps.TransactionCommand, g.Transaction...),
		ShippingQuery:      shippingQuery,
		StockReservation:   NewStockReservationRepository(deps.DB),
		Stats:              NewOrderStatsRepository(deps.DB),
		StatsByMerchant:    NewOrderStatsByMerchantRepository(deps.DB),
		Outbox:             NewOutboxRepository(deps.DB),
	}
}
