package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryCommandHandler struct {
	pbcategories.UnimplementedCategoryCommandServiceServer
	service service.CategoryCommandService
	logger  logger.LoggerInterface
}

func NewCategoryCommandHandler(service service.CategoryCommandService, logger logger.LoggerInterface) pbcategories.CategoryCommandServiceServer {
	return &categoryCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (h *categoryCommandHandler) Create(ctx context.Context, request *pbcategories.CreateCategoryRequest) (*pbcategories.ApiResponseCategory, error) {
	slug := request.GetSlugCategory()
	req := &requests.CreateCategoryRequest{
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		SlugCategory:  &slug,
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := h.service.Create(ctx, req)
	if err != nil {
		return nil, category_errors.ErrGrpcCreateCategory
	}

	return &pbcategories.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategories.CategoryResponse),
	}, nil
}

func (h *categoryCommandHandler) Update(ctx context.Context, request *pbcategories.UpdateCategoryRequest) (*pbcategories.ApiResponseCategory, error) {
	id := int(request.GetCategoryId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	slug := request.GetSlugCategory()
	req := &requests.UpdateCategoryRequest{
		CategoryID:    &id,
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		SlugCategory:  &slug,
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := h.service.Update(ctx, req)
	if err != nil {
		return nil, category_errors.ErrGrpcUpdateCategory
	}

	return &pbcategories.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategories.CategoryResponse),
	}, nil
}

func (h *categoryCommandHandler) TrashedCategory(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.Trash(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategories.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    mapToCategoryDeleteAtResponse(category),
	}, nil
}

func (h *categoryCommandHandler) RestoreCategory(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.Restore(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategories.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    mapToCategoryDeleteAtResponse(category),
	}, nil
}

func (h *categoryCommandHandler) DeleteCategoryPermanent(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	_, err := h.service.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbcategories.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (h *categoryCommandHandler) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pbcategories.ApiResponseCategoryAll, error) {
	_, err := h.service.RestoreAll(ctx)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategories.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restored all categories",
	}, nil
}

func (h *categoryCommandHandler) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcategories.ApiResponseCategoryAll, error) {
	_, err := h.service.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbcategories.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully deleted all categories permanently",
	}, nil
}
