package order_item

import (
	"context"

	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	order_item_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the contract consumers depend on for order-item reads.
type QueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*mongo.OrderItem, error)
}

// CommandRepository is the contract consumers depend on for order-item writes.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*mongo.OrderItem, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*mongo.OrderItem, error)
	Trash(ctx context.Context, orderID int) (*mongo.OrderItem, error)
	Restore(ctx context.Context, orderID int) (*mongo.OrderItem, error)
	DeletePermanent(ctx context.Context, orderID int) (bool, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	// CalculateTotalPrice is served by the order-item command service; it is
	// exposed here because it is a command-service RPC.
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated order-item query/command clients.
type Repository struct {
	query   pborder_item.OrderItemQueryServiceClient
	command pborder_item.OrderItemCommandServiceClient
	guard   *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *Repository) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// NewAdapter wraps the generated order-item clients. Resilience
// (timeout/circuit-breaker/bulkhead) is applied per call through the optional
// dependency guard; a nil guard is a passthrough.
func NewAdapter(query pborder_item.OrderItemQueryServiceClient, command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	a := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// New wraps the generated order-item query/command clients. It is the canonical
// constructor; modules build it inside their repositories with the guard
// options they were given.
func New(query pborder_item.OrderItemQueryServiceClient, command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return NewAdapter(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface.
func NewQueryAdapter(query pborder_item.OrderItemQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewAdapter(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository surface.
func NewCommandAdapter(command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return NewAdapter(nil, command, opts...)
}

func (a *Repository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*mongo.OrderItem, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponsesOrderItem, error) {
		return a.query.FindOrderItemByOrder(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}

	items := make([]*mongo.OrderItem, 0, len(res.Data))
	for _, item := range res.Data {
		items = append(items, toModel(item))
	}

	return items, nil
}

func (a *Repository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.CalculateTotalPriceResponse, error) {
		return a.command.CalculateTotalPrice(ctx, &pborder_item.CalculateTotalPriceRequest{OrderId: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrCalculateTotalPrice.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}

	total := int32(res.TotalPrice)
	return &total, nil
}

func (a *Repository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*mongo.OrderItem, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return a.command.CreateOrderItem(ctx, &pborder_item.CreateOrderItemRecordRequest{
			OrderId:   int32(req.OrderID),
			ProductId: int32(req.ProductID),
			Quantity:  int32(req.Quantity),
			Price:     int32(req.Price),
		})
	})
	if err != nil {
		return nil, order_item_errors.ErrCreateOrderItem.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*mongo.OrderItem, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return a.command.UpdateOrderItem(ctx, &pborder_item.UpdateOrderItemRecordRequest{
			OrderItemId: int32(req.OrderItemID),
			Quantity:    int32(req.Quantity),
			Price:       int32(req.Price),
		})
	})
	if err != nil {
		return nil, order_item_errors.ErrUpdateOrderItem.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) Trash(ctx context.Context, orderID int) (*mongo.OrderItem, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return a.command.TrashOrderItem(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrTrashedOrderItem.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) Restore(ctx context.Context, orderID int) (*mongo.OrderItem, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return a.command.RestoreOrderItem(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrRestoreOrderItem.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) DeletePermanent(ctx context.Context, orderID int) (bool, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItemDelete, error) {
		return a.command.DeleteOrderItemPermanent(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	if res == nil {
		return false, nil
	}
	return res.Status == "success", nil
}

func (a *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItemDelete, error) {
		return a.command.DeleteOrderItemByOrderPermanent(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	if res == nil {
		return false, nil
	}
	return res.Status == "success", nil
}

func (a *Repository) RestoreAll(ctx context.Context) (bool, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItemAll, error) {
		return a.command.RestoreAllOrdersItem(ctx, &emptypb.Empty{})
	})
	if err != nil {
		return false, order_item_errors.ErrRestoreAllOrderItem.WithInternal(err)
	}
	if res == nil {
		return false, nil
	}
	return res.Status == "success", nil
}

func (a *Repository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder_item.ApiResponseOrderItemAll, error) {
		return a.command.DeleteAllPermanentOrdersItem(ctx, &emptypb.Empty{})
	})
	if err != nil {
		return false, order_item_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	if res == nil {
		return false, nil
	}
	return res.Status == "success", nil
}

func toModel(item *pborder_item.OrderItemResponse) *mongo.OrderItem {
	if item == nil {
		return nil
	}
	return &mongo.OrderItem{
		OrderItemID: item.Id,
		OrderID:     item.OrderId,
		ProductID:   item.ProductId,
		Quantity:    item.Quantity,
		Price:       item.Price,
	}
}
