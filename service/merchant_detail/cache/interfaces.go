package cache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailListRow, *int, bool)
	SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantDetailListRow, total *int)

	GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailActiveRow, *int, bool)
	SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantDetailActiveRow, total *int)

	GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailTrashedRow, *int, bool)
	SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantDetailTrashedRow, total *int)

	GetCachedMerchantDetail(ctx context.Context, id int) (*salemongo.MerchantDetailWithJoinsRow, bool)
	SetCachedMerchantDetail(ctx context.Context, data *salemongo.MerchantDetailWithJoinsRow)
}

type MerchantDetailCommandCache interface {
	DeleteMerchantDetailCache(ctx context.Context, id int)
}
