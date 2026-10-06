package cache

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	orderItemAllCacheKey     = "order_item:all:page:%d:pageSize:%d:search:%s"
	orderItemActiveCacheKey  = "order_item:active:page:%d:pageSize:%d:search:%s"
	orderItemTrashedCacheKey = "order_item:trashed:page:%d:pageSize:%d:search:%s"
	orderItemByOrderCacheKey = "order_item:order:%d"

	ttlDefault = 5 * time.Minute
)

type orderItemCacheResponseDB struct {
	Data         []*salemongo.OrderItemListRow `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type orderItemActiveCacheResponseDB struct {
	Data         []*salemongo.OrderItemActiveRow `json:"data"`
	TotalRecords *int                            `json:"total_records"`
}

type orderItemTrashedCacheResponseDB struct {
	Data         []*salemongo.OrderItemTrashedRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type orderItemQueryCache struct {
	store *cache.CacheStore
}

func NewOrderItemQueryCache(store *cache.CacheStore) *orderItemQueryCache {
	return &orderItemQueryCache{store: store}
}

func (o *orderItemQueryCache) GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemListRow, *int, bool) {
	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[orderItemCacheResponseDB](ctx, o.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (o *orderItemQueryCache) SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*salemongo.OrderItemListRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.OrderItemListRow{}
	}

	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &orderItemCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, o.store, key, payload, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemActiveRow, *int, bool) {
	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[orderItemActiveCacheResponseDB](ctx, o.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (o *orderItemQueryCache) SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*salemongo.OrderItemActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.OrderItemActiveRow{}
	}

	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, o.store, key, payload, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemTrashedRow, *int, bool) {
	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderItemTrashedCacheResponseDB](ctx, o.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (o *orderItemQueryCache) SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*salemongo.OrderItemTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*salemongo.OrderItemTrashedRow{}
	}

	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, o.store, key, payload, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItems(ctx context.Context, orderID int) ([]*salemongo.OrderItemRow, bool) {
	key := fmt.Sprintf(orderItemByOrderCacheKey, orderID)
	result, found := cache.GetFromCache[[]*salemongo.OrderItemRow](ctx, o.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (o *orderItemQueryCache) SetCachedOrderItems(ctx context.Context, data []*salemongo.OrderItemRow) {
	if len(data) == 0 {
		return
	}

	key := fmt.Sprintf(orderItemByOrderCacheKey, data[0].OrderID)
	cache.SetToCache(ctx, o.store, key, &data, ttlDefault)
}
