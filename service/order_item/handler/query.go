package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/service"
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type orderItemQueryHandler struct {
	pborder_items.UnimplementedOrderItemQueryServiceServer
	orderItemService service.OrderItemQueryService
	logger           logger.LoggerInterface
}

func NewOrderItemQueryHandler(orderItemService service.OrderItemQueryService, logger logger.LoggerInterface) *orderItemQueryHandler {
	return &orderItemQueryHandler{
		orderItemService: orderItemService,
		logger:           logger,
	}
}

func (h *orderItemQueryHandler) FindAll(ctx context.Context, request *pborder_items.FindAllOrderItemRequest) (*pborder_items.ApiResponsePaginationOrderItem, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := h.orderItemService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrderItems := make([]*pborder_items.OrderItemResponse, len(orderItems))
	for i, item := range orderItems {
		pbOrderItems[i] = mapToProtoOrderItemResponse(item)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborder_items.ApiResponsePaginationOrderItem{
		Status:     "success",
		Message:    "Successfully fetched order items",
		Data:       pbOrderItems,
		Pagination: paginationMeta,
	}, nil
}

func (h *orderItemQueryHandler) FindByActive(ctx context.Context, request *pborder_items.FindAllOrderItemRequest) (*pborder_items.ApiResponsePaginationOrderItemDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := h.orderItemService.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrderItems := make([]*pborder_items.OrderItemResponseDeleteAt, len(orderItems))
	for i, item := range orderItems {
		pbOrderItems[i] = mapToProtoOrderItemResponseDeleteAt(item)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborder_items.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order items",
		Data:       pbOrderItems,
		Pagination: paginationMeta,
	}, nil
}

func (h *orderItemQueryHandler) FindByTrashed(ctx context.Context, request *pborder_items.FindAllOrderItemRequest) (*pborder_items.ApiResponsePaginationOrderItemDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := h.orderItemService.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrderItems := make([]*pborder_items.OrderItemResponseDeleteAt, len(orderItems))
	for i, item := range orderItems {
		pbOrderItems[i] = mapToProtoOrderItemResponseDeleteAt(item)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborder_items.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order items",
		Data:       pbOrderItems,
		Pagination: paginationMeta,
	}, nil
}

func (h *orderItemQueryHandler) FindOrderItemByOrder(ctx context.Context, request *pborder_items.FindByIdOrderItemRequest) (*pborder_items.ApiResponsesOrderItem, error) {
	id := int(request.GetId())

	orderItems, err := h.orderItemService.FindByOrder(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrderItems := make([]*pborder_items.OrderItemResponse, len(orderItems))
	for i, item := range orderItems {
		pbOrderItems[i] = mapToProtoOrderItemResponse(item)
	}

	return &pborder_items.ApiResponsesOrderItem{
		Status:  "success",
		Message: "Successfully fetched order items by order",
		Data:    pbOrderItems,
	}, nil
}
