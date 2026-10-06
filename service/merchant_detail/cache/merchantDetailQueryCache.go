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
	merchantDetailAllCacheKey     = "merchant_detail:all:page:%d:pageSize:%d:search:%s"
	merchantDetailByIdCacheKey    = "merchant_detail:id:%d"
	merchantDetailActiveCacheKey  = "merchant_detail:active:page:%d:pageSize:%d:search:%s"
	merchantDetailTrashedCacheKey = "merchant_detail:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantDetailCacheResponseDB struct {
	Data         []*salemongo.MerchantDetailListRow `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantDetailActiveCacheResponseDB struct {
	Data         []*salemongo.MerchantDetailActiveRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantDetailTrashedCacheResponseDB struct {
	Data         []*salemongo.MerchantDetailTrashedRow `json:"data"`
	TotalRecords *int                                  `json:"total_records"`
}

type merchantDetailQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDetailQueryCache(store *cache.CacheStore) *merchantDetailQueryCache {
	return &merchantDetailQueryCache{
		store: store,
	}
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailListRow, *int, bool) {
	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDetailCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantDetailListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantDetailListRow{}
	}

	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDetailActiveCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantDetailActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantDetailActiveRow{}
	}

	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDetailTrashedCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantDetailTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.MerchantDetailTrashedRow{}
	}

	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetail(ctx context.Context, id int) (*salemongo.MerchantDetailWithJoinsRow, bool) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)

	result, found := cache.GetFromCache[salemongo.MerchantDetailWithJoinsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetail(ctx context.Context, data *salemongo.MerchantDetailWithJoinsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailByIdCacheKey, data.MerchantDetailID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
