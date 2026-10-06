package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantDetailCommandRepository struct {
	db *mongo.Database
}

func NewMerchantDetailCommandRepository(db *mongo.Database) MerchantDetailCommandRepository {
	return &merchantDetailCommandRepository{db: db}
}

func (r *merchantDetailCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_details")
}

func derefDetailTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func detailToRow(d *salemongo.MerchantDetail) *salemongo.MerchantDetailRow {
	return &salemongo.MerchantDetailRow{
		MerchantDetailID: d.MerchantDetailID,
		MerchantID:       d.MerchantID,
		DisplayName:      d.DisplayName,
		CoverImageUrl:    d.CoverImageUrl,
		LogoUrl:          d.LogoUrl,
		ShortDescription: d.ShortDescription,
		WebsiteUrl:       d.WebsiteUrl,
		CreatedAt:        derefDetailTime(d.CreatedAt),
		UpdatedAt:        derefDetailTime(d.UpdatedAt),
	}
}

func (r *merchantDetailCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantDetailRequest) (*salemongo.MerchantDetailRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "merchant_details")
	if err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail
	}

	displayName := req.DisplayName
	cover := req.CoverImageUrl
	logo := req.LogoUrl
	shortDesc := req.ShortDescription
	website := req.WebsiteUrl

	doc := &salemongo.MerchantDetail{
		MerchantDetailID: id,
		MerchantID:       int32(req.MerchantID),
		DisplayName:      &displayName,
		CoverImageUrl:    &cover,
		LogoUrl:          &logo,
		ShortDescription: &shortDesc,
		WebsiteUrl:       &website,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail
	}
	return detailToRow(doc), nil
}

func (r *merchantDetailCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantDetailRequest) (*salemongo.MerchantDetailRow, error) {
	now := time.Now()
	detailID := int32(*req.MerchantDetailID)
	filter := bson.M{"merchant_detail_id": detailID, "deleted_at": nil}

	var doc salemongo.MerchantDetail
	if err := r.coll().FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}

	displayName := req.DisplayName
	cover := req.CoverImageUrl
	logo := req.LogoUrl
	shortDesc := req.ShortDescription
	website := req.WebsiteUrl

	update := bson.M{"$set": bson.M{
		"display_name":      displayName,
		"cover_image_url":   cover,
		"logo_url":          logo,
		"short_description": shortDesc,
		"website_url":       website,
		"updated_at":        now,
	}}

	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail
	}
	return detailToRow(&doc), nil
}

func (r *merchantDetailCommandRepository) Trash(ctx context.Context, merchantID int) (*salemongo.MerchantDetail, error) {
	now := time.Now()
	filter := bson.M{"merchant_id": int32(merchantID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var doc salemongo.MerchantDetail
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchantdetail_errors.ErrTrashMerchantDetail
	}
	return &doc, nil
}

func (r *merchantDetailCommandRepository) Restore(ctx context.Context, merchantID int) (*salemongo.MerchantDetail, error) {
	filter := bson.M{"merchant_id": int32(merchantID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var doc salemongo.MerchantDetail
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail
	}
	return &doc, nil
}

func (r *merchantDetailCommandRepository) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"merchant_id": int32(merchantID)}); err != nil {
		return false, merchantdetail_errors.ErrDeletePermanentMerchantDetail
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchantdetail_errors.ErrDeleteAllPermanentMerchantDetails
	}
	return true, nil
}
