package cache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryCache interface {
	GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyListRow, *int, bool)
	SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantPolicyListRow, total *int)

	GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyActiveRow, *int, bool)
	SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantPolicyActiveRow, total *int)

	GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyTrashedRow, *int, bool)
	SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*salemongo.MerchantPolicyTrashedRow, total *int)

	GetCachedMerchantPolicy(ctx context.Context, id int) (*salemongo.MerchantPolicyRow, bool)
	SetCachedMerchantPolicy(ctx context.Context, data *salemongo.MerchantPolicyRow)
}

type MerchantPoliciesCommandCache interface {
	DeleteMerchantPolicyCache(ctx context.Context, merchantID int)
}
