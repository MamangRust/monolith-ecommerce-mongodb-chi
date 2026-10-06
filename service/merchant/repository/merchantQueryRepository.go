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

type merchantQueryRepository struct {
	coll *mongo.Collection
}

func NewMerchantQueryRepository(coll *mongo.Collection) *merchantQueryRepository {
	return &merchantQueryRepository{coll: coll}
}

func fmtTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func merchantToResult(m *salemongo.Merchant, totalCount int64) *MerchantResult {
	return &MerchantResult{
		MerchantID:   m.MerchantID,
		UserID:       m.UserID,
		Name:         m.Name,
		Description:  m.Description,
		Address:      m.Address,
		ContactEmail: m.ContactEmail,
		ContactPhone: m.ContactPhone,
		Status:       m.Status,
		CreatedAt:    strPtr(fmtTime(m.CreatedAt)),
		UpdatedAt:    strPtr(fmtTime(m.UpdatedAt)),
		DeletedAt:    strPtr(fmtTime(m.DeletedAt)),
		TotalCount:   totalCount,
	}
}

func strPtr(s string) *string {
	return &s
}

func merchantFilter(req *requests.FindAllMerchant, trashed bool) bson.M {
	filter := bson.M{}
	if trashed {
		filter["deleted_at"] = bson.M{"$ne": nil}
	} else {
		filter["deleted_at"] = nil
	}
	if req.Search != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	return filter
}

func (r *merchantQueryRepository) findMerchants(ctx context.Context, req *requests.FindAllMerchant, trashed bool) ([]*MerchantResult, error) {
	offset := int64((req.Page - 1) * req.PageSize)
	filter := merchantFilter(req, trashed)

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "merchant_id", Value: -1}}).
		SetSkip(offset).
		SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var merchants []salemongo.Merchant
	if err := cursor.All(ctx, &merchants); err != nil {
		return nil, err
	}

	var results []*MerchantResult
	for i := range merchants {
		results = append(results, merchantToResult(&merchants[i], totalCount))
	}
	return results, nil
}

func (r *merchantQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	results, err := r.findMerchants(ctx, req, false)
	if err != nil {
		return nil, merchant_errors.ErrFindAllMerchants.WithInternal(err)
	}
	return results, nil
}

func (r *merchantQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	results, err := r.findMerchants(ctx, req, false)
	if err != nil {
		return nil, merchant_errors.ErrFindActiveMerchants.WithInternal(err)
	}
	return results, nil
}

func (r *merchantQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	results, err := r.findMerchants(ctx, req, true)
	if err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchants.WithInternal(err)
	}
	return results, nil
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchantID int) (*MerchantResult, error) {
	var merchant salemongo.Merchant
	err := r.coll.FindOne(ctx, bson.M{"merchant_id": merchantID, "deleted_at": nil}).Decode(&merchant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	return merchantToResult(&merchant, 1), nil
}
