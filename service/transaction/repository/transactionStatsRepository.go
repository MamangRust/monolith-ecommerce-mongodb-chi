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

type transactionStatsRepository struct {
	db *mongo.Database
}

func NewTransactionStatsRepository(db *mongo.Database) *transactionStatsRepository {
	return &transactionStatsRepository{db: db}
}

func (r *transactionStatsRepository) coll() *mongo.Collection {
	return r.db.Collection("transactions")
}

func monthName(m time.Month) string {
	return time.Date(0, m, 1, 0, 0, 0, 0, time.UTC).Format("Jan")
}

func (r *transactionStatsRepository) aggregateMonthlyAmount(ctx context.Context, status string, year, month int) (string, string, int64, int32, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at":     nil,
			"payment_status": status,
			"created_at":     bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_count":  bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return "", "", 0, 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalCount  int64 `bson:"total_count"`
		TotalAmount int32 `bson:"total_amount"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return "", "", 0, 0, err
		}
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006"),
		monthName(time.Month(month)), result.TotalCount, result.TotalAmount, nil
}

func (r *transactionStatsRepository) aggregateYearlyAmount(ctx context.Context, status string, year int) (string, int64, int32, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at":     nil,
			"payment_status": status,
			"created_at":     bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_count":  bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return "", 0, 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalCount  int64 `bson:"total_count"`
		TotalAmount int32 `bson:"total_amount"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return "", 0, 0, err
		}
	}
	return start.Format("2006"), result.TotalCount, result.TotalAmount, nil
}

func (r *transactionStatsRepository) GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, error) {
	year, total, amount, _, err := r.aggregateMonthlyAmountWrapper(ctx, "success", req.Year, req.Month)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess.WithInternal(err)
	}
	_ = total
	results := []*salemongo.TransactionMonthlyAmountSuccessRow{{
		Year:         year,
		Month:        monthName(time.Month(req.Month)),
		TotalSuccess: total,
		TotalAmount:  amount,
	}}
	return results, nil
}

func (r *transactionStatsRepository) aggregateMonthlyAmountWrapper(ctx context.Context, status string, year, month int) (string, int64, int32, string, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at":     nil,
			"payment_status": status,
			"created_at":     bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_count":  bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$amount"},
		}}},
	})
	if err != nil {
		return "", 0, 0, "", err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalCount  int64 `bson:"total_count"`
		TotalAmount int32 `bson:"total_amount"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return "", 0, 0, "", err
		}
	}
	return start.Format("2006"), result.TotalCount, result.TotalAmount, monthName(time.Month(month)), nil
}

func (r *transactionStatsRepository) GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, error) {
	year, total, amount, _, err := r.aggregateMonthlyAmountWrapper(ctx, "failed", req.Year, req.Month)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailed.WithInternal(err)
	}
	results := []*salemongo.TransactionMonthlyAmountFailedRow{{
		Year:        year,
		Month:       monthName(time.Month(req.Month)),
		TotalFailed: total,
		TotalAmount: amount,
	}}
	return results, nil
}

func (r *transactionStatsRepository) GetYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, error) {
	yearStr, total, amount, err := r.aggregateYearlyAmount(ctx, "success", year)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess.WithInternal(err)
	}
	results := []*salemongo.TransactionYearlyAmountSuccessRow{{
		Year:         yearStr,
		TotalSuccess: total,
		TotalAmount:  amount,
	}}
	return results, nil
}

func (r *transactionStatsRepository) GetYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, error) {
	yearStr, total, amount, err := r.aggregateYearlyAmount(ctx, "failed", year)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed.WithInternal(err)
	}
	results := []*salemongo.TransactionYearlyAmountFailedRow{{
		Year:        yearStr,
		TotalFailed: total,
		TotalAmount: amount,
	}}
	return results, nil
}

func (r *transactionStatsRepository) aggregateMethod(ctx context.Context, status string, year, month int) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error) {
	match := bson.M{"deleted_at": nil, "payment_status": status}
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

	var results []*salemongo.TransactionMonthlyMethodSuccessRow
	for _, row := range rows {
		results = append(results, &salemongo.TransactionMonthlyMethodSuccessRow{
			Month:             monthName(time.Month(max(month, 1))),
			Year:              time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"),
			PaymentMethod:     row.ID.Method,
			TotalTransactions: row.TotalTransactions,
			TotalAmount:       row.TotalAmount,
		})
	}
	return results, nil
}

func (r *transactionStatsRepository) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error) {
	results, err := r.aggregateMethod(ctx, "success", req.Year, req.Month)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsRepository) GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, error) {
	results, err := r.aggregateMethod(ctx, "failed", req.Year, req.Month)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsRepository) GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, error) {
	monthly, err := r.aggregateMethod(ctx, "success", year, 0)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodSuccessRow
	for _, m := range monthly {
		results = append(results, &salemongo.TransactionYearlyMethodSuccessRow{
			Year:              m.Year,
			PaymentMethod:     m.PaymentMethod,
			TotalTransactions: m.TotalTransactions,
			TotalAmount:       m.TotalAmount,
		})
	}
	return results, nil
}

func (r *transactionStatsRepository) GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, error) {
	monthly, err := r.aggregateMethod(ctx, "failed", year, 0)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodFailedRow
	for _, m := range monthly {
		results = append(results, &salemongo.TransactionYearlyMethodFailedRow{
			Year:              m.Year,
			PaymentMethod:     m.PaymentMethod,
			TotalTransactions: m.TotalTransactions,
			TotalAmount:       m.TotalAmount,
		})
	}
	return results, nil
}
