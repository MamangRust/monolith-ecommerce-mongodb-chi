package repository

import (
	"context"

	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoListRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.MerchantBusinessInfoRow, error)
}

type MerchantBusinessCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantBusinessInformationRequest) (*salemongo.MerchantBusinessInfoRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*salemongo.MerchantBusinessInfoRow, error)
	Trash(ctx context.Context, merchant_id int) (*salemongo.MerchantBusinessInformation, error)
	Restore(ctx context.Context, merchant_id int) (*salemongo.MerchantBusinessInformation, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository
