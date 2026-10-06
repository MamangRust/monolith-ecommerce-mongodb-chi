package cache

import (
	"context"
	"fmt"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	merchantPolicyAllCacheKey     = "merchant_policy:all:page:%d:pageSize:%d:search:%s"
	merchantPolicyByIdCacheKey    = "merchant_policy:id:%d"
	merchantPolicyActiveCacheKey  = "merchant_policy:active:page:%d:pageSize:%d:search:%s"
	merchantPolicyTrashedCacheKey = "merchant_policy:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantPolicyCacheResponseDB struct {
	Data         []*salemongo.MerchantPolicyListRow `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantPolicyActiveCacheResponseDB struct {
	Data         []*salemongo.MerchantPolicyActiveRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantPolicyTrashedCacheResponseDB struct {
	Data         []*salemongo.MerchantPolicyTrashedRow `json:"data"`
	TotalRecords *int                                  `json:"total_records"`
}

type merchantPoliciesQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantPoliciesQueryCache(store *cache.CacheStore) MerchantPoliciesQueryCache {
	return &merchantPoliciesQueryCache{
		store: store,
	}
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyListRow, *int, bool) {
	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantPolicyCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantPolicyListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantPolicyListRow{}
	}

	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantPolicyCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantPolicyActiveCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantPolicyActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantPolicyActiveRow{}
	}

	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantPolicyActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantPolicyTrashedCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantPolicyTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantPolicyTrashedRow{}
	}

	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantPolicyTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicy(ctx context.Context, id int) (*salemongo.MerchantPolicyRow, bool) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)

	result, found := cache.GetFromCache[salemongo.MerchantPolicyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicy(ctx context.Context, data *salemongo.MerchantPolicyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantPolicyByIdCacheKey, data.MerchantPolicyID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
