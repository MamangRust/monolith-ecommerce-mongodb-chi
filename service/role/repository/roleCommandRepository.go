package repository

import (
	"context"
	"errors"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roleCommandRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewRoleCommandRepository(db *mongo.Database, coll *mongo.Collection) *roleCommandRepository {
	return &roleCommandRepository{db: db, coll: coll}
}

func (r *roleCommandRepository) Create(ctx context.Context, req *requests.CreateRoleRequest) (*salemongo.Role, error) {
	roleID, err := salemongo.NextID(ctx, r.db, "roles")
	if err != nil {
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}

	now := time.Now()
	role := &salemongo.Role{
		RoleID:    roleID,
		RoleName:  req.Name,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = r.coll.InsertOne(ctx, role)
	if err != nil {
		if salemongo.IsDuplicateKeyError(err) {
			return nil, role_errors.ErrRoleConflict
		}
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}
	return role, nil
}

func (r *roleCommandRepository) Update(ctx context.Context, req *requests.UpdateRoleRequest) (*salemongo.Role, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"role_id": *req.ID, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"role_name":  req.Name,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var role salemongo.Role
	if err := result.Decode(&role); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, role_errors.ErrRoleNotFound
		}
		if salemongo.IsDuplicateKeyError(err) {
			return nil, role_errors.ErrRoleConflict
		}
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) Trash(ctx context.Context, id int) (*salemongo.Role, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"role_id": id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var role salemongo.Role
	if err := result.Decode(&role); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, role_errors.ErrRoleNotFound
		}
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) Restore(ctx context.Context, id int) (*salemongo.Role, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"role_id": id},
		bson.M{"$unset": bson.M{"deleted_at": ""}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var role salemongo.Role
	if err := result.Decode(&role); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, role_errors.ErrRoleNotFound
		}
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) DeletePermanent(ctx context.Context, roleID int) (bool, error) {
	if _, err := r.coll.DeleteOne(ctx, bson.M{"role_id": roleID, "deleted_at": bson.M{"$ne": nil}}); err != nil {
		return false, role_errors.ErrDeleteRolePermanent.WithInternal(err)
	}
	return true, nil
}

func (r *roleCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, role_errors.ErrRestoreAllRoles.WithInternal(err)
	}
	return true, nil
}

func (r *roleCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, role_errors.ErrDeleteAllRoles.WithInternal(err)
	}
	return true, nil
}
