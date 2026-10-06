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
	merchantBusinessAllCacheKey     = "merchant_business:all:page:%d:pageSize:%d:search:%s"
	merchantBusinessByIdCacheKey    = "merchant_business:id:%d"
	merchantBusinessActiveCacheKey  = "merchant_business:active:page:%d:pageSize:%d:search:%s"
	merchantBusinessTrashedCacheKey = "merchant_business:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantBusinessCacheResponseDB struct {
	Data         []*salemongo.MerchantBusinessInfoListRow `json:"data"`
	TotalRecords *int                                     `json:"total_records"`
}

type merchantBusinessActiveCacheResponseDB struct {
	Data         []*salemongo.MerchantBusinessInfoActiveRow `json:"data"`
	TotalRecords *int                                       `json:"total_records"`
}

type merchantBusinessTrashedCacheResponseDB struct {
	Data         []*salemongo.MerchantBusinessInfoTrashedRow `json:"data"`
	TotalRecords *int                                        `json:"total_records"`
}

type merchantBusinessQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantBusinessQueryCache(store *cache.CacheStore) *merchantBusinessQueryCache {
	return &merchantBusinessQueryCache{
		store: store,
	}
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoListRow, *int, bool) {
	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantBusinessCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantBusinessInfoListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantBusinessInfoListRow{}
	}

	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantBusinessActiveCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantBusinessInfoActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantBusinessInfoActiveRow{}
	}

	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantBusinessTrashedCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantBusinessInfoTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantBusinessInfoTrashedRow{}
	}

	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusiness(ctx context.Context, id int) (*salemongo.MerchantBusinessInfoRow, bool) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)

	result, found := cache.GetFromCache[salemongo.MerchantBusinessInfoRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusiness(ctx context.Context, data *salemongo.MerchantBusinessInfoRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantBusinessByIdCacheKey, data.MerchantBusinessInfoID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
