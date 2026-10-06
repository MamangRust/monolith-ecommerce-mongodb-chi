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

type orderQueryHandler struct {
	pborders.UnimplementedOrderQueryServiceServer
	orderQuery service.OrderQueryService
	logger     logger.LoggerInterface
}

func NewOrderQueryHandler(
	orderQuery service.OrderQueryService,
	logger logger.LoggerInterface,
) pborders.OrderQueryServiceServer {
	return &orderQueryHandler{
		orderQuery: orderQuery,
		logger:     logger,
	}
}

func (s *orderQueryHandler) FindAll(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrder, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrders := make([]*pborders.OrderResponse, len(orders))
	for i, order := range orders {
		pbOrders[i] = mapToProtoOrderResponse(order)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborders.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       pbOrders,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderQueryHandler) FindById(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrder, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborders.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully fetched order",
		Data:    mapToProtoOrderResponse(order),
	}, nil
}

func (s *orderQueryHandler) FindByActive(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrderDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrders := make([]*pborders.OrderResponseDeleteAt, len(orders))
	for i, order := range orders {
		pbOrders[i] = mapToProtoOrderResponseDeleteAt(order)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborders.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order",
		Data:       pbOrders,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderQueryHandler) FindByTrashed(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrderDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrders := make([]*pborders.OrderResponseDeleteAt, len(orders))
	for i, order := range orders {
		pbOrders[i] = mapToProtoOrderResponseDeleteAt(order)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborders.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order",
		Data:       pbOrders,
		Pagination: paginationMeta,
	}, nil
}
