package repository

import (
	"context"
	"errors"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	resettoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/reset_token_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type resetTokenRepository struct {
	coll *mongo.Collection
}

func NewResetTokenRepository(coll *mongo.Collection) *resetTokenRepository {
	return &resetTokenRepository{coll: coll}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*salemongo.ResetToken, error) {
	var result salemongo.ResetToken
	err := r.coll.FindOne(ctx, bson.M{"token": code}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resettoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*salemongo.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}

	resetToken := &salemongo.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}

	_, err = r.coll.InsertOne(ctx, resetToken)
	if err != nil {
		return nil, resettoken_errors.ErrCreateResetToken.WithInternal(err)
	}
	return resetToken, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"user_id": int64(user_id)})
	if err != nil {
		return resettoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
