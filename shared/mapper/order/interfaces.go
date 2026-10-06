package orderapimapper

import (
	pborders "github.com/MamangRust/monolith-ecommerce-pb/orders"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderBaseResponseMapper interface {
	ToResponseOrder(order *pborders.OrderResponse) *response.OrderResponse
	ToResponsesOrder(orders []*pborders.OrderResponse) []*response.OrderResponse
	ToApiResponseOrder(pbResponse *pborders.ApiResponseOrder) *response.ApiResponseOrder
}

type OrderQueryResponseMapper interface {
	OrderBaseResponseMapper
	ToApiResponsesOrder(pbResponse *pborders.ApiResponsesOrder) *response.ApiResponsesOrder
	ToApiResponsePaginationOrder(pbResponse *pborders.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder
	ToApiResponsePaginationOrderDeleteAt(pbResponse *pborders.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt
}

type OrderCommandResponseMapper interface {
	OrderBaseResponseMapper
	ToResponseOrderDeleteAt(order *pborders.OrderResponseDeleteAt) *response.OrderResponseDeleteAt
	ToResponsesOrderDeleteAt(orders []*pborders.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt
	ToApiResponseOrderDeleteAt(pbResponse *pborders.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt
	ToApiResponseOrderDelete(pbResponse *pborders.ApiResponseOrderDelete) *response.ApiResponseOrderDelete
	ToApiResponseOrderAll(pbResponse *pborders.ApiResponseOrderAll) *response.ApiResponseOrderAll
}

type OrderStatsResponseMapper interface {
	ToOrderMonthlyPrice(category *pborders.OrderMonthlyResponse) *response.OrderMonthlyResponse
	ToOrderMonthlyPrices(c []*pborders.OrderMonthlyResponse) []*response.OrderMonthlyResponse
	ToOrderYearlyPrice(category *pborders.OrderYearlyResponse) *response.OrderYearlyResponse
	ToOrderYearlyPrices(c []*pborders.OrderYearlyResponse) []*response.OrderYearlyResponse
	ToResponseOrderMonthlyTotalRevenue(c *pborders.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderMonthlyTotalRevenues(c []*pborders.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenue(c *pborders.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenues(c []*pborders.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse

	ToApiResponseMonthlyOrder(pbResponse *pborders.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly
	ToApiResponseYearlyOrder(pbResponse *pborders.ApiResponseOrderYearly) *response.ApiResponseOrderYearly
	ToApiResponseMonthlyTotalRevenue(pbResponse *pborders.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue
	ToApiResponseYearlyTotalRevenue(pbResponse *pborders.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue
}
