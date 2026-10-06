package cache

import (
	"context"
	"fmt"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	transactionMonthlyAmountSuccessKey = "transaction:month:amount:success:month:%d:year:%d"
	transactionMonthlyAmountFailedKey  = "transaction:month:amount:failed:month:%d:year:%d"
	transactionYearlyAmountSuccessKey  = "transaction:year:amount:success:year:%d"
	transactionYearlyAmountFailedKey   = "transaction:year:amount:failed:year:%d"
	transactionMonthlyMethodSuccessKey = "transaction:month:method:success:month:%d:year:%d"
	transactionMonthlyMethodFailedKey  = "transaction:month:method:failed:month:%d:year:%d"
	transactionYearlyMethodSuccessKey  = "transaction:year:method:success:year:%d"
	transactionYearlyMethodFailedKey   = "transaction:year:method:failed:year:%d"
)

type transactionStatsCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsCache(store *cache.CacheStore) *transactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthlyAmountSuccessKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyAmountSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyAmountSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearlyAmountSuccessKey, year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyAmountSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyAmountSuccess(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyAmountSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthlyAmountFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyAmountFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyAmountFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, bool) {
	key := fmt.Sprintf(transactionYearlyAmountFailedKey, year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyAmountFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyAmountFailed(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyAmountFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthlyMethodSuccessKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyMethodSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyMethodSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearlyMethodSuccessKey, year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyMethodSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyMethodSuccess(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyMethodSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthlyMethodFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyMethodFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyMethodFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, bool) {
	key := fmt.Sprintf(transactionYearlyMethodFailedKey, year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyMethodFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyMethodFailed(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyMethodFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
