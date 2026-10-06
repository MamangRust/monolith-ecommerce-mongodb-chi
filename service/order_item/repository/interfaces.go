package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemListRow, error)
	FindActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemTrashedRow, error)
	FindOrderItemByOrder(ctx context.Context, order_id int) ([]*salemongo.OrderItemRow, error)
}

type OrderItemCommandRepository interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItemRow, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItemRow, error)
	Trash(ctx context.Context, orderItemID int) (*salemongo.OrderItem, error)
	Restore(ctx context.Context, orderItemID int) (*salemongo.OrderItem, error)
	DeletePermanent(ctx context.Context, orderItemID int) (bool, error)
	DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	CalculateTotalPrice(ctx context.Context, orderID int) (int, error)
}
