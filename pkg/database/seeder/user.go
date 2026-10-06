package seeder

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type userSeeder struct {
	db     *mongo.Database
	hash   hash.HashPassword
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserSeeder(db *mongo.Database, hash hash.HashPassword, ctx context.Context, logger logger.LoggerInterface) *userSeeder {
	return &userSeeder{
		db:     db,
		hash:   hash,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userSeeder) Seed() error {
	coll := r.db.Collection("users")
	for i := 1; i <= 10; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())
		rawPassword := fmt.Sprintf("password%d", i)

		hashedPassword, err := r.hash.HashPassword(rawPassword)
		if err != nil {
			r.logger.Error("failed to hash password", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to hash password for user %d: %w", i, err)
		}

		user := &salemongo.User{
			Firstname: fmt.Sprintf("User%d", i),
			Lastname:  fmt.Sprintf("Last%d", i),
			Email:     email,
			Password:  hashedPassword,
		}

		if _, err := coll.InsertOne(r.ctx, user); err != nil {
			r.logger.Error("failed to seed user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to seed user %d: %w", i, err)
		}

		if i > 5 {
			now := time.Now()
			if _, err := coll.UpdateOne(r.ctx,
				bson.M{"user_id": user.UserID},
				bson.M{"$set": bson.M{"deleted_at": now}},
			); err != nil {
				r.logger.Error("failed to trash user", zap.Int("user", i), zap.Error(err))
				return fmt.Errorf("failed to trash user %d: %w", i, err)
			}
		}
	}

	r.logger.Info("User seeding completed successfully")
	return nil
}
