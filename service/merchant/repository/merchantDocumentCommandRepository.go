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

type merchantDocumentCommandRepository struct {
	coll *mongo.Collection
}

func NewMerchantDocumentCommandRepository(coll *mongo.Collection) *merchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{coll: coll}
}

func (r *merchantDocumentCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*salemongo.MerchantDocument, error) {
	now := time.Now()
	doc := &salemongo.MerchantDocument{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		UploadedAt:   &now,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return doc, nil
}

func (r *merchantDocumentCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*salemongo.MerchantDocument, error) {
	now := time.Now()
	filter := bson.M{"document_id": *request.DocumentID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"document_type": request.DocumentType,
		"document_url":  request.DocumentUrl,
		"note":          stringPtr(request.Note),
		"updated_at":    now,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc salemongo.MerchantDocument
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*salemongo.MerchantDocument, error) {
	now := time.Now()
	filter := bson.M{"document_id": *request.DocumentID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"status":     request.Status,
		"updated_at": now,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc salemongo.MerchantDocument
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) Trash(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error) {
	now := time.Now()
	filter := bson.M{"document_id": documentID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	var doc salemongo.MerchantDocument
	if err := r.coll.FindOne(ctx, bson.M{"document_id": documentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) Restore(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error) {
	filter := bson.M{"document_id": documentID, "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil, "updated_at": time.Now()}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	var doc salemongo.MerchantDocument
	if err := r.coll.FindOne(ctx, bson.M{"document_id": documentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"document_id": documentID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, merchant_errors.ErrMerchantNotFound
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return true, nil
}
