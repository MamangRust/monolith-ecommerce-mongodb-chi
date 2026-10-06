package cache

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	categoryStatsByIdMonthTotalPriceCacheKey = "category:stats:byid:%d:month:%d:year:%d"
	categoryStatsByIdYearTotalPriceCacheKey  = "category:stats:byid:%d:year:%d"
	categoryStatsByIdMonthPriceCacheKey      = "category:stats:byid:%d:month:%d"
	categoryStatsByIdYearPriceCacheKey       = "category:stats:byid:%d:year:%d"
)

type categoryStatsByIdCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsByIdCache(store *cache.CacheStore) *categoryStatsByIdCache {
	return &categoryStatsByIdCache{store: store}
}

func (s *categoryStatsByIdCache) GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryRevenueRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*salemongo.CategoryRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyRevenueRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryYearlyRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*salemongo.CategoryYearlyRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryStatsRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryStatsRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*salemongo.CategoryStatsRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyStatsRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryYearlyStatsRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*salemongo.CategoryYearlyStatsRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
