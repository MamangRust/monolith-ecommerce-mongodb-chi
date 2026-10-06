package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
	pbshipping_addresses "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ShippingQuery   pbshipping_addresses.ShippingQueryServiceServer
	ShippingCommand pbshipping_addresses.ShippingCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ShippingQuery:   NewShippingQueryHandler(deps.Service.ShippingAddressQuery, deps.Logger),
		ShippingCommand: NewShippingCommandHandler(deps.Service.ShippingAddressCommand, deps.Logger),
	}
}
