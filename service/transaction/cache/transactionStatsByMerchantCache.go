package cache

import (
	"context"
	"fmt"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	transactonMonthlyAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthlyAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"
	transactonYearlyAmountSuccessByMerchantKey  = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearlyAmountFailedByMerchantKey   = "transaction:year:amount:failed:merchant:%d:year:%d"
	transactonMonthlyMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthlyMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"
	transactonYearlyMethodSuccessByMerchantKey  = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearlyMethodFailedByMerchantKey   = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) *transactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthlyAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*salemongo.TransactionMonthlyAmountSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearlyAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyAmountSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*salemongo.TransactionYearlyAmountSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthlyAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyAmountFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*salemongo.TransactionMonthlyAmountFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearlyAmountFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyAmountFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*salemongo.TransactionYearlyAmountFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyAmountFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthlyMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*salemongo.TransactionMonthlyMethodSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearlyMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyMethodSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*salemongo.TransactionYearlyMethodSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthlyMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyMethodFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*salemongo.TransactionMonthlyMethodFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearlyMethodFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyMethodFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*salemongo.TransactionYearlyMethodFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
