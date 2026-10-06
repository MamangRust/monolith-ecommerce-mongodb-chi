package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

type categoryStatsByIdHandler struct {
	pbcategories.UnimplementedCategoryStatsByIdServiceServer
	categoryStatsById service.CategoryStatsByIdService
	logger            logger.LoggerInterface
}

func NewCategoryStatsByIdHandler(categoryStatsById service.CategoryStatsByIdService, logger logger.LoggerInterface) CategoryStatsByIdHandler {
	return &categoryStatsByIdHandler{
		categoryStatsById: categoryStatsById,
		logger:            logger,
	}
}

func (h *categoryStatsByIdHandler) FindMonthlyTotalPricesById(ctx context.Context, req *pbcategories.FindYearMonthTotalPriceById) (*pbcategories.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategories.CategoriesMonthlyTotalPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategories.CategoriesMonthlyTotalPriceResponse)
	}

	return &pbcategories.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindYearlyTotalPricesById(ctx context.Context, req *pbcategories.FindYearTotalPriceById) (*pbcategories.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategories.CategoriesYearlyTotalPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategories.CategoriesYearlyTotalPriceResponse)
	}

	return &pbcategories.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindMonthPriceById(ctx context.Context, req *pbcategories.FindYearCategoryById) (*pbcategories.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategories.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategories.CategoryMonthPriceResponse)
	}

	return &pbcategories.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Category monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindYearPriceById(ctx context.Context, req *pbcategories.FindYearCategoryById) (*pbcategories.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindYearPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategories.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategories.CategoryYearPriceResponse)
	}

	return &pbcategories.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Category yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}
