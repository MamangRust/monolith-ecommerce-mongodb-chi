package repository

import (
	"context"
	"errors"
	"regexp"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userQueryRepository struct {
	coll *mongo.Collection
}

func NewUserQueryRepository(coll *mongo.Collection) *userQueryRepository {
	return &userQueryRepository{coll: coll}
}

func userSearchFilter(search string) bson.M {
	rx := bson.M{"$regex": regexp.QuoteMeta(search), "$options": "i"}
	return bson.M{"$or": bson.A{
		bson.M{"firstname": rx},
		bson.M{"lastname": rx},
		bson.M{"email": rx},
	}}
}

func (r *userQueryRepository) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter = userSearchFilter(req.Search)
		filter["deleted_at"] = nil
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*UserResult
	for cursor.Next(ctx) {
		var user salemongo.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		results = append(results, &UserResult{
			UserID:     user.UserID,
			Firstname:  user.Firstname,
			Lastname:   user.Lastname,
			Email:      user.Email,
			Password:   user.Password,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
			DeletedAt:  user.DeletedAt,
			TotalCount: total,
		})
	}
	return results, nil
}

func (r *userQueryRepository) FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	return r.FindAll(ctx, req)
}

func (r *userQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter = userSearchFilter(req.Search)
		filter["deleted_at"] = bson.M{"$ne": nil}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*UserResult
	for cursor.Next(ctx) {
		var user salemongo.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		results = append(results, &UserResult{
			UserID:     user.UserID,
			Firstname:  user.Firstname,
			Lastname:   user.Lastname,
			Email:      user.Email,
			Password:   user.Password,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
			DeletedAt:  user.DeletedAt,
			TotalCount: total,
		})
	}
	return results, nil
}

func (r *userQueryRepository) FindByID(ctx context.Context, user_id int) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"user_id": user_id, "deleted_at": nil}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByIDWithPassword(ctx context.Context, user_id int) (*salemongo.User, error) {
	return r.FindByID(ctx, user_id)
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"email": email, "deleted_at": nil}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByEmailWithPassword(ctx context.Context, email string) (*salemongo.User, error) {
	return r.FindByEmail(ctx, email)
}

func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, code string) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"verification_code": code}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}
