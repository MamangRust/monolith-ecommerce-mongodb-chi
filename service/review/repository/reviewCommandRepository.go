package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewCommandRepository struct {
	db *mongo.Database
}

func NewReviewCommandRepository(db *mongo.Database) ReviewCommandRepository {
	return &reviewCommandRepository{db: db}
}

func (r *reviewCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("reviews")
}

func (r *reviewCommandRepository) Create(ctx context.Context, request *requests.CreateReviewRequest) (*salemongo.Review, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "reviews")
	if err != nil {
		return nil, review_errors.ErrCreateReview
	}

	review := &salemongo.Review{
		ReviewID:  id,
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Comment:   request.Comment,
		Rating:    int32(request.Rating),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if _, err := r.coll().InsertOne(ctx, review); err != nil {
		return nil, review_errors.ErrCreateReview
	}
	return review, nil
}

func (r *reviewCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewRequest) (*salemongo.Review, error) {
	now := time.Now()
	reviewID := int32(*request.ReviewID)
	filter := bson.M{"review_id": reviewID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"comment":    request.Comment,
		"rating":     int32(request.Rating),
		"updated_at": now,
	}}

	var review salemongo.Review
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&review); err != nil {
		return nil, review_errors.ErrUpdateReview
	}
	return &review, nil
}

func (r *reviewCommandRepository) Trash(ctx context.Context, reviewID int) (*salemongo.Review, error) {
	now := time.Now()
	filter := bson.M{"review_id": int32(reviewID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var review salemongo.Review
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&review); err != nil {
		return nil, review_errors.ErrTrashReview
	}
	return &review, nil
}

func (r *reviewCommandRepository) Restore(ctx context.Context, reviewID int) (*salemongo.Review, error) {
	filter := bson.M{"review_id": int32(reviewID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var review salemongo.Review
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&review); err != nil {
		return nil, review_errors.ErrRestoreReview
	}
	return &review, nil
}

func (r *reviewCommandRepository) DeletePermanent(ctx context.Context, reviewID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"review_id": int32(reviewID)}); err != nil {
		return false, review_errors.ErrDeleteReviewPermanent
	}
	return true, nil
}

func (r *reviewCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, review_errors.ErrRestoreAllReviews
	}
	return true, nil
}

func (r *reviewCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview
	}
	return true, nil
}
