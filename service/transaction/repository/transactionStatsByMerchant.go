package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionStatsByMerchantRepository struct {
	db *mongo.Database
}

func NewTransactionStatsByMerchantRepository(db *mongo.Database) *transactionStatsByMerchantRepository {
	return &transactionStatsByMerchantRepository{db: db}
}

func (r *transactionStatsByMerchantRepository) coll() *mongo.Collection {
	return r.db.Collection("transactions")
}

func (r *transactionStatsByMerchantRepository) aggregateMonthlyAmountByMerchant(ctx context.Context, status string, year, month, merchantID int) (int64, int32, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at":     nil,
			"payment_status": status,
			"merchant_id":    int32(merchantID),
			"created_at":     bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_count":  bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalCount  int64 `bson:"total_count"`
		TotalAmount int32 `bson:"total_amount"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, 0, err
		}
	}
	return result.TotalCount, result.TotalAmount, nil
}

func (r *transactionStatsByMerchantRepository) aggregateYearlyAmountByMerchant(ctx context.Context, status string, year, merchantID int) (int64, int32, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at":     nil,
			"payment_status": status,
			"merchant_id":    int32(merchantID),
			"created_at":     bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_count":  bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalCount  int64 `bson:"total_count"`
		TotalAmount int32 `bson:"total_amount"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, 0, err
		}
	}
	return result.TotalCount, result.TotalAmount, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error) {
	total, amount, err := r.aggregateMonthlyAmountByMerchant(ctx, "success", req.Year, req.Month, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant.WithInternal(err)
	}
	results := []*salemongo.TransactionMonthlyAmountSuccessByMerchantRow{{
		Year:         time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"),
		Month:        monthName(time.Month(req.Month)),
		TotalSuccess: total,
		TotalAmount:  amount,
	}}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error) {
	total, amount, err := r.aggregateMonthlyAmountByMerchant(ctx, "failed", req.Year, req.Month, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant.WithInternal(err)
	}
	results := []*salemongo.TransactionMonthlyAmountFailedByMerchantRow{{
		Year:        time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"),
		Month:       monthName(time.Month(req.Month)),
		TotalFailed: total,
		TotalAmount: amount,
	}}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error) {
	total, amount, err := r.aggregateYearlyAmountByMerchant(ctx, "success", req.Year, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant.WithInternal(err)
	}
	results := []*salemongo.TransactionYearlyAmountSuccessByMerchantRow{{
		Year:         time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"),
		TotalSuccess: total,
		TotalAmount:  amount,
	}}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error) {
	total, amount, err := r.aggregateYearlyAmountByMerchant(ctx, "failed", req.Year, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant.WithInternal(err)
	}
	results := []*salemongo.TransactionYearlyAmountFailedByMerchantRow{{
		Year:        time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"),
		TotalFailed: total,
		TotalAmount: amount,
	}}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) aggregateMethodByMerchant(ctx context.Context, status string, year, month, merchantID int) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, error) {
	match := bson.M{"deleted_at": nil, "payment_status": status, "merchant_id": int32(merchantID)}
	if year > 0 {
		start := time.Date(year, time.Month(max(month, 1)), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		if month <= 0 {
			end = start.AddDate(1, 0, 0)
		}
		match["created_at"] = bson.M{"$gte": start, "$lt": end}
	}

	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":                bson.M{"method": "$payment_method"},
			"total_transactions": bson.M{"$sum": 1},
			"total_amount":       bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rows []struct {
		ID struct {
			Method string `bson:"method"`
		} `bson:"_id"`
		TotalTransactions int64   `bson:"total_transactions"`
		TotalAmount       float64 `bson:"total_amount"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	var results []*salemongo.TransactionMonthlyMethodSuccessByMerchantRow
	for _, row := range rows {
		results = append(results, &salemongo.TransactionMonthlyMethodSuccessByMerchantRow{
			Month:             monthName(time.Month(max(month, 1))),
			Year:              time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"),
			PaymentMethod:     row.ID.Method,
			TotalTransactions: row.TotalTransactions,
			TotalAmount:       row.TotalAmount,
		})
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, error) {
	results, err := r.aggregateMethodByMerchant(ctx, "success", req.Year, req.Month, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodFailedByMerchantRow, error) {
	results, err := r.aggregateMethodByMerchant(ctx, "failed", req.Year, req.Month, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodSuccessByMerchantRow, error) {
	monthly, err := r.aggregateMethodByMerchant(ctx, "success", req.Year, 0, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodSuccessByMerchantRow
	for _, m := range monthly {
		results = append(results, &salemongo.TransactionYearlyMethodSuccessByMerchantRow{
			Year:              m.Year,
			PaymentMethod:     m.PaymentMethod,
			TotalTransactions: m.TotalTransactions,
			TotalAmount:       m.TotalAmount,
		})
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodFailedByMerchantRow, error) {
	monthly, err := r.aggregateMethodByMerchant(ctx, "failed", req.Year, 0, req.MerchantID)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodFailedByMerchantRow
	for _, m := range monthly {
		results = append(results, &salemongo.TransactionYearlyMethodFailedByMerchantRow{
			Year:              m.Year,
			PaymentMethod:     m.PaymentMethod,
			TotalTransactions: m.TotalTransactions,
			TotalAmount:       m.TotalAmount,
		})
	}
	return results, nil
}
