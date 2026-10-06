package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
	Outbox                     OutboxService
}

type Deps struct {
	Kafka         *kafka.Kafka
	DB            *mongo.Database
	Cache         *mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		TransactionQuery: NewTransactionQueryService(&TransactionQueryServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionQueryCache,
			Repository:    deps.Repositories.TransactionQuery,
			Logger:        deps.Logger,
		}),
		TransactionCommand: NewTransactionCommandService(&TransactionCommandServiceDeps{
			Observability:      deps.Observability,
			Kafka:              deps.Kafka,
			DB:                 deps.DB,
			Cache:              deps.Cache.TransactionCommandCache,
			TransactionQuery:   deps.Repositories.TransactionQuery,
			TransactionCommand: deps.Repositories.TransactionCommand,
			UserQuery:          deps.Repositories.UserQuery,
			MerchantQuery:      deps.Repositories.MerchantQuery,
			OrderQuery:         deps.Repositories.OrderQuery,
			OrderItem:          deps.Repositories.OrderItem,
			ShippingAddress:    deps.Repositories.ShippingAddress,
			Outbox:             deps.Repositories.Outbox,
			Logger:             deps.Logger,
		}),
		TransactionStats: NewTransactionStatsService(&TransactionStatsServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionStatsCache,
			Repository:    deps.Repositories.TransactionStats,
			Logger:        deps.Logger,
		}),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(&TransactionStatsByMerchantServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionStatsByMerchantCache,
			Repository:    deps.Repositories.TransactionStatsByMerchant,
			Logger:        deps.Logger,
		}),
		Outbox: NewOutboxService(deps.Repositories.Outbox, deps.Kafka, deps.Logger),
	}
}
