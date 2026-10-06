package cache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data []*repository.OrderResult, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data []*repository.OrderResult, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data []*repository.OrderResult, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*salemongo.Order, bool)
	SetCachedOrderCache(ctx context.Context, data *salemongo.Order)

	GetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*repository.OrderResult, *int, bool)
	SetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant, data []*repository.OrderResult, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, orderID int)
	InvalidateOrderCache(ctx context.Context)
}

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.OrderMonthlyRevenueRow, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data []*salemongo.OrderMonthlyRevenueRow)
	GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*salemongo.OrderYearlyRevenueRow, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, data []*salemongo.OrderYearlyRevenueRow)
	GetMonthlyOrderCache(ctx context.Context, year int) ([]*salemongo.OrderMonthlyOrderRow, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, data []*salemongo.OrderMonthlyOrderRow)
	GetYearlyOrderCache(ctx context.Context, year int) ([]*salemongo.OrderYearlyOrderRow, bool)
	SetYearlyOrderCache(ctx context.Context, year int, data []*salemongo.OrderYearlyOrderRow)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.OrderMonthlyRevenueRow, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data []*salemongo.OrderMonthlyRevenueRow)
	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.OrderYearlyRevenueRow, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data []*salemongo.OrderYearlyRevenueRow)
	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.OrderMonthlyOrderRow, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data []*salemongo.OrderMonthlyOrderRow)
	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.OrderYearlyOrderRow, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data []*salemongo.OrderYearlyOrderRow)
}
