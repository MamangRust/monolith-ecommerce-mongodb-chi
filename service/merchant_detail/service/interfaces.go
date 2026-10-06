package service

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailListRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailTrashedRow, *int, error)
	FindByID(ctx context.Context, user_id int) (*salemongo.MerchantDetailWithJoinsRow, error)
}

type MerchantDetailCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*salemongo.MerchantDetailRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*salemongo.MerchantDetailRow, error)
	Trash(ctx context.Context, merchant_id int) (*salemongo.MerchantDetail, error)
	Restore(ctx context.Context, merchant_id int) (*salemongo.MerchantDetail, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandService interface {
	Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*salemongo.MerchantSocialMediaLinkRow, error)
	Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*salemongo.MerchantSocialMediaLinkRow, error)
	Trash(ctx context.Context, socialID int) (bool, error)
	Restore(ctx context.Context, socialID int) (bool, error)
	DeletePermanent(ctx context.Context, socialID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
