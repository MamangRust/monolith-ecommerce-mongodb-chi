package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantDocumentQueryRepository struct {
	coll *mongo.Collection
}

func NewMerchantDocumentQueryRepository(coll *mongo.Collection) *merchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{coll: coll}
}

func fmtTimePtr(t interface{ IsZero() bool }) *string { return nil }

func docToResult(d *salemongo.MerchantDocument, totalCount int64) *MerchantDocumentResult {
	res := &MerchantDocumentResult{
		DocumentID:   d.DocumentID,
		MerchantID:   d.MerchantID,
		DocumentType: d.DocumentType,
		DocumentUrl:  d.DocumentUrl,
		Status:       d.Status,
		Note:         d.Note,
		TotalCount:   totalCount,
	}
	if d.UploadedAt != nil {
		s := d.UploadedAt.Format("2006-01-02 15:04:05")
		res.UploadedAt = &s
	}
	if d.CreatedAt != nil {
		s := d.CreatedAt.Format("2006-01-02 15:04:05")
		res.CreatedAt = &s
	}
	if d.UpdatedAt != nil {
		s := d.UpdatedAt.Format("2006-01-02 15:04:05")
		res.UpdatedAt = &s
	}
	if d.DeletedAt != nil {
		s := d.DeletedAt.Format("2006-01-02 15:04:05")
		res.DeletedAt = &s
	}
	return res
}

func docFilter(req *requests.FindAllMerchantDocuments, trashed bool) bson.M {
	filter := bson.M{}
	if trashed {
		filter["deleted_at"] = bson.M{"$ne": nil}
	} else {
		filter["deleted_at"] = nil
	}
	if req.Search != "" {
		filter["$or"] = bson.A{
			bson.M{"document_type": bson.M{"$regex": req.Search, "$options": "i"}},
			bson.M{"status": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	return filter
}

func (r *merchantDocumentQueryRepository) findDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, trashed bool) ([]*MerchantDocumentResult, *int, error) {
	offset := int64((req.Page - 1) * req.PageSize)
	filter := docFilter(req, trashed)

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "document_id", Value: -1}}).
		SetSkip(offset).
		SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var docs []salemongo.MerchantDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	var results []*MerchantDocumentResult
	for i := range docs {
		results = append(results, docToResult(&docs[i], totalCount))
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *merchantDocumentQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	return r.findDocuments(ctx, req, false)
}

func (r *merchantDocumentQueryRepository) FindByID(ctx context.Context, id int) (*MerchantDocumentResult, error) {
	var doc salemongo.MerchantDocument
	err := r.coll.FindOne(ctx, bson.M{"document_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	return docToResult(&doc, 1), nil
}

func (r *merchantDocumentQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	return r.findDocuments(ctx, req, false)
}

func (r *merchantDocumentQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	return r.findDocuments(ctx, req, true)
}
