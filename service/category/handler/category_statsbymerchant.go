package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

type categoryStatsByMerchantHandler struct {
	pbcategories.UnimplementedCategoryStatsByMerchantServiceServer
	categoryStatsByMerchant service.CategoryStatsByMerchantService
	logger                  logger.LoggerInterface
}

func NewCategoryStatsByMerchantHandler(categoryStatsByMerchant service.CategoryStatsByMerchantService, logger logger.LoggerInterface) CategoryStatsByMerchantHandler {
	return &categoryStatsByMerchantHandler{
		categoryStatsByMerchant: categoryStatsByMerchant,
		logger:                  logger,
	}
}

func (h *categoryStatsByMerchantHandler) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pbcategories.FindYearMonthTotalPriceByMerchant) (*pbcategories.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
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

func (h *categoryStatsByMerchantHandler) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pbcategories.FindYearTotalPriceByMerchant) (*pbcategories.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(ctx, &reqService)
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

func (h *categoryStatsByMerchantHandler) FindMonthPriceByMerchant(ctx context.Context, req *pbcategories.FindYearCategoryByMerchant) (*pbcategories.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategories.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategories.CategoryMonthPriceResponse)
	}

	return &pbcategories.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByMerchantHandler) FindYearPriceByMerchant(ctx context.Context, req *pbcategories.FindYearCategoryByMerchant) (*pbcategories.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategories.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategories.CategoryYearPriceResponse)
	}

	return &pbcategories.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}
