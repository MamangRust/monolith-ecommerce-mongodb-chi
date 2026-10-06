package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindByID(ctx context.Context, categoryID int) (*salemongo.Category, error)
}

type CategoryCommandService interface {
	Create(ctx context.Context, req *requests.CreateCategoryRequest) (*salemongo.Category, error)
	Update(ctx context.Context, req *requests.UpdateCategoryRequest) (*salemongo.Category, error)
	Trash(ctx context.Context, categoryID int) (*salemongo.Category, error)
	Restore(ctx context.Context, categoryID int) (*salemongo.Category, error)
	DeletePermanent(ctx context.Context, categoryID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type CategoryStatsService interface {
	FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryRevenueRow, error)
	FindYearlyTotalPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRevenueRow, error)
	FindMonthPrice(ctx context.Context, year int) ([]*salemongo.CategoryStatsRow, error)
	FindYearPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyStatsRow, error)
}

type CategoryStatsByIdService interface {
	FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryRevenueRow, error)
	FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyRevenueRow, error)
	FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryStatsRow, error)
	FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyStatsRow, error)
}

type CategoryStatsByMerchantService interface {
	FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryRevenueRow, error)
	FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyRevenueRow, error)
	FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryStatsRow, error)
	FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyStatsRow, error)
}
