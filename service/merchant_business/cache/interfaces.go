package cache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryCache interface {
	GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoListRow, *int, bool)
	SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantBusinessInfoListRow, total *int)

	GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoActiveRow, *int, bool)
	SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantBusinessInfoActiveRow, total *int)

	GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoTrashedRow, *int, bool)
	SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantBusinessInfoTrashedRow, total *int)

	GetCachedMerchantBusiness(ctx context.Context, id int) (*salemongo.MerchantBusinessInfoRow, bool)
	SetCachedMerchantBusiness(ctx context.Context, data *salemongo.MerchantBusinessInfoRow)
}

type MerchantBusinessCommandCache interface {
	DeleteMerchantBusinessCache(ctx context.Context, id int)
}
