package seeder

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type orderSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *orderSeeder {
	return &orderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *orderSeeder) Seed() error {
	for i := 1; i <= 8; i++ {
		order := &salemongo.Order{
			MerchantID: int32(i),
			UserID:     int32(i),
			TotalPrice: int32(10000 * i),
		}
		if _, err := r.db.Collection("orders").InsertOne(r.ctx, order); err != nil {
			r.logger.Error("failed to create order", zap.Error(err))
			return err
		}

		orderItem := &salemongo.OrderItem{
			OrderID:   order.OrderID,
			ProductID: int32(i),
			Quantity:  int32(i),
			Price:     int32(10000),
		}
		if _, err := r.db.Collection("order_items").InsertOne(r.ctx, orderItem); err != nil {
			r.logger.Error("failed to create order item", zap.Error(err))
			return err
		}
	}

	r.logger.Info("order & order-item successfully seeded")

	return nil
}
