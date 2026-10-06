package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantCommandRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMerchantCommandRepository(db *mongo.Database, coll *mongo.Collection) *merchantCommandRepository {
	return &merchantCommandRepository{db: db, coll: coll}
}

func (r *merchantCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantRequest) (*salemongo.Merchant, error) {
	merchantID, err := salemongo.NextID(ctx, r.db, "merchants")
	if err != nil {
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}

	now := time.Now()
	merchant := &salemongo.Merchant{
		MerchantID:   merchantID,
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Status:       "active",
		Description:  stringPtr(request.Description),
		Address:      stringPtr(request.Address),
		ContactEmail: stringPtr(request.ContactEmail),
		ContactPhone: stringPtr(request.ContactPhone),
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	if _, err := r.coll.InsertOne(ctx, merchant); err != nil {
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}
	return merchant, nil
}

func (r *merchantCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*salemongo.Merchant, error) {
	now := time.Now()
	filter := bson.M{"merchant_id": *request.MerchantID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"name":          request.Name,
		"description":   stringPtr(request.Description),
		"address":       stringPtr(request.Address),
		"contact_email": stringPtr(request.ContactEmail),
		"contact_phone": stringPtr(request.ContactPhone),
		"status":        request.Status,
		"updated_at":    now,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var merchant salemongo.Merchant
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&merchant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) Trash(ctx context.Context, merchantID int) (*salemongo.Merchant, error) {
	now := time.Now()
	filter := bson.M{"merchant_id": merchantID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	var merchant salemongo.Merchant
	if err := r.coll.FindOne(ctx, bson.M{"merchant_id": merchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) Restore(ctx context.Context, merchantID int) (*salemongo.Merchant, error) {
	filter := bson.M{"merchant_id": merchantID, "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil, "updated_at": time.Now()}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	var merchant salemongo.Merchant
	if err := r.coll.FindOne(ctx, bson.M{"merchant_id": merchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	if _, err := r.coll.DeleteOne(ctx, bson.M{"merchant_id": merchantID, "deleted_at": bson.M{"$ne": nil}}); err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanent.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchants.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchants.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*salemongo.Merchant, error) {
	now := time.Now()
	filter := bson.M{"merchant_id": *request.MerchantID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"status": request.Status, "updated_at": now}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var merchant salemongo.Merchant
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&merchant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &merchant, nil
}
