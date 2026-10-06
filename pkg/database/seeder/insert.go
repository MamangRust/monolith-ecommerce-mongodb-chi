package seeder

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func insertMany(ctx context.Context, db *mongo.Database, collName string, docs []interface{}, log logger.LoggerInterface) error {
	coll := db.Collection(collName)
	if len(docs) == 0 {
		return nil
	}
	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		log.Error("Failed to bulk insert", zap.String("collection", collName), zap.Error(err))
		return err
	}
	log.Info("Seeded collection", zap.String("collection", collName), zap.Int("count", len(docs)))
	return nil
}
