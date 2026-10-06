package repository

import (
	"context"
	"errors"
	"regexp"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roleQueryRepository struct {
	db           *mongo.Database
	roleColl     *mongo.Collection
	userRoleColl *mongo.Collection
}

func NewRoleQueryRepository(db *mongo.Database, roleColl *mongo.Collection, userRoleColl *mongo.Collection) *roleQueryRepository {
	return &roleQueryRepository{db: db, roleColl: roleColl, userRoleColl: userRoleColl}
}

func (r *roleQueryRepository) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": regexp.QuoteMeta(req.Search), "$options": "i"}
	}

	total, err := r.roleColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.roleColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*RoleResult
	for cursor.Next(ctx) {
		var role salemongo.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, err
		}
		results = append(results, &RoleResult{
			RoleID:     role.RoleID,
			RoleName:   role.RoleName,
			CreatedAt:  role.CreatedAt,
			UpdatedAt:  role.UpdatedAt,
			DeletedAt:  role.DeletedAt,
			TotalCount: total,
		})
	}
	return results, nil
}

func (r *roleQueryRepository) FindActive(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	return r.FindAll(ctx, req)
}

func (r *roleQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": regexp.QuoteMeta(req.Search), "$options": "i"}
	}

	total, err := r.roleColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.roleColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*RoleResult
	for cursor.Next(ctx) {
		var role salemongo.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, err
		}
		results = append(results, &RoleResult{
			RoleID:     role.RoleID,
			RoleName:   role.RoleName,
			CreatedAt:  role.CreatedAt,
			UpdatedAt:  role.UpdatedAt,
			DeletedAt:  role.DeletedAt,
			TotalCount: total,
		})
	}
	return results, nil
}

func (r *roleQueryRepository) FindByID(ctx context.Context, id int) (*salemongo.Role, error) {
	var role salemongo.Role
	err := r.roleColl.FindOne(ctx, bson.M{"role_id": id, "deleted_at": nil}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*salemongo.Role, error) {
	var role salemongo.Role
	err := r.roleColl.FindOne(ctx, bson.M{"role_name": name, "deleted_at": nil}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, userID int) ([]*salemongo.Role, error) {
	// Find role_ids from user_roles
	cursor, err := r.userRoleColl.Find(ctx, bson.M{"user_id": int32(userID)})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var userRoles []salemongo.UserRole
	if err := cursor.All(ctx, &userRoles); err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}

	if len(userRoles) == 0 {
		return []*salemongo.Role{}, nil
	}

	// Get role_ids
	roleIDs := make([]int32, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// Find roles
	cursor2, err := r.roleColl.Find(ctx, bson.M{
		"role_id":    bson.M{"$in": roleIDs},
		"deleted_at": nil,
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	defer cursor2.Close(ctx)

	var roles []*salemongo.Role
	if err := cursor2.All(ctx, &roles); err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roles, nil
}
