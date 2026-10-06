package seeder

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type transactionSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	transactions := []salemongo.Transaction{
		{OrderID: 1, MerchantID: 1, PaymentMethod: "credit_card", Amount: 120000, PaymentStatus: "success"},
		{OrderID: 2, MerchantID: 2, PaymentMethod: "bank_transfer", Amount: 45000, PaymentStatus: "success"},
		{OrderID: 3, MerchantID: 3, PaymentMethod: "ewallet", Amount: 78900, PaymentStatus: "success"},
		{OrderID: 4, MerchantID: 4, PaymentMethod: "cash_on_delivery", Amount: 32000, PaymentStatus: "success"},
		{OrderID: 5, MerchantID: 5, PaymentMethod: "credit_card", Amount: 99999, PaymentStatus: "failed"},
		{OrderID: 6, MerchantID: 6, PaymentMethod: "bank_transfer", Amount: 150000, PaymentStatus: "failed"},
		{OrderID: 7, MerchantID: 7, PaymentMethod: "ewallet", Amount: 51000, PaymentStatus: "failed"},
		{OrderID: 8, MerchantID: 8, PaymentMethod: "credit_card", Amount: 67500, PaymentStatus: "failed"},
	}

	for _, tx := range transactions {
		if _, err := r.db.Collection("transactions").InsertOne(r.ctx, tx); err != nil {
			r.logger.Error("failed to seed transaction", zap.Error(err))
			return err
		}
	}

	r.logger.Info("transaction successfully seeded")

	return nil
}
