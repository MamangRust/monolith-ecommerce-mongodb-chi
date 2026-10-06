package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	OrderQuery           OrderQueryService
	OrderCommand         OrderCommandService
	OrderStats           OrderStatsService
	OrderStatsByMerchant OrderStatsByMerchantService
	Outbox               OutboxService
}

type Deps struct {
	Kafka         *kafka.Kafka
	Cache         *mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		OrderQuery: NewOrderQueryService(&OrderQueryServiceDeps{
			Observability:   deps.Observability,
			Cache:           deps.Cache.OrderQueryCache,
			OrderRepository: deps.Repositories.OrderQuery,
			Logger:          deps.Logger,
		}),
		OrderCommand: NewOrderCommandService(&OrderCommandServiceDeps{
			Observability:                deps.Observability,
			Cache:                        deps.Cache.OrderCommandCache,
			UserQueryRepository:          deps.Repositories.UserQuery,
			ProductQueryRepository:       deps.Repositories.ProductQuery,
			ProductCommandRepository:     deps.Repositories.ProductCommand,
			OrderQueryRepository:         deps.Repositories.OrderQuery,
			OrderCommandRepository:       deps.Repositories.OrderCommand,
			OrderItemQueryRepository:     deps.Repositories.OrderItemQuery,
			OrderItemCommandRepository:   deps.Repositories.OrderItemCommand,
			MerchantQueryRepository:      deps.Repositories.MerchantQuery,
			ShippingAddressRepository:    deps.Repositories.ShippingAddress,
			TransactionCommandRepository: deps.Repositories.TransactionCommand,
			ShippingQueryRepository:      deps.Repositories.ShippingQuery,
			StockReservationRepository:   deps.Repositories.StockReservation,
			Outbox:                       deps.Repositories.Outbox,
			Logger:                       deps.Logger,
		}),
		OrderStats: NewOrderStatsService(&OrderStatsServiceDeps{
			Observability:        deps.Observability,
			Cache:                deps.Cache.OrderStatsCache,
			OrderStatsRepository: deps.Repositories.Stats,
			Logger:               deps.Logger,
		}),
		OrderStatsByMerchant: NewOrderStatsByMerchantService(&OrderStatsByMerchantServiceDeps{
			Observability:                  deps.Observability,
			Cache:                          deps.Cache.OrderStatsByMerchantCache,
			OrderStatsByMerchantRepository: deps.Repositories.StatsByMerchant,
			Logger:                         deps.Logger,
		}),
		Outbox: NewOutboxService(deps.Repositories.Outbox, deps.Kafka, deps.Logger),
	}
}
