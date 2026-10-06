package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/service"
	pbsliders "github.com/MamangRust/monolith-ecommerce-pb/sliders"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	SliderQuery   pbsliders.SliderQueryServiceServer
	SliderCommand pbsliders.SliderCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		SliderQuery:   NewSliderQueryHandler(deps.Service.SliderQuery, deps.Logger),
		SliderCommand: NewSliderCommandHandler(deps.Service.SliderCommand, deps.Logger),
	}
}
