package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	order_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderStatsByMerchantRepository struct {
	db *mongo.Database
}

func NewOrderStatsByMerchantRepository(db *mongo.Database) OrderStatsByMerchantRepository {
	return &orderStatsByMerchantRepository{db: db}
}

func (r *orderStatsByMerchantRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.OrderMonthlyRevenueRow, error) {
	currentStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevStart := currentStart.AddDate(0, -1, 0)
	end := currentStart.AddDate(0, 1, -1)

	match := bson.M{
		"deleted_at":  nil,
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": prevStart, "$lte": end},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$created_at"},
				"month": bson.M{"$month": "$created_at"},
			},
			"total_revenue": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id.year", Value: -1}, {Key: "_id.month", Value: -1}}}},
		{{Key: "$project", Value: bson.M{
			"_id":           0,
			"year":          bson.M{"$toString": "$_id.year"},
			"month":         bson.M{"$toString": "$_id.month"},
			"total_revenue": bson.M{"$toInt": "$total_revenue"},
		}}},
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*salemongo.OrderMonthlyRevenueRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.OrderYearlyRevenueRow, error) {
	match := bson.M{
		"deleted_at":  nil,
		"merchant_id": int32(req.MerchantID),
		"created_at": bson.M{
			"$gte": time.Date(int(req.Year)-1, 1, 1, 0, 0, 0, 0, time.UTC),
			"$lte": time.Date(int(req.Year), 12, 31, 23, 59, 59, 0, time.UTC),
		},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":           bson.M{"year": bson.M{"$year": "$created_at"}},
			"total_revenue": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id.year", Value: -1}}}},
		{{Key: "$project", Value: bson.M{
			"_id":           0,
			"year":          bson.M{"$toString": "$_id.year"},
			"total_revenue": bson.M{"$toInt": "$total_revenue"},
		}}},
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*salemongo.OrderYearlyRevenueRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *orderStatsByMerchantRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.OrderMonthlyOrderRow, error) {
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(req.Year, 12, 31, 23, 59, 59, 0, time.UTC)

	match := bson.M{
		"deleted_at":  nil,
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": start, "$lte": end},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "order_items",
			"localField":   "order_id",
			"foreignField": "order_id",
			"as":           "items",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$items", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$group", Value: bson.M{
			"_id":              bson.M{"month": bson.M{"$month": "$created_at"}},
			"order_count":      bson.M{"$addToSet": "$order_id"},
			"total_revenue":    bson.M{"$sum": "$total_price"},
			"total_items_sold": bson.M{"$sum": bson.M{"$ifNull": []interface{}{"$items.quantity", 0}}},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":              0,
			"month":            bson.M{"$toString": "$_id.month"},
			"order_count":      bson.M{"$size": "$order_count"},
			"total_revenue":    bson.M{"$toDouble": "$total_revenue"},
			"total_items_sold": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "month", Value: 1}}}},
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*salemongo.OrderMonthlyOrderRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.OrderYearlyOrderRow, error) {
	start := time.Date(int(req.Year)-4, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(int(req.Year), 12, 31, 23, 59, 59, 0, time.UTC)

	match := bson.M{
		"deleted_at":  nil,
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": start, "$lte": end},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "order_items",
			"localField":   "order_id",
			"foreignField": "order_id",
			"as":           "items",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$items", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$group", Value: bson.M{
			"_id":                  bson.M{"year": bson.M{"$year": "$created_at"}},
			"order_count":          bson.M{"$addToSet": "$order_id"},
			"total_revenue":        bson.M{"$sum": "$total_price"},
			"total_items_sold":     bson.M{"$sum": bson.M{"$ifNull": []interface{}{"$items.quantity", 0}}},
			"unique_products_sold": bson.M{"$addToSet": bson.M{"$ifNull": []interface{}{"$items.product_id", nil}}},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":                  0,
			"year":                 bson.M{"$toString": "$_id.year"},
			"order_count":          bson.M{"$size": "$order_count"},
			"total_revenue":        bson.M{"$toDouble": "$total_revenue"},
			"total_items_sold":     1,
			"unique_products_sold": bson.M{"$size": "$unique_products_sold"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}}}},
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*salemongo.OrderYearlyOrderRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant.WithInternal(err)
	}
	return results, nil
}
