package cache

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	categoryStatsMonthTotalPriceCacheKey = "category:stats:month:%d:year:%d"
	categoryStatsYearTotalPriceCacheKey  = "category:stats:year:%d"
	categoryStatsMonthPriceCacheKey      = "category:stats:month:%d"
	categoryStatsYearPriceCacheKey       = "category:stats:year:%d"
)

type categoryStatsCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsCache(store *cache.CacheStore) *categoryStatsCache {
	return &categoryStatsCache{store: store}
}

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryRevenueRow, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*salemongo.CategoryRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRevenueRow, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryYearlyRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*salemongo.CategoryYearlyRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryStatsRow, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryStatsRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(ctx context.Context, year int, data []*salemongo.CategoryStatsRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryYearlyStatsRow, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryYearlyStatsRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(ctx context.Context, year int, data []*salemongo.CategoryYearlyStatsRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
