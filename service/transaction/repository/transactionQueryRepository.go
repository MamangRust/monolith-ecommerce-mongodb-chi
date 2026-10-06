package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transactionQueryRepository struct {
	db *mongo.Database
}

func NewTransactionQueryRepository(db *mongo.Database) TransactionQueryRepository {
	return &transactionQueryRepository{db: db}
}

func (r *transactionQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("transactions")
}

func (r *transactionQueryRepository) findTxns(ctx context.Context, filter bson.M, page, pageSize int) ([]*TransactionResult, error) {
	page = max(page, 1)
	pageSize = max(page, 1)
	offset := (page - 1) * pageSize

	totalCount, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.Transaction
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	var results []*TransactionResult
	for _, t := range docs {
		result := toTransactionResult(t)
		result.TotalCount = totalCount
		results = append(results, result)
	}
	return results, nil
}

func transactionFilter(req *requests.FindAllTransaction, trashed bool) bson.M {
	filter := bson.M{}
	if trashed {
		filter["deleted_at"] = bson.M{"$ne": nil}
	} else {
		filter["deleted_at"] = nil
	}
	if req.Search != "" {
		filter["$or"] = bson.A{
			bson.M{"payment_method": bson.M{"$regex": req.Search, "$options": "i"}},
			bson.M{"payment_status": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	return filter
}

func (r *transactionQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	results, err := r.findTxns(ctx, transactionFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, transaction_errors.ErrFindAllTransactions.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindActive(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	results, err := r.findTxns(ctx, transactionFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, transaction_errors.ErrFindByActive.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	results, err := r.findTxns(ctx, transactionFilter(req, true), req.Page, req.PageSize)
	if err != nil {
		return nil, transaction_errors.ErrFindByTrashed.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*TransactionResult, error) {
	filter := bson.M{"merchant_id": int32(req.MerchantID), "deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = bson.A{
			bson.M{"payment_method": bson.M{"$regex": req.Search, "$options": "i"}},
			bson.M{"payment_status": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	results, err := r.findTxns(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, transaction_errors.ErrFindByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindByID(ctx context.Context, transactionID int) (*TransactionResult, error) {
	var t salemongo.Transaction
	if err := r.coll().FindOne(ctx, bson.M{"transaction_id": int32(transactionID), "deleted_at": nil}).Decode(&t); err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	return toTransactionResult(&t), nil
}

func (r *transactionQueryRepository) FindByOrderID(ctx context.Context, orderID int) (*TransactionResult, error) {
	var t salemongo.Transaction
	if err := r.coll().FindOne(ctx, bson.M{"order_id": int32(orderID), "deleted_at": nil}).Decode(&t); err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	return toTransactionResult(&t), nil
}
