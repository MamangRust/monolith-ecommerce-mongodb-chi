package handler

import (
	"context"

	pborders "github.com/MamangRust/monolith-ecommerce-pb/orders"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderQueryHandler interface {
	pborders.OrderQueryServiceServer
}

type OrderCommandHandler interface {
	pborders.OrderCommandServiceServer
}

type OrderStatsHandler interface {
	pborders.OrderStatsServiceServer
}

type OrderStatsByMerchantHandler interface {
	pborders.OrderStatsByMerchantServiceServer
}

type OrderHandleGrpc interface {
	FindMonthlyTotalRevenue(ctx context.Context, request *pborders.FindYearMonthTotalRevenue) (*pborders.ApiResponseOrderMonthlyTotalRevenue, error)
	FindYearlyTotalRevenue(ctx context.Context, request *pborders.FindYearTotalRevenue) (*pborders.ApiResponseOrderYearlyTotalRevenue, error)

	FindMonthlyTotalRevenueByMerchant(ctx context.Context, request *pborders.FindYearMonthTotalRevenueByMerchant) (*pborders.ApiResponseOrderMonthlyTotalRevenue, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, request *pborders.FindYearTotalRevenueByMerchant) (*pborders.ApiResponseOrderYearlyTotalRevenue, error)

	FindMonthlyRevenue(ctx context.Context, request *pborders.FindYearOrder) (*pborders.ApiResponseOrderMonthly, error)
	FindYearlyRevenue(ctx context.Context, request *pborders.FindYearOrder) (*pborders.ApiResponseOrderYearly, error)

	FindMonthlyRevenueByMerchant(ctx context.Context, request *pborders.FindYearOrderByMerchant) (*pborders.ApiResponseOrderMonthly, error)
	FindYearlyRevenueByMerchant(ctx context.Context, request *pborders.FindYearOrderByMerchant) (*pborders.ApiResponseOrderYearly, error)

	FindAll(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrder, error)
	FindById(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrder, error)

	FindByActive(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrderDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrderDeleteAt, error)

	Create(ctx context.Context, request *pborders.CreateOrderRequest) (*pborders.ApiResponseOrder, error)
	Update(ctx context.Context, request *pborders.UpdateOrderRequest) (*pborders.ApiResponseOrder, error)
	TrashedOrder(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrderDeleteAt, error)
	RestoreOrder(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrderDeleteAt, error)
	DeleteOrderPermanent(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrderDelete, error)
	RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborders.ApiResponseOrderAll, error)
	DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pborders.ApiResponseOrderAll, error)
}
