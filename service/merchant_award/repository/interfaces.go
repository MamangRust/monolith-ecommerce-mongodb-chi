package repository

import (
	"context"

	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationListRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.MerchantCertificationRow, error)
}

type MerchantAwardCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantCertificationOrAwardRequest) (*salemongo.MerchantCertificationRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*salemongo.MerchantCertificationRow, error)
	Trash(ctx context.Context, merchant_id int) (*salemongo.MerchantCertificationsAndAward, error)
	Restore(ctx context.Context, merchant_id int) (*salemongo.MerchantCertificationsAndAward, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository
