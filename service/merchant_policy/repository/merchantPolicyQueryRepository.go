package repository

import (
	"context"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantPoliciesQueryRepository struct {
	db *mongo.Database
}

func NewMerchantPoliciesQueryRepository(db *mongo.Database) MerchantPoliciesQueryRepository {
	return &merchantPoliciesQueryRepository{db: db}
}

func (r *merchantPoliciesQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_policies")
}

func (r *merchantPoliciesQueryRepository) findPolicies(ctx context.Context, filter bson.M, page, pageSize int) ([]*salemongo.MerchantPolicyListRow, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, merchantpolicy_errors.ErrFindMerchantPolicyByID
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "merchant_policy_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, merchantpolicy_errors.ErrFindMerchantPolicyByID
	}
	defer cursor.Close(ctx)

	var policies []*salemongo.MerchantPolicy
	if err := cursor.All(ctx, &policies); err != nil {
		return nil, merchantpolicy_errors.ErrFindMerchantPolicyByID
	}

	var results []*salemongo.MerchantPolicyListRow
	for _, p := range policies {
		results = append(results, &salemongo.MerchantPolicyListRow{
			MerchantPolicyRow: salemongo.MerchantPolicyRow{
				MerchantPolicyID: p.MerchantPolicyID,
				MerchantID:       p.MerchantID,
				PolicyType:       p.PolicyType,
				Title:            p.Title,
				Description:      p.Description,
				CreatedAt:        derefTime(p.CreatedAt),
				UpdatedAt:        derefTime(p.UpdatedAt),
			},
			TotalCount: total,
		})
	}
	return results, nil
}

func (r *merchantPoliciesQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyListRow, error) {
	return r.findPolicies(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
}

func (r *merchantPoliciesQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyActiveRow, error) {
	listRows, err := r.findPolicies(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	var active []*salemongo.MerchantPolicyActiveRow
	for _, row := range listRows {
		active = append(active, &salemongo.MerchantPolicyActiveRow{
			MerchantPolicyRow: row.MerchantPolicyRow,
			TotalCount:        row.TotalCount,
		})
	}
	return active, nil
}

func (r *merchantPoliciesQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyTrashedRow, error) {
	listRows, err := r.findPolicies(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	var trashed []*salemongo.MerchantPolicyTrashedRow
	for _, row := range listRows {
		trashed = append(trashed, &salemongo.MerchantPolicyTrashedRow{
			MerchantPolicyRow: row.MerchantPolicyRow,
			TotalCount:        row.TotalCount,
		})
	}
	return trashed, nil
}

func (r *merchantPoliciesQueryRepository) FindByID(ctx context.Context, id int) (*salemongo.MerchantPolicyRow, error) {
	var p salemongo.MerchantPolicy
	err := r.coll().FindOne(ctx, bson.M{"merchant_policy_id": int32(id)}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchantpolicy_errors.ErrMerchantPolicyNotFound
		}
		return nil, merchantpolicy_errors.ErrFindMerchantPolicyByID
	}
	return &salemongo.MerchantPolicyRow{
		MerchantPolicyID: p.MerchantPolicyID,
		MerchantID:       p.MerchantID,
		PolicyType:       p.PolicyType,
		Title:            p.Title,
		Description:      p.Description,
		CreatedAt:        derefTime(p.CreatedAt),
		UpdatedAt:        derefTime(p.UpdatedAt),
	}, nil
}
