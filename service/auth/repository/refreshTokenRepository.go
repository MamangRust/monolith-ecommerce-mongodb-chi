package repository

import (
	"context"
	"errors"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	refreshtoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/refresh_token_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type refreshTokenRepository struct {
	coll *mongo.Collection
}

func NewRefreshTokenRepository(coll *mongo.Collection) *refreshTokenRepository {
	return &refreshTokenRepository{coll: coll}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*salemongo.RefreshToken, error) {
	var result salemongo.RefreshToken
	err := r.coll.FindOne(ctx, bson.M{"token": token}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, refreshtoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*salemongo.RefreshToken, error) {
	var result salemongo.RefreshToken
	err := r.coll.FindOne(ctx, bson.M{"user_id": int32(user_id)}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, refreshtoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*salemongo.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	now := time.Now()
	refreshToken := &salemongo.RefreshToken{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	_, err = r.coll.InsertOne(ctx, refreshToken)
	if err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken.WithInternal(err)
	}
	return refreshToken, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*salemongo.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": int32(req.UserId)},
		bson.M{"$set": bson.M{
			"token":      req.Token,
			"expiration": expirationTime,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updated salemongo.RefreshToken
	if err := result.Decode(&updated); err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}
	return &updated, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"user_id": int32(user_id)})
	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
