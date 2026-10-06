package handler

import (
	"math"
	"time"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func mapToProtoOrderItemResponse(item interface{}) *pborder_items.OrderItemResponse {
	var id, orderId, productId int32
	var quantity, price int32
	var createdAt, updatedAt string

	switch v := item.(type) {
	case *salemongo.OrderItem:
		id, orderId, productId = v.OrderItemID, v.OrderID, v.ProductID
		quantity, price = v.Quantity, v.Price
		createdAt, updatedAt = formatTimestamp(v.CreatedAt), formatTimestamp(v.UpdatedAt)
	case *salemongo.OrderItemRow:
		id, orderId, productId = v.OrderItemID, v.OrderID, v.ProductID
		quantity, price = v.Quantity, v.Price
		createdAt, updatedAt = formatTimestamp(v.CreatedAt), formatTimestamp(v.UpdatedAt)
	case *salemongo.OrderItemListRow:
		id, orderId, productId = v.OrderItemID, v.OrderID, v.ProductID
		quantity, price = v.Quantity, v.Price
		createdAt, updatedAt = formatTimestamp(v.CreatedAt), formatTimestamp(v.UpdatedAt)
	default:
		return nil
	}

	return &pborder_items.OrderItemResponse{
		Id: id, OrderId: orderId, ProductId: productId,
		Quantity: quantity, Price: price,
		CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

func mapToProtoOrderItemResponseDeleteAt(item interface{}) *pborder_items.OrderItemResponseDeleteAt {
	var res *pborder_items.OrderItemResponseDeleteAt
	var deletedAt interface{}

	switch v := item.(type) {
	case *salemongo.OrderItem:
		res = &pborder_items.OrderItemResponseDeleteAt{
			Id: v.OrderItemID, OrderId: v.OrderID, ProductId: v.ProductID,
			Quantity: v.Quantity, Price: v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt), UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.OrderItemActiveRow:
		res = &pborder_items.OrderItemResponseDeleteAt{
			Id: v.OrderItemID, OrderId: v.OrderID, ProductId: v.ProductID,
			Quantity: v.Quantity, Price: v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt), UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if deletedAt != nil {
		if dt, ok := deletedAt.(*time.Time); ok && dt != nil && !dt.IsZero() {
			res.DeletedAt = &wrapperspb.StringValue{Value: dt.Format("2006-01-02 15:04:05.000")}
		}
	}

	return res
}
