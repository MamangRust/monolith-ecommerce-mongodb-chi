package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewDetailCommandRepository struct {
	db *mongo.Database
}

func NewReviewDetailCommandRepository(db *mongo.Database) ReviewDetailCommandRepository {
	return &reviewDetailCommandRepository{db: db}
}

func (r *reviewDetailCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("review_details")
}

func derefRevDetailTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func reviewDetailToRow(d *salemongo.ReviewDetail) *salemongo.ReviewDetailRow {
	return &salemongo.ReviewDetailRow{
		ReviewDetailID: d.ReviewDetailID,
		ReviewID:       d.ReviewID,
		Type:           d.Type,
		Url:            d.Url,
		Caption:        d.Caption,
		CreatedAt:      derefRevDetailTime(d.CreatedAt),
		UpdatedAt:      derefRevDetailTime(d.UpdatedAt),
	}
}

func (r *reviewDetailCommandRepository) Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*salemongo.ReviewDetailRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "review_details")
	if err != nil {
		return nil, reviewdetail_errors.ErrCreateReviewDetail
	}

	caption := request.Caption
	doc := &salemongo.ReviewDetail{
		ReviewDetailID: id,
		ReviewID:       int32(request.ReviewID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        &caption,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, reviewdetail_errors.ErrCreateReviewDetail
	}
	return reviewDetailToRow(doc), nil
}

func (r *reviewDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*salemongo.ReviewDetailRow, error) {
	now := time.Now()
	detailID := int32(*request.ReviewDetailID)
	filter := bson.M{"review_detail_id": detailID, "deleted_at": nil}

	var doc salemongo.ReviewDetail
	if err := r.coll().FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, reviewdetail_errors.ErrReviewDetailNotFound
	}

	caption := request.Caption
	update := bson.M{"$set": bson.M{
		"type":       request.Type,
		"url":        request.Url,
		"caption":    caption,
		"updated_at": now,
	}}

	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, reviewdetail_errors.ErrUpdateReviewDetail
	}
	return reviewDetailToRow(&doc), nil
}

func (r *reviewDetailCommandRepository) Trash(ctx context.Context, reviewDetailID int) (*salemongo.ReviewDetail, error) {
	now := time.Now()
	filter := bson.M{"review_detail_id": int32(reviewDetailID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var doc salemongo.ReviewDetail
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, reviewdetail_errors.ErrTrashedReviewDetail
	}
	return &doc, nil
}

func (r *reviewDetailCommandRepository) Restore(ctx context.Context, reviewDetailID int) (*salemongo.ReviewDetail, error) {
	filter := bson.M{"review_detail_id": int32(reviewDetailID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var doc salemongo.ReviewDetail
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, reviewdetail_errors.ErrRestoreReviewDetail
	}
	return &doc, nil
}

func (r *reviewDetailCommandRepository) DeletePermanent(ctx context.Context, reviewDetailID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"review_detail_id": int32(reviewDetailID)}); err != nil {
		return false, reviewdetail_errors.ErrDeleteReviewDetailPermanent
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, reviewdetail_errors.ErrRestoreAllReviewDetails
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, reviewdetail_errors.ErrDeleteAllReviewDetails
	}
	return true, nil
}
