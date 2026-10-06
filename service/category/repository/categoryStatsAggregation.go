package repository

import (
	"context"
	"strconv"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var monthAbbreviations = [...]string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// categoryStatsFilter carries the optional narrowing filters for the stats aggregations.
type categoryStatsFilter struct {
	categoryID *int32
	merchantID *int32
}

// aggregateCategoryStats joins order items to orders, products and categories to
// compute per-category stats for the given year, grouped per month when monthly is
// true, otherwise for the whole year. Monthly rows and yearly rows are returned so
// callers can pick the shape they need.
func aggregateCategoryStats(ctx context.Context, db *mongo.Database, year int, monthly bool, filter categoryStatsFilter) ([]*salemongo.CategoryStatsRow, []*salemongo.CategoryYearlyStatsRow, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	match := bson.M{
		"deleted_at":       nil,
		"order.deleted_at": nil,
		"order.created_at": bson.M{"$gte": start, "$lte": end},
	}
	if filter.categoryID != nil {
		match["product.category_id"] = *filter.categoryID
	}
	if filter.merchantID != nil {
		match["order.merchant_id"] = *filter.merchantID
	}

	groupID := bson.M{"category_id": "$product.category_id", "category_name": "$category.name"}
	if monthly {
		groupID["month"] = bson.M{"$month": "$order.created_at"}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from": "orders", "localField": "order_id", "foreignField": "order_id", "as": "order",
		}}},
		{{Key: "$unwind", Value: "$order"}},
		{{Key: "$lookup", Value: bson.M{
			"from": "products", "localField": "product_id", "foreignField": "product_id", "as": "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$match", Value: match}},
		{{Key: "$lookup", Value: bson.M{
			"from": "categories", "localField": "product.category_id", "foreignField": "category_id", "as": "category",
		}}},
		{{Key: "$unwind", Value: "$category"}},
		{{Key: "$group", Value: bson.M{
			"_id":           groupID,
			"order_ids":     bson.M{"$addToSet": "$order_id"},
			"product_ids":   bson.M{"$addToSet": "$product_id"},
			"items_sold":    bson.M{"$sum": "$quantity"},
			"total_revenue": bson.M{"$sum": "$order.total_price"},
		}}},
		{{Key: "$project", Value: bson.M{
			"category_id":          "$_id.category_id",
			"category_name":        "$_id.category_name",
			"month":                "$_id.month",
			"order_count":          bson.M{"$size": "$order_ids"},
			"items_sold":           1,
			"total_revenue":        1,
			"unique_products_sold": bson.M{"$size": "$product_ids"},
		}}},
		{{Key: "$sort", Value: bson.M{"month": 1, "total_revenue": -1}}},
	}

	cursor, err := db.Collection("order_items").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var rows []struct {
		CategoryID         int32   `bson:"category_id"`
		CategoryName       string  `bson:"category_name"`
		Month              int32   `bson:"month"`
		OrderCount         int64   `bson:"order_count"`
		ItemsSold          int64   `bson:"items_sold"`
		TotalRevenue       float64 `bson:"total_revenue"`
		UniqueProductsSold int64   `bson:"unique_products_sold"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, nil, err
	}

	yearStr := strconv.Itoa(year)
	var monthlyRows []*salemongo.CategoryStatsRow
	var yearlyRows []*salemongo.CategoryYearlyStatsRow
	for _, row := range rows {
		if monthly {
			monthlyRows = append(monthlyRows, &salemongo.CategoryStatsRow{
				Month:              monthAbbreviations[row.Month],
				Year:               yearStr,
				CategoryID:         row.CategoryID,
				CategoryName:       row.CategoryName,
				OrderCount:         row.OrderCount,
				ItemsSold:          row.ItemsSold,
				TotalRevenue:       row.TotalRevenue,
				UniqueProductsSold: row.UniqueProductsSold,
			})
			continue
		}
		yearlyRows = append(yearlyRows, &salemongo.CategoryYearlyStatsRow{
			Year:               yearStr,
			CategoryID:         row.CategoryID,
			CategoryName:       row.CategoryName,
			OrderCount:         row.OrderCount,
			ItemsSold:          row.ItemsSold,
			TotalRevenue:       row.TotalRevenue,
			UniqueProductsSold: row.UniqueProductsSold,
		})
	}
	return monthlyRows, yearlyRows, nil
}

// aggregateCategoryRevenue sums order revenue within [start, end] narrowed by the
// optional category/merchant filters.
func aggregateCategoryRevenue(ctx context.Context, db *mongo.Database, start, end time.Time, filter categoryStatsFilter) (int32, error) {
	match := bson.M{
		"deleted_at":       nil,
		"order.deleted_at": nil,
		"order.created_at": bson.M{"$gte": start, "$lte": end},
	}
	if filter.categoryID != nil {
		match["product.category_id"] = *filter.categoryID
	}
	if filter.merchantID != nil {
		match["order.merchant_id"] = *filter.merchantID
	}

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from": "orders", "localField": "order_id", "foreignField": "order_id", "as": "order",
		}}},
		{{Key: "$unwind", Value: "$order"}},
		{{Key: "$lookup", Value: bson.M{
			"from": "products", "localField": "product_id", "foreignField": "product_id", "as": "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$order.total_price"},
		}}},
	}

	cursor, err := db.Collection("order_items").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var res []struct {
		Total int32 `bson:"total"`
	}
	if err := cursor.All(ctx, &res); err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, nil
	}
	return res[0].Total, nil
}
