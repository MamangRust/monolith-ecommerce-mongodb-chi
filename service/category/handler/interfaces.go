package handler

import (
	"context"

	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryQueryHandler interface {
	pbcategories.CategoryQueryServiceServer
}

type CategoryCommandHandler interface {
	pbcategories.CategoryCommandServiceServer
}

type CategoryStatsHandler interface {
	pbcategories.CategoryStatsServiceServer
}

type CategoryStatsByIdHandler interface {
	pbcategories.CategoryStatsByIdServiceServer
}

type CategoryStatsByMerchantHandler interface {
	pbcategories.CategoryStatsByMerchantServiceServer
}

type CategoryHandleGrpc interface {
	FindAll(ctx context.Context, request *pbcategories.FindAllCategoryRequest) (*pbcategories.ApiResponsePaginationCategory, error)
	FindById(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategory, error)
	FindByActive(ctx context.Context, request *pbcategories.FindAllCategoryRequest) (*pbcategories.ApiResponsePaginationCategoryDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pbcategories.FindAllCategoryRequest) (*pbcategories.ApiResponsePaginationCategoryDeleteAt, error)
	Create(ctx context.Context, request *pbcategories.CreateCategoryRequest) (*pbcategories.ApiResponseCategory, error)
	Update(ctx context.Context, request *pbcategories.UpdateCategoryRequest) (*pbcategories.ApiResponseCategory, error)
	TrashedCategory(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDeleteAt, error)
	RestoreCategory(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDeleteAt, error)
	DeleteCategoryPermanent(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDelete, error)
	RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pbcategories.ApiResponseCategoryAll, error)
	DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcategories.ApiResponseCategoryAll, error)
}
