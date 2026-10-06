package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryResult struct {
	CategoryID    int32
	Name          string
	Description   *string
	SlugCategory  *string
	ImageCategory *string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	TotalCount    int64
}

type CategoryQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error)
	FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error)
	FindByID(ctx context.Context, category_id int) (*salemongo.Category, error)
	FindByIDTrashed(ctx context.Context, category_id int) (*salemongo.Category, error)
}

type CategoryCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateCategoryRequest) (*salemongo.Category, error)
	Update(ctx context.Context, request *requests.UpdateCategoryRequest) (*salemongo.Category, error)
	Trash(ctx context.Context, category_id int) (*salemongo.Category, error)
	Restore(ctx context.Context, category_id int) (*salemongo.Category, error)
	DeletePermanent(ctx context.Context, category_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryRevenueRow, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRevenueRow, error)
	GetMonthPrice(ctx context.Context, year int) ([]*salemongo.CategoryStatsRow, error)
	GetYearPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyStatsRow, error)
}

type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryRevenueRow, error)
	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyRevenueRow, error)
	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryStatsRow, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyStatsRow, error)
}

type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryRevenueRow, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyRevenueRow, error)
	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryStatsRow, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyStatsRow, error)
}
