package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryStatsByMerchantRepository struct {
	db *mongo.Database
}

func NewCategoryStatsByMerchantRepository(db *mongo.Database) *categoryStatsByMerchantRepository {
	return &categoryStatsByMerchantRepository{db: db}
}

func (r *categoryStatsByMerchantRepository) GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryRevenueRow, error) {
	filter := categoryStatsFilter{merchantID: int32Ptr(req.MerchantID)}

	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	currentRevenue, err := aggregateCategoryRevenue(ctx, r.db, currentMonthStart, currentMonthEnd, filter)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceByMerchant.WithInternal(err)
	}
	prevRevenue, err := aggregateCategoryRevenue(ctx, r.db, prevMonthStart, prevMonthEnd, filter)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceByMerchant.WithInternal(err)
	}

	return []*salemongo.CategoryRevenueRow{
		{Year: formatYear(currentMonthStart), Month: currentMonthStart.Month().String(), TotalRevenue: float64(currentRevenue)},
		{Year: formatYear(prevMonthStart), Month: prevMonthStart.Month().String(), TotalRevenue: float64(prevRevenue)},
	}, nil
}

func (r *categoryStatsByMerchantRepository) GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyRevenueRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(req.Year, 12, 31, 23, 59, 59, 0, time.UTC)

	totalRevenue, err := aggregateCategoryRevenue(ctx, r.db, yearStart, yearEnd, categoryStatsFilter{merchantID: int32Ptr(req.MerchantID)})
	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesByMerchant.WithInternal(err)
	}

	return []*salemongo.CategoryYearlyRevenueRow{
		{Year: formatYear(yearStart), TotalRevenue: float64(totalRevenue)},
	}, nil
}

func (r *categoryStatsByMerchantRepository) GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryStatsRow, error) {
	monthlyRows, _, err := aggregateCategoryStats(ctx, r.db, req.Year, true, categoryStatsFilter{merchantID: int32Ptr(req.MerchantID)})
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceByMerchant.WithInternal(err)
	}
	return monthlyRows, nil
}

func (r *categoryStatsByMerchantRepository) GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyStatsRow, error) {
	_, yearlyRows, err := aggregateCategoryStats(ctx, r.db, req.Year, false, categoryStatsFilter{merchantID: int32Ptr(req.MerchantID)})
	if err != nil {
		return nil, category_errors.ErrGetYearPriceByMerchant.WithInternal(err)
	}
	return yearlyRows, nil
}
