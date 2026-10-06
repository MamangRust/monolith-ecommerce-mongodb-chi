package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/service"
	pborders "github.com/MamangRust/monolith-ecommerce-pb/orders"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
)

type orderStatsByMerchantHandler struct {
	pborders.UnimplementedOrderStatsByMerchantServiceServer
	orderStatsByMerchant service.OrderStatsByMerchantService
	logger               logger.LoggerInterface
}

func NewOrderStatsByMerchantHandler(
	orderStatsByMerchant service.OrderStatsByMerchantService,
	logger logger.LoggerInterface,
) OrderStatsByMerchantHandler {
	return &orderStatsByMerchantHandler{
		orderStatsByMerchant: orderStatsByMerchant,
		logger:               logger,
	}
}

func (s *orderStatsByMerchantHandler) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pborders.FindYearMonthTotalRevenueByMerchant) (*pborders.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalRevenueMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborders.OrderMonthlyTotalRevenueResponse
	for _, method := range methods {
		data = append(data, &pborders.OrderMonthlyTotalRevenueResponse{
			Year:         method.Year,
			Month:        method.Month,
			TotalRevenue: int32(method.TotalRevenue),
		})
	}

	return &pborders.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (s *orderStatsByMerchantHandler) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pborders.FindYearTotalRevenueByMerchant) (*pborders.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalRevenueMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindYearlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborders.OrderYearlyTotalRevenueResponse
	for _, method := range methods {
		data = append(data, &pborders.OrderYearlyTotalRevenueResponse{
			Year:         method.Year,
			TotalRevenue: int32(method.TotalRevenue),
		})
	}

	return &pborders.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *orderStatsByMerchantHandler) FindMonthlyRevenueByMerchant(ctx context.Context, request *pborders.FindYearOrderByMerchant) (*pborders.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindMonthlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborders.OrderMonthlyResponse
	for _, item := range res {
		data = append(data, &pborders.OrderMonthlyResponse{
			Month:          item.Month,
			OrderCount:     int32(item.OrderCount),
			TotalRevenue:   int32(item.TotalRevenue),
			TotalItemsSold: int32(item.TotalItemsSold),
		})
	}

	return &pborders.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue by merchant data retrieved",
		Data:    data,
	}, nil
}

func (s *orderStatsByMerchantHandler) FindYearlyRevenueByMerchant(ctx context.Context, request *pborders.FindYearOrderByMerchant) (*pborders.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindYearlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborders.OrderYearlyResponse
	for _, item := range res {
		data = append(data, &pborders.OrderYearlyResponse{
			Year:               item.Year,
			OrderCount:         int32(item.OrderCount),
			TotalRevenue:       int32(item.TotalRevenue),
			TotalItemsSold:     int32(item.TotalItemsSold),
			UniqueProductsSold: int32(item.UniqueProductsSold),
		})
	}

	return &pborders.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue by merchant data retrieved",
		Data:    data,
	}, nil
}

func (s *orderStatsByMerchantHandler) mustEmbedUnimplementedOrderStatsByMerchantServiceServer() {}
