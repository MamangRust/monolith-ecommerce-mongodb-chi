package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantBusinessCommandRepository struct {
	db *mongo.Database
}

func NewMerchantBusinessCommandRepository(db *mongo.Database) MerchantBusinessCommandRepository {
	return &merchantBusinessCommandRepository{db: db}
}

func (r *merchantBusinessCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_business_information")
}

func derefBizTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func bizToRow(b *salemongo.MerchantBusinessInformation) *salemongo.MerchantBusinessInfoRow {
	return &salemongo.MerchantBusinessInfoRow{
		MerchantBusinessInfoID: b.MerchantBusinessInfoID,
		MerchantID:             b.MerchantID,
		BusinessType:           b.BusinessType,
		TaxID:                  b.TaxID,
		EstablishedYear:        b.EstablishedYear,
		NumberOfEmployees:      b.NumberOfEmployees,
		WebsiteUrl:             b.WebsiteUrl,
		CreatedAt:              derefBizTime(b.CreatedAt),
		UpdatedAt:              derefBizTime(b.UpdatedAt),
	}
}

func (r *merchantBusinessCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantBusinessInformationRequest) (*salemongo.MerchantBusinessInfoRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "merchant_business_information")
	if err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness
	}

	bizType := req.BusinessType
	taxID := req.TaxID
	establishedYear := int32(req.EstablishedYear)
	employees := int32(req.NumberOfEmployees)
	website := req.WebsiteUrl

	doc := &salemongo.MerchantBusinessInformation{
		MerchantBusinessInfoID: id,
		MerchantID:             int32(req.MerchantID),
		BusinessType:           &bizType,
		TaxID:                  &taxID,
		EstablishedYear:        &establishedYear,
		NumberOfEmployees:      &employees,
		WebsiteUrl:             &website,
		CreatedAt:              &now,
		UpdatedAt:              &now,
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness
	}
	return bizToRow(doc), nil
}

func (r *merchantBusinessCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantBusinessInformationRequest) (*salemongo.MerchantBusinessInfoRow, error) {
	now := time.Now()
	bizID := int32(*req.MerchantBusinessInfoID)
	filter := bson.M{"merchant_business_info_id": bizID, "deleted_at": nil}

	var doc salemongo.MerchantBusinessInformation
	if err := r.coll().FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}

	bizType := req.BusinessType
	taxID := req.TaxID
	establishedYear := int32(req.EstablishedYear)
	employees := int32(req.NumberOfEmployees)
	website := req.WebsiteUrl

	update := bson.M{"$set": bson.M{
		"business_type":       bizType,
		"tax_id":              taxID,
		"established_year":    establishedYear,
		"number_of_employees": employees,
		"website_url":         website,
		"updated_at":          now,
	}}

	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness
	}
	return bizToRow(&doc), nil
}

func (r *merchantBusinessCommandRepository) Trash(ctx context.Context, merchantID int) (*salemongo.MerchantBusinessInformation, error) {
	now := time.Now()
	filter := bson.M{"merchant_id": int32(merchantID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var doc salemongo.MerchantBusinessInformation
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchantbusiness_errors.ErrTrashMerchantBusiness
	}
	return &doc, nil
}

func (r *merchantBusinessCommandRepository) Restore(ctx context.Context, merchantID int) (*salemongo.MerchantBusinessInformation, error) {
	filter := bson.M{"merchant_id": int32(merchantID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var doc salemongo.MerchantBusinessInformation
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness
	}
	return &doc, nil
}

func (r *merchantBusinessCommandRepository) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"merchant_id": int32(merchantID)}); err != nil {
		return false, merchantbusiness_errors.ErrDeletePermanentMerchantBusiness
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchantbusiness_errors.ErrRestoreAllMerchantBusinesses
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchantbusiness_errors.ErrDeleteAllPermanentMerchantBusinesses
	}
	return true, nil
}
