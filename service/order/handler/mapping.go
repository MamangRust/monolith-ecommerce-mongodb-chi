package handler

import (
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"math"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pborders "github.com/MamangRust/monolith-ecommerce-pb/orders"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case time.Time:
		return convert.FormatTimeRFC3339(&t)
	case string:
		return t
	}
	return ""
}

func mapToProtoOrderResponse(m interface{}) *pborders.OrderResponse {
	switch v := m.(type) {
	case *salemongo.Order:
		return &pborders.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.OrderListRow:
		return &pborders.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.OrderRow:
		return &pborders.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoOrderResponseDeleteAt(m interface{}) *pborders.OrderResponseDeleteAt {
	var res *pborders.OrderResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *salemongo.Order:
		res = &pborders.OrderResponseDeleteAt{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.OrderActiveRow:
		res = &pborders.OrderResponseDeleteAt{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if val := formatTimestamp(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}

func mapToProtoOrderMonthlyTotalRevenueResponse(m interface{}) *pborders.OrderMonthlyTotalRevenueResponse {
	switch v := m.(type) {
	case *salemongo.OrderMonthlyRevenueRow:
		return &pborders.OrderMonthlyTotalRevenueResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	default:
		return nil
	}
}

func mapToProtoOrderYearlyTotalRevenueResponse(m interface{}) *pborders.OrderYearlyTotalRevenueResponse {
	switch v := m.(type) {
	case *salemongo.OrderYearlyRevenueRow:
		return &pborders.OrderYearlyTotalRevenueResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	default:
		return nil
	}
}

func mapToProtoOrderMonthlyResponse(m interface{}) *pborders.OrderMonthlyResponse {
	switch v := m.(type) {
	case *salemongo.OrderMonthlyOrderRow:
		return &pborders.OrderMonthlyResponse{
			Month:          v.Month,
			OrderCount:     int32(v.OrderCount),
			TotalRevenue:   int32(v.TotalRevenue),
			TotalItemsSold: int32(v.TotalItemsSold),
		}
	default:
		return nil
	}
}

func mapToProtoOrderYearlyResponse(m interface{}) *pborders.OrderYearlyResponse {
	switch v := m.(type) {
	case *salemongo.OrderYearlyOrderRow:
		return &pborders.OrderYearlyResponse{
			Year:               v.Year,
			OrderCount:         int32(v.OrderCount),
			TotalRevenue:       int32(v.TotalRevenue),
			TotalItemsSold:     int32(v.TotalItemsSold),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	default:
		return nil
	}
}
