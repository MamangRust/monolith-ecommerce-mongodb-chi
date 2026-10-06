package cache

import (
	"context"
	"fmt"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	monthlyTotalRevenueCacheKeyByMerchant = "order:monthly:totalRevenue:merchant:%d:month:%d:year:%d"
	yearlyTotalRevenueCacheKeyByMerchant  = "order:yearly:totalRevenue:merchant:%d:year:%d"
	monthlyOrderCacheKeyByMerchant        = "order:monthly:order:merchant:%d:year:%d"
	yearlyOrderCacheKeyByMerchant         = "order:yearly:order:merchant:%d:year:%d"
)

type orderStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewOrderStatsByMerchantCache(store *cache.CacheStore) *orderStatsByMerchantCache {
	return &orderStatsByMerchantCache{store: store}
}

func (s *orderStatsByMerchantCache) GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.OrderMonthlyRevenueRow, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.OrderMonthlyRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data []*salemongo.OrderMonthlyRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.OrderYearlyRevenueRow, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.OrderYearlyRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data []*salemongo.OrderYearlyRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.OrderMonthlyOrderRow, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.OrderMonthlyOrderRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data []*salemongo.OrderMonthlyOrderRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.OrderYearlyOrderRow, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.OrderYearlyOrderRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data []*salemongo.OrderYearlyOrderRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
