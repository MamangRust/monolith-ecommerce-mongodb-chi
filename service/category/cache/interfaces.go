package cache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryCache(ctx context.Context, id int) (*salemongo.Category, bool)
	SetCachedCategoryCache(ctx context.Context, data *salemongo.Category)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryRevenueRow, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*salemongo.CategoryRevenueRow)

	GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRevenueRow, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*salemongo.CategoryYearlyRevenueRow)

	GetCachedMonthPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryStatsRow, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, data []*salemongo.CategoryStatsRow)

	GetCachedYearPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryYearlyStatsRow, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, data []*salemongo.CategoryYearlyStatsRow)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryRevenueRow, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*salemongo.CategoryRevenueRow)

	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyRevenueRow, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*salemongo.CategoryYearlyRevenueRow)

	GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryStatsRow, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*salemongo.CategoryStatsRow)

	GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyStatsRow, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*salemongo.CategoryYearlyStatsRow)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryRevenueRow, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data []*salemongo.CategoryRevenueRow)

	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyRevenueRow, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data []*salemongo.CategoryYearlyRevenueRow)

	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryStatsRow, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data []*salemongo.CategoryStatsRow)

	GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyStatsRow, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data []*salemongo.CategoryYearlyStatsRow)
}
