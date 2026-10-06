package cache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemListRow, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*salemongo.OrderItemListRow, total *int)
	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemActiveRow, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*salemongo.OrderItemActiveRow, total *int)
	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemTrashedRow, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*salemongo.OrderItemTrashedRow, total *int)
	GetCachedOrderItems(ctx context.Context, orderID int) ([]*salemongo.OrderItemRow, bool)
	SetCachedOrderItems(ctx context.Context, data []*salemongo.OrderItemRow)
}

type OrderItemCommandCache interface {
	InvalidateOrderItemCache(ctx context.Context) error
}
