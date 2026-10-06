package orderitemapimapper

import (
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemBaseResponseMapper interface {
	ToResponseOrderItem(orderItem *pborder_items.OrderItemResponse) *response.OrderItemResponse
	ToResponsesOrderItem(orderItems []*pborder_items.OrderItemResponse) []*response.OrderItemResponse
}

type OrderItemQueryResponseMapper interface {
	OrderItemBaseResponseMapper
	ToApiResponseOrderItem(pbResponse *pborder_items.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToApiResponsesOrderItem(pbResponse *pborder_items.ApiResponsesOrderItem) *response.ApiResponsesOrderItem
	ToApiResponsePaginationOrderItem(pbResponse *pborder_items.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pborder_items.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
}

type OrderItemCommandResponseMapper interface {
	OrderItemBaseResponseMapper
	ToApiResponseOrderItem(pbResponse *pborder_items.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToResponseOrderItemDeleteAt(orderItem *pborder_items.OrderItemResponseDeleteAt) *response.OrderItemResponseDeleteAt
	ToResponsesOrderItemDeleteAt(orderItems []*pborder_items.OrderItemResponseDeleteAt) []*response.OrderItemResponseDeleteAt
	ToApiResponseOrderItemDelete(pbResponse *pborder_items.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete
	ToApiResponseOrderItemAll(pbResponse *pborder_items.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pborder_items.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
}
