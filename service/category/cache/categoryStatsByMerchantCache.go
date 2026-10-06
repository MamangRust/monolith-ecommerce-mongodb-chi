package cache

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	categoryStatsByMerchantMonthTotalPriceCacheKey = "category:stats:bymerchant:%d:month:%d:year:%d"
	categoryStatsByMerchantYearTotalPriceCacheKey  = "category:stats:bymerchant:%d:year:%d"
	categoryStatsByMerchantMonthPriceCacheKey      = "category:stats:bymerchant:%d:month:%d"
	categoryStatsByMerchantYearPriceCacheKey       = "category:stats:bymerchant:%d:year:%d"
)

type categoryStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsByMerchantCache(store *cache.CacheStore) *categoryStatsByMerchantCache {
	return &categoryStatsByMerchantCache{store: store}
}

func (s *categoryStatsByMerchantCache) GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryRevenueRow, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data []*salemongo.CategoryRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyRevenueRow, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryYearlyRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data []*salemongo.CategoryYearlyRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryStatsRow, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryStatsRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data []*salemongo.CategoryStatsRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyStatsRow, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CategoryYearlyStatsRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data []*salemongo.CategoryYearlyStatsRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
