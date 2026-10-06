package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryStatsByIdRepository struct {
	db *mongo.Database
}

func NewCategoryStatsByIdRepository(db *mongo.Database) *categoryStatsByIdRepository {
	return &categoryStatsByIdRepository{db: db}
}

func (r *categoryStatsByIdRepository) GetMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryRevenueRow, error) {
	filter := categoryStatsFilter{categoryID: int32Ptr(req.CategoryID)}

	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	currentRevenue, err := aggregateCategoryRevenue(ctx, r.db, currentMonthStart, currentMonthEnd, filter)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceById.WithInternal(err)
	}
	prevRevenue, err := aggregateCategoryRevenue(ctx, r.db, prevMonthStart, prevMonthEnd, filter)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceById.WithInternal(err)
	}

	return []*salemongo.CategoryRevenueRow{
		{Year: formatYear(currentMonthStart), Month: currentMonthStart.Month().String(), TotalRevenue: float64(currentRevenue)},
		{Year: formatYear(prevMonthStart), Month: prevMonthStart.Month().String(), TotalRevenue: float64(prevRevenue)},
	}, nil
}

func (r *categoryStatsByIdRepository) GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyRevenueRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year, 12, 31, 23, 59, 59, 0, time.UTC)

	totalRevenue, err := aggregateCategoryRevenue(ctx, r.db, yearStart, yearEnd, categoryStatsFilter{categoryID: int32Ptr(req.CategoryID)})
	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesById.WithInternal(err)
	}

	return []*salemongo.CategoryYearlyRevenueRow{
		{Year: formatYear(yearStart), TotalRevenue: float64(totalRevenue)},
	}, nil
}

func (r *categoryStatsByIdRepository) GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryStatsRow, error) {
	monthlyRows, _, err := aggregateCategoryStats(ctx, r.db, req.Year, true, categoryStatsFilter{categoryID: int32Ptr(req.CategoryID)})
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceById.WithInternal(err)
	}
	return monthlyRows, nil
}

func (r *categoryStatsByIdRepository) GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyStatsRow, error) {
	_, yearlyRows, err := aggregateCategoryStats(ctx, r.db, req.Year, false, categoryStatsFilter{categoryID: int32Ptr(req.CategoryID)})
	if err != nil {
		return nil, category_errors.ErrGetYearPriceById.WithInternal(err)
	}
	return yearlyRows, nil
}

func int32Ptr(v int) *int32 {
	converted := int32(v)
	return &converted
}
