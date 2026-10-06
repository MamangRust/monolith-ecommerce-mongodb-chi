package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	CategoryQuery           pbcategories.CategoryQueryServiceServer
	CategoryCommand         pbcategories.CategoryCommandServiceServer
	CategoryStats           CategoryStatsHandler
	CategoryStatsById       CategoryStatsByIdHandler
	CategoryStatsByMerchant CategoryStatsByMerchantHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		CategoryQuery:           NewCategoryQueryHandler(deps.Service.CategoryQuery, deps.Logger),
		CategoryCommand:         NewCategoryCommandHandler(deps.Service.CategoryCommand, deps.Logger),
		CategoryStats:           NewCategoryStatsHandler(deps.Service.CategoryStats, deps.Logger),
		CategoryStatsById:       NewCategoryStatsByIdHandler(deps.Service.CategoryStatsById, deps.Logger),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantHandler(deps.Service.CategoryStatsByMerchant, deps.Logger),
	}
}
