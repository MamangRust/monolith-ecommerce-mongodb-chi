package handler

import (
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
)

type OrderItemQueryHandler interface {
	pborder_items.OrderItemQueryServiceServer
}

type OrderItemCommandHandler interface {
	pborder_items.OrderItemCommandServiceServer
}
