package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_social_link_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantSocialMediaCommandRepository struct {
	db *mongo.Database
}

func NewMerchantSocialLinkCommandRepository(db *mongo.Database) MerchantSocialLinkCommandRepository {
	return &merchantSocialMediaCommandRepository{db: db}
}

func (r *merchantSocialMediaCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_social_media_links")
}

func socialToRow(s *salemongo.MerchantSocialMediaLink) *salemongo.MerchantSocialMediaLinkRow {
	return &salemongo.MerchantSocialMediaLinkRow{
		MerchantSocialID: s.MerchantSocialID,
		MerchantDetailID: s.MerchantDetailID,
		Platform:         s.Platform,
		Url:              s.Url,
		CreatedAt:        derefDetailTime(s.CreatedAt),
		UpdatedAt:        derefDetailTime(s.UpdatedAt),
	}
}

func (r *merchantSocialMediaCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*salemongo.MerchantSocialMediaLinkRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "merchant_social_media_links")
	if err != nil {
		return nil, merchant_social_link_errors.ErrFailedCreateMerchantSocialLink
	}

	doc := &salemongo.MerchantSocialMediaLink{
		MerchantSocialID: id,
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, merchant_social_link_errors.ErrFailedCreateMerchantSocialLink
	}
	return socialToRow(doc), nil
}

func (r *merchantSocialMediaCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*salemongo.MerchantSocialMediaLinkRow, error) {
	now := time.Now()
	filter := bson.M{"merchant_social_id": int32(req.ID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"merchant_detail_id": int32(*req.MerchantDetailID),
		"platform":           req.Platform,
		"url":                req.Url,
		"updated_at":         now,
	}}

	var doc salemongo.MerchantSocialMediaLink
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return socialToRow(&doc), nil
}

func (r *merchantSocialMediaCommandRepository) Trash(ctx context.Context, id int) (bool, error) {
	now := time.Now()
	filter := bson.M{"merchant_social_id": int32(id), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var doc salemongo.MerchantSocialMediaLink
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialMediaCommandRepository) Restore(ctx context.Context, id int) (bool, error) {
	filter := bson.M{"merchant_social_id": int32(id), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var doc salemongo.MerchantSocialMediaLink
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialMediaCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"merchant_social_id": int32(id)}); err != nil {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialMediaCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	return err == nil, err
}

func (r *merchantSocialMediaCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	return err == nil, err
}
