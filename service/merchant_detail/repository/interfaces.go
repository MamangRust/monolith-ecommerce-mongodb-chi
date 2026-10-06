package repository

import (
	"context"

	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailListRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.MerchantDetailWithJoinsRow, error)
}

type MerchantDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*salemongo.MerchantDetailRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*salemongo.MerchantDetailRow, error)
	Trash(ctx context.Context, merchant_id int) (*salemongo.MerchantDetail, error)
	Restore(ctx context.Context, merchant_id int) (*salemongo.MerchantDetail, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantSocialRequest) (*salemongo.MerchantSocialMediaLinkRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantSocialRequest) (*salemongo.MerchantSocialMediaLinkRow, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	Trash(ctx context.Context, id int) (bool, error)
	Restore(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository
