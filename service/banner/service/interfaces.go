package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, error)
	FindByID(ctx context.Context, bannerID int) (*repository.BannerResult, error)
}

type BannerCommandService interface {
	Create(ctx context.Context, req *requests.CreateBannerRequest) (*salemongo.Banner, error)
	Update(ctx context.Context, req *requests.UpdateBannerRequest) (*salemongo.Banner, error)

	Trash(ctx context.Context, bannerID int) (*salemongo.Banner, error)
	Restore(ctx context.Context, bannerID int) (*salemongo.Banner, error)
	DeletePermanent(ctx context.Context, bannerID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
