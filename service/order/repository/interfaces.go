package repository

import (
	"context"
	"time"

	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order_item"
	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	shippingadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/shipping_address"
	transactionadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/transaction"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/outbox"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderResult struct {
	OrderID    int32
	UserID     int32
	MerchantID int32
	TotalPrice int32
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type StockReservationResult struct {
	ReservationID int32
	OrderID       int32
	ProductID     int32
	Quantity      int32
	Status        string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// UserQueryRepository delegates to the shared user adapter.
type UserQueryRepository = useradapter.QueryRepository

// ProductQueryRepository delegates to the shared product adapter.
type ProductQueryRepository = productadapter.QueryRepository

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository

// ProductCommandRepository delegates to the shared product adapter (stock
// adjustments); it is guarded like every other outbound dependency.
type ProductCommandRepository = productadapter.CommandRepository

// ShippingAddressCommandRepository delegates to the shared shipping-address
// adapter.
type ShippingAddressCommandRepository = shippingadapter.CommandRepository

// ShippingQueryRepository delegates to the shared shipping-address adapter
// (query side); the order command flow reads the current address through it.
type ShippingQueryRepository = shippingadapter.QueryRepository

// TransactionCommandRepository delegates to the shared transaction adapter.
type TransactionCommandRepository = transactionadapter.CommandRepository

// OrderItemQueryRepository is the read surface of the shared order-item
// adapter: queries plus CalculateTotalPrice, which the order-item command
// service serves.
type OrderItemQueryRepository interface {
	orderitemadapter.QueryRepository
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

// OrderItemCommandRepository delegates to the shared order-item adapter.
type OrderItemCommandRepository = orderitemadapter.CommandRepository

type OrderCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*salemongo.Order, error)
	Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*salemongo.Order, error)
	Trash(ctx context.Context, order_id int) (*salemongo.Order, error)
	Restore(ctx context.Context, order_id int) (*salemongo.Order, error)
	DeletePermanent(ctx context.Context, order_id int) (bool, error)
	DeletePermanentWithChildren(ctx context.Context, order_id int) (bool, error)
	FindTrashedByID(ctx context.Context, order_id int) (*salemongo.Order, error)
	FindTrashed(ctx context.Context) ([]*salemongo.Order, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type OrderQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*OrderResult, error)
	FindByID(ctx context.Context, order_id int) (*salemongo.Order, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, topic, key string, payload []byte) (*outbox.OutboxEvent, error)
	GetPending(ctx context.Context, limit int) ([]*outbox.OutboxEvent, error)
	Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*outbox.OutboxEvent, error)
	MarkDelivered(ctx context.Context, outboxID int64) (*outbox.OutboxEvent, error)
	MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*outbox.OutboxEvent, error)
	MarkDead(ctx context.Context, outboxID int64) (*outbox.OutboxEvent, error)
	DeleteOld(ctx context.Context, cutoff time.Time) (int64, error)
}

type StockReservationRepository interface {
	GetByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderStockReservation, error)
	Upsert(ctx context.Context, orderID, productID, quantity int) (*salemongo.OrderStockReservation, error)
	UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*salemongo.OrderStockReservation, error)
	Release(ctx context.Context, orderID, productID int) (*salemongo.OrderStockReservation, error)
	Reserve(ctx context.Context, orderID, productID int) (*salemongo.OrderStockReservation, error)
	GetReservedForTrashedOrders(ctx context.Context) ([]*salemongo.OrderStockReservation, error)
	GetReleasedForTrashedOrders(ctx context.Context) ([]*salemongo.OrderStockReservation, error)
	DeleteByOrder(ctx context.Context, orderID int) error
	DeleteByOrderProduct(ctx context.Context, orderID, productID int) error
	DeleteAllForTrashedOrders(ctx context.Context) error
	GetReleasedForActiveOrders(ctx context.Context) ([]*salemongo.OrderStockReservation, error)
	DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error)
	DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error)
}

type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.OrderMonthlyRevenueRow, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]*salemongo.OrderYearlyRevenueRow, error)
	GetMonthlyOrder(ctx context.Context, year int) ([]*salemongo.OrderMonthlyOrderRow, error)
	GetYearlyOrder(ctx context.Context, year int) ([]*salemongo.OrderYearlyOrderRow, error)
}

type OrderStatsByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.OrderMonthlyRevenueRow, error)
	GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.OrderYearlyRevenueRow, error)
	GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.OrderMonthlyOrderRow, error)
	GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.OrderYearlyOrderRow, error)
}
