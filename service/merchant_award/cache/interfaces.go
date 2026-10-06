package cache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryCache interface {
	GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationListRow, *int, bool)
	SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantCertificationListRow, total *int)

	GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationActiveRow, *int, bool)
	SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantCertificationActiveRow, total *int)

	GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationTrashedRow, *int, bool)
	SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantCertificationTrashedRow, total *int)

	GetCachedMerchantAward(ctx context.Context, id int) (*salemongo.MerchantCertificationRow, bool)
	SetCachedMerchantAward(ctx context.Context, data *salemongo.MerchantCertificationRow)
}

type MerchantAwardCommandCache interface {
	DeleteMerchantAwardCache(ctx context.Context, id int)
}
