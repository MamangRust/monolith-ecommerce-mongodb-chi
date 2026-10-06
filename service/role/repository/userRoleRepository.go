package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

type userRoleRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewUserRoleRepository(db *mongo.Database, coll *mongo.Collection) UserRoleRepository {
	return &userRoleRepository{db: db, coll: coll}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error) {
	userRoleID, err := salemongo.NextID(ctx, r.db, "user_roles")
	if err != nil {
		return nil, role_errors.ErrAssignRole.WithInternal(err)
	}

	now := time.Now()
	userRole := &salemongo.UserRole{
		UserRoleID: userRoleID,
		UserID:     int32(req.UserId),
		RoleID:     int32(req.RoleId),
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	_, err = r.coll.InsertOne(ctx, userRole)
	if err != nil {
		return nil, role_errors.ErrAssignRole.WithInternal(err)
	}
	return userRole, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{
		"user_id": int32(req.UserId),
		"role_id": int32(req.RoleId),
	})
	if err != nil {
		return role_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}
