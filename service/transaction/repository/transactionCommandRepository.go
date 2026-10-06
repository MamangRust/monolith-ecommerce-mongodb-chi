package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transactionCommandRepository struct {
	db *mongo.Database
}

func NewTransactionCommandRepository(db *mongo.Database) TransactionCommandRepository {
	return &transactionCommandRepository{db: db}
}

func (r *transactionCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("transactions")
}

func toTransactionResult(t *salemongo.Transaction) *TransactionResult {
	if t == nil {
		return nil
	}
	return &TransactionResult{
		TransactionID: t.TransactionID,
		OrderID:       t.OrderID,
		MerchantID:    t.MerchantID,
		PaymentMethod: t.PaymentMethod,
		Amount:        t.Amount,
		PaymentStatus: t.PaymentStatus,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		DeletedAt:     t.DeletedAt,
	}
}

func (r *transactionCommandRepository) Create(ctx context.Context, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "transactions")
	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}

	paymentStatus := "pending"
	if request.PaymentStatus != nil {
		paymentStatus = *request.PaymentStatus
	}

	txn := &salemongo.Transaction{
		TransactionID: id,
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: paymentStatus,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if _, err := r.coll().InsertOne(ctx, txn); err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}
	return toTransactionResult(txn), nil
}

func (r *transactionCommandRepository) CreateInTx(ctx context.Context, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	// MongoDB doesn't support gorm-style transactions here; insert is atomic.
	return r.Create(ctx, request)
}

func (r *transactionCommandRepository) Update(ctx context.Context, request *requests.UpdateTransactionRequest) (*TransactionResult, error) {
	now := time.Now()
	txnID := int32(*request.TransactionID)
	filter := bson.M{"transaction_id": txnID, "deleted_at": nil}

	var txn salemongo.Transaction
	if err := r.coll().FindOne(ctx, filter).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}

	paymentStatus := "pending"
	if request.PaymentStatus != nil {
		paymentStatus = *request.PaymentStatus
	}

	update := bson.M{"$set": bson.M{
		"merchant_id":    int32(request.MerchantID),
		"payment_method": request.PaymentMethod,
		"amount":         int32(request.Amount),
		"order_id":       int32(request.OrderID),
		"payment_status": paymentStatus,
		"updated_at":     now,
	}}

	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrUpdateTransaction
	}
	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) Trash(ctx context.Context, transactionID int) (*TransactionResult, error) {
	now := time.Now()
	filter := bson.M{"transaction_id": int32(transactionID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var txn salemongo.Transaction
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrTrashTransaction
	}
	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) Restore(ctx context.Context, transactionID int) (*TransactionResult, error) {
	filter := bson.M{"transaction_id": int32(transactionID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var txn salemongo.Transaction
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrRestoreTransaction
	}
	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) DeletePermanent(ctx context.Context, transactionID int) (bool, error) {
	result, err := r.coll().DeleteMany(ctx, bson.M{"transaction_id": int32(transactionID)})
	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}
	if result.DeletedCount == 0 {
		return false, transaction_errors.ErrTransactionNotFound
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"order_id": int32(orderID)}); err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}
	return true, nil
}

func (r *transactionCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactions
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent
	}
	return true, nil
}
