package seeder

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type roleSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	randomRoles := []string{"Cashier", "Manager", "Admin", "Supplier"}

	totalRoles := len(randomRoles)

	for i, roleName := range randomRoles {
		role := &salemongo.Role{RoleName: roleName}
		if _, err := r.db.Collection("roles").InsertOne(r.ctx, role); err != nil {
			r.logger.Error("failed to seed role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %d (%s): %w", i+1, roleName, err)
		}
	}

	r.logger.Debug("role seeded successfully", zap.Int("totalRoles", totalRoles))
	return nil
}
