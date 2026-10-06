package orderapimapper

import (
	pborders "github.com/MamangRust/monolith-ecommerce-pb/orders"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type orderCommandResponseMapper struct{}

func NewOrderCommandResponseMapper() OrderCommandResponseMapper {
	return &orderCommandResponseMapper{}
}

func (o *orderCommandResponseMapper) ToResponseOrder(order *pborders.OrderResponse) *response.OrderResponse {
	if order == nil {
		return nil
	}
	return &response.OrderResponse{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderCommandResponseMapper) ToResponsesOrder(orders []*pborders.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}
	return mappedOrders
}

func (o *orderCommandResponseMapper) ToApiResponseOrder(pbResponse *pborders.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderCommandResponseMapper) ToResponseOrderDeleteAt(order *pborders.OrderResponseDeleteAt) *response.OrderResponseDeleteAt {
	var deletedAt string
	if order.DeletedAt != nil {
		deletedAt = order.DeletedAt.Value
	}

	return &response.OrderResponseDeleteAt{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  &deletedAt,
	}
}

func (o *orderCommandResponseMapper) ToResponsesOrderDeleteAt(orders []*pborders.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrderDeleteAt(order))
	}
	return mappedOrders
}

func (o *orderCommandResponseMapper) ToApiResponseOrderDeleteAt(pbResponse *pborders.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt {
	return &response.ApiResponseOrderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderDeleteAt(pbResponse.Data),
	}
}

func (o *orderCommandResponseMapper) ToApiResponseOrderDelete(pbResponse *pborders.ApiResponseOrderDelete) *response.ApiResponseOrderDelete {
	return &response.ApiResponseOrderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderCommandResponseMapper) ToApiResponseOrderAll(pbResponse *pborders.ApiResponseOrderAll) *response.ApiResponseOrderAll {
	return &response.ApiResponseOrderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}
