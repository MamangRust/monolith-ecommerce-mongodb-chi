package service

import (
	"context"

	repository "github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// F5: legacy OLTP order stats services were removed; stats are served by
// service/stats_reader from ClickHouse.

type OrderQueryService interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*repository.OrderResult, *int, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*repository.OrderResult, *int, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrder,
	) ([]*repository.OrderResult, *int, error)

	FindByID(
		ctx context.Context,
		order_id int,
	) (*salemongo.Order, error)
}

type OrderCommandService interface {
	Create(
		ctx context.Context,
		request *requests.CreateOrderRequest,
	) (*salemongo.Order, error)

	Update(
		ctx context.Context,
		request *requests.UpdateOrderRequest,
	) (*salemongo.Order, error)

	Trash(
		ctx context.Context,
		order_id int,
	) (*salemongo.Order, error)

	Restore(
		ctx context.Context,
		order_id int,
	) (*salemongo.Order, error)

	DeletePermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	// ReconcileStockReservations repairs drift between reservation status and order
	// lifecycle caused by failed compensation (e.g. released reservation on an active
	// order, or reserved reservation on a trashed order). It returns the number of
	// reservations repaired.
	ReconcileStockReservations(ctx context.Context) (*ReconcileResult, error)

	// CleanupIdempotencyRecords applies the retention policy to the idempotency ledger
	// and to released reservations of trashed orders older than the retention window.
	CleanupIdempotencyRecords(ctx context.Context, retentionDays int) (*CleanupResult, error)
}

type OrderStatsService interface {
	FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.OrderMonthlyRevenueRow, error)
	FindYearlyTotalRevenue(ctx context.Context, year int) ([]*salemongo.OrderYearlyRevenueRow, error)
	FindMonthlyOrder(ctx context.Context, year int) ([]*salemongo.OrderMonthlyOrderRow, error)
	FindYearlyOrder(ctx context.Context, year int) ([]*salemongo.OrderYearlyOrderRow, error)
}

type OrderStatsByMerchantService interface {
	FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.OrderMonthlyRevenueRow, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.OrderYearlyRevenueRow, error)
	FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.OrderMonthlyOrderRow, error)
	FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.OrderYearlyOrderRow, error)
}
