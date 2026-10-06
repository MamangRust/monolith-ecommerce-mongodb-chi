package seeder

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type userRoleSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserRoleSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *userRoleSeeder {
	return &userRoleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userRoleSeeder) Seed() error {
	// Fetch users
	usersColl := r.db.Collection("users")
	cursor, err := usersColl.Find(r.ctx, bson.M{})
	if err != nil {
		r.logger.Error("failed to fetch users", zap.Error(err))
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	var users []salemongo.User
	if err := cursor.All(r.ctx, &users); err != nil {
		return fmt.Errorf("failed to decode users: %w", err)
	}
	if len(users) > 20 {
		users = users[:20]
	}

	// Fetch roles
	rolesColl := r.db.Collection("roles")
	cursor, err = rolesColl.Find(r.ctx, bson.M{})
	if err != nil {
		r.logger.Error("failed to fetch roles", zap.Error(err))
		return fmt.Errorf("failed to fetch roles: %w", err)
	}
	var roles []salemongo.Role
	if err := cursor.All(r.ctx, &roles); err != nil {
		return fmt.Errorf("failed to decode roles: %w", err)
	}
	if len(roles) > 4 {
		roles = roles[:4]
	}

	if len(users) == 0 || len(roles) == 0 {
		r.logger.Debug("no users or roles available for seeding")
		return nil
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	userRolesColl := r.db.Collection("user_roles")
	for _, user := range users {
		role := roles[rand.Intn(len(roles))]

		userRole := &salemongo.UserRole{
			UserID: user.UserID,
			RoleID: role.RoleID,
		}
		if _, err := userRolesColl.InsertOne(r.ctx, userRole); err != nil {
			r.logger.Error("failed to assign role to user", zap.String("user", user.Email), zap.String("role", role.RoleName), zap.Error(err))
			return fmt.Errorf("failed to assign role %s to user %s: %w", role.RoleName, user.Email, err)
		}
	}

	r.logger.Info("user roles assigned successfully")
	return nil
}
