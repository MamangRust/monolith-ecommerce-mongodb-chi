package handler

import (
	"math"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	pbcarts "github.com/MamangRust/monolith-ecommerce-pb/carts"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
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

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return formatTime(*t)
}

func mapToProtoCartResponse(m interface{}) *pbcarts.CartResponse {
	switch v := m.(type) {
	case *salemongo.Cart:
		return &pbcarts.CartResponse{
			Id:        v.CartID,
			UserId:    v.UserID,
			ProductId: v.ProductID,
			Name:      v.Name,
			Price:     v.Price,
			Image:     v.Image,
			Quantity:  v.Quantity,
			Weight:    v.Weight,
			CreatedAt: formatTimePtr(v.CreatedAt),
			UpdatedAt: formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.CartListRow:
		return &pbcarts.CartResponse{
			Id:        v.CartID,
			UserId:    v.UserID,
			ProductId: v.ProductID,
			Name:      v.Name,
			Price:     v.Price,
			Image:     v.Image,
			Quantity:  v.Quantity,
			Weight:    v.Weight,
			CreatedAt: formatTime(v.CreatedAt),
			UpdatedAt: formatTime(v.UpdatedAt),
		}
	case *repository.CartCreateResult:
		return &pbcarts.CartResponse{
			Id:        v.CartID,
			UserId:    v.UserID,
			ProductId: v.ProductID,
			Name:      v.Name,
			Price:     v.Price,
			Image:     v.Image,
			Quantity:  v.Quantity,
			Weight:    v.Weight,
		}
	default:
		return nil
	}
}
