package repository

import (
	"context"
	"time"

	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/order_item"
	shippingadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type TransactionResult struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	PaymentMethod string
	Amount        int32
	PaymentStatus string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	TotalCount    int64
}

// UserQueryRepository delegates to the shared user adapter.
type UserQueryRepository = useradapter.QueryRepository

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository

// OrderItemRepository delegates to the shared order-item adapter (query side).
type OrderItemRepository = orderitemadapter.QueryRepository

// OrderQueryRepository delegates to the shared order adapter.
type OrderQueryRepository = orderadapter.QueryRepository

// ShippingAddressQueryRepository delegates to the shared shipping-address
// adapter; transactions resolve the address by order via FindByOrder.
type ShippingAddressQueryRepository = shippingadapter.QueryRepository

type TransactionQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*TransactionResult, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*TransactionResult, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*TransactionResult, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*TransactionResult, error)

	FindByID(
		ctx context.Context,
		transaction_id int,
	) (*TransactionResult, error)

	FindByOrderID(
		ctx context.Context,
		order_id int,
	) (*TransactionResult, error)
}

type TransactionCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*TransactionResult, error)

	CreateInTx(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*TransactionResult, error)

	Update(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*TransactionResult, error)

	Trash(
		ctx context.Context,
		transaction_id int,
	) (*TransactionResult, error)

	Restore(
		ctx context.Context,
		transaction_id int,
	) (*TransactionResult, error)

	DeletePermanent(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type TransactionStatsRepository interface {
	GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, error)
	GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, error)
	GetYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, error)
	GetYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, error)
	GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error)
	GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, error)
	GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, error)
	GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, error)
}

type TransactionStatsByMerchantRepository interface {
	GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error)
	GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error)
	GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error)
	GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error)
	GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, error)
	GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodFailedByMerchantRow, error)
	GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodSuccessByMerchantRow, error)
	GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodFailedByMerchantRow, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, topic, key string, payload []byte) (*OutboxEventResult, error)
	CreateInTx(ctx context.Context, tx interface{}, topic, key string, payload []byte) (*OutboxEventResult, error)
	GetPending(ctx context.Context, limit int) ([]*OutboxEventResult, error)
	Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*OutboxEventResult, error)
	MarkDelivered(ctx context.Context, outboxID int64) (*OutboxEventResult, error)
	MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*OutboxEventResult, error)
	MarkDead(ctx context.Context, outboxID int64) (*OutboxEventResult, error)
	DeleteOld(ctx context.Context, cutoff time.Time) (int64, error)
}

type OutboxEventResult struct {
	OutboxID      int64
	Topic         string
	EventKey      string
	Payload       []byte
	Status        string
	Attempts      int32
	NextAttemptAt time.Time
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
