package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryStatsRepository struct {
	db *mongo.Database
}

func NewCategoryStatsRepository(db *mongo.Database) *categoryStatsRepository {
	return &categoryStatsRepository{db: db}
}

func (r *categoryStatsRepository) GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryRevenueRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	orders := r.db.Collection("orders")
	orderItems := r.db.Collection("order_items")

	// Compute revenue for current month
	currentRevenue, err := r.monthRevenue(ctx, orders, orderItems, currentMonthStart, currentMonthEnd)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPrice.WithInternal(err)
	}

	// Compute revenue for previous month
	prevRevenue, err := r.monthRevenue(ctx, orders, orderItems, prevMonthStart, prevMonthEnd)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPrice.WithInternal(err)
	}

	results := []*salemongo.CategoryRevenueRow{
		{Year: formatYear(currentMonthStart), Month: currentMonthStart.Month().String(), TotalRevenue: float64(currentRevenue)},
		{Year: formatYear(prevMonthStart), Month: prevMonthStart.Month().String(), TotalRevenue: float64(prevRevenue)},
	}
	return results, nil
}

func (r *categoryStatsRepository) monthRevenue(ctx context.Context, orders, orderItems *mongo.Collection, start, end time.Time) (int32, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at": nil,
			"created_at": bson.M{"$gte": start, "$lte": end},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "order_items",
			"localField":   "order_id",
			"foreignField": "order_id",
			"as":           "items",
		}}},
		{{Key: "$unwind", Value: "$items"}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$total_price"},
		}}},
	}
	cursor, err := orders.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	type result struct {
		Total int32 `bson:"total"`
	}
	var res []result
	if err := cursor.All(ctx, &res); err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, nil
	}
	return res[0].Total, nil
}

func formatYear(t time.Time) string {
	return t.Format("2006")
}

func (r *categoryStatsRepository) GetYearlyTotalPrices(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRevenueRow, error) {
	orders := r.db.Collection("orders")
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at": nil,
			"created_at": bson.M{"$gte": yearStart, "$lte": yearEnd},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$total_price"},
		}}},
	}
	cursor, err := orders.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPrices.WithInternal(err)
	}
	defer cursor.Close(ctx)

	type result struct {
		Total int32 `bson:"total"`
	}
	var res []result
	if err := cursor.All(ctx, &res); err != nil {
		return nil, category_errors.ErrGetYearlyTotalPrices.WithInternal(err)
	}

	totalRevenue := 0
	if len(res) > 0 {
		totalRevenue = int(res[0].Total)
	}
	return []*salemongo.CategoryYearlyRevenueRow{
		{Year: formatYear(yearStart), TotalRevenue: float64(totalRevenue)},
	}, nil
}

func (r *categoryStatsRepository) GetMonthPrice(ctx context.Context, year int) ([]*salemongo.CategoryStatsRow, error) {
	monthlyRows, _, err := aggregateCategoryStats(ctx, r.db, year, true, categoryStatsFilter{})
	if err != nil {
		return nil, category_errors.ErrGetMonthPrice.WithInternal(err)
	}
	return monthlyRows, nil
}

func (r *categoryStatsRepository) GetYearPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyStatsRow, error) {
	_, yearlyRows, err := aggregateCategoryStats(ctx, r.db, year, false, categoryStatsFilter{})
	if err != nil {
		return nil, category_errors.ErrGetYearPrice.WithInternal(err)
	}
	return yearlyRows, nil
}
