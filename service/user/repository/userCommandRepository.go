package repository

import (
	"context"
	"errors"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userCommandRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewUserCommandRepository(db *mongo.Database, coll *mongo.Collection) *userCommandRepository {
	return &userCommandRepository{db: db, coll: coll}
}

func (r *userCommandRepository) Create(ctx context.Context, request *requests.CreateUserRequest) (*salemongo.User, error) {
	userID, err := salemongo.NextID(ctx, r.db, "users")
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	now := time.Now()
	user := &salemongo.User{
		UserID:    userID,
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = r.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	return user, nil
}

func (r *userCommandRepository) Update(ctx context.Context, request *requests.UpdateUserRequest) (*salemongo.User, error) {
	now := time.Now()
	update := bson.M{
		"firstname":  request.FirstName,
		"lastname":   request.LastName,
		"email":      request.Email,
		"updated_at": now,
	}
	if request.Password != "" {
		update["password"] = request.Password
	}

	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": *request.UserID, "deleted_at": nil},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": user_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"is_verified": is_verified, "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) UpdatePassword(ctx context.Context, user_id int, password string) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": user_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"password": password, "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) Trash(ctx context.Context, user_id int) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": user_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) Restore(ctx context.Context, user_id int) (*salemongo.User, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": user_id, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) DeletePermanent(ctx context.Context, user_id int) (bool, error) {
	if _, err := r.coll.DeleteOne(ctx, bson.M{"user_id": user_id, "deleted_at": bson.M{"$ne": nil}}); err != nil {
		return false, user_errors.ErrDeleteUserPermanent.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, user_errors.ErrRestoreAllUsers.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, user_errors.ErrDeleteAllUsers.WithInternal(err)
	}
	return true, nil
}
