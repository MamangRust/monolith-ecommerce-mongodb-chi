package repository

import (
	"context"

	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyListRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantPolicyTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.MerchantPolicyRow, error)
}

type MerchantPoliciesCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*salemongo.MerchantPolicyRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*salemongo.MerchantPolicyRow, error)
	Trash(ctx context.Context, merchant_id int) (*salemongo.MerchantPolicy, error)
	Restore(ctx context.Context, merchant_id int) (*salemongo.MerchantPolicy, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository
