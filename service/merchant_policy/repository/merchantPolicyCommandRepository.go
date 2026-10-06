package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantPoliciesCommandRepository struct {
	db *mongo.Database
}

func NewMerchantPoliciesCommandRepository(db *mongo.Database) MerchantPoliciesCommandRepository {
	return &merchantPoliciesCommandRepository{db: db}
}

func (r *merchantPoliciesCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_policies")
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func toPolicyRow(p *salemongo.MerchantPolicy) *salemongo.MerchantPolicyRow {
	return &salemongo.MerchantPolicyRow{
		MerchantPolicyID: p.MerchantPolicyID,
		MerchantID:       p.MerchantID,
		PolicyType:       p.PolicyType,
		Title:            p.Title,
		Description:      p.Description,
		CreatedAt:        derefTime(p.CreatedAt),
		UpdatedAt:        derefTime(p.UpdatedAt),
	}
}

func (r *merchantPoliciesCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantPolicyRequest) (*salemongo.MerchantPolicyRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "merchant_policies")
	if err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy
	}
	p := &salemongo.MerchantPolicy{
		MerchantPolicyID: id,
		MerchantID:       int32(req.MerchantID),
		PolicyType:       req.PolicyType,
		Title:            req.Title,
		Description:      req.Description,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	if _, err := r.coll().InsertOne(ctx, p); err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy
	}
	return toPolicyRow(p), nil
}

func (r *merchantPoliciesCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantPolicyRequest) (*salemongo.MerchantPolicyRow, error) {
	now := time.Now()
	filter := bson.M{"merchant_policy_id": int32(*req.MerchantPolicyID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"policy_type": req.PolicyType,
		"title":       req.Title,
		"description": req.Description,
		"updated_at":  now,
	}}
	var p salemongo.MerchantPolicy
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, opts).Decode(&p); err != nil {
		return nil, merchantpolicy_errors.ErrUpdateMerchantPolicy
	}
	return toPolicyRow(&p), nil
}

func (r *merchantPoliciesCommandRepository) Trash(ctx context.Context, id int) (*salemongo.MerchantPolicy, error) {
	now := time.Now()
	filter := bson.M{"merchant_policy_id": int32(id), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var p salemongo.MerchantPolicy
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, opts).Decode(&p); err != nil {
		return nil, merchantpolicy_errors.ErrTrashedMerchantPolicy
	}
	return &p, nil
}

func (r *merchantPoliciesCommandRepository) Restore(ctx context.Context, id int) (*salemongo.MerchantPolicy, error) {
	filter := bson.M{"merchant_policy_id": int32(id), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var p salemongo.MerchantPolicy
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, opts).Decode(&p); err != nil {
		return nil, merchantpolicy_errors.ErrRestoreMerchantPolicy
	}
	return &p, nil
}

func (r *merchantPoliciesCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	res, err := r.coll().DeleteOne(ctx, bson.M{"merchant_policy_id": int32(id)})
	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteMerchantPolicyPermanent
	}
	return res.DeletedCount > 0, nil
}

func (r *merchantPoliciesCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	return err == nil, err
}

func (r *merchantPoliciesCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	return err == nil, err
}
