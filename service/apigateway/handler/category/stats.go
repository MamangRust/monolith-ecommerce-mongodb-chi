package categoryhandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	category_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/category"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/httpx"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/category"
	"github.com/go-chi/chi/v5"
)

type categoryStatsHandlerApi struct {
	statsClient           pb.CategoryStatsServiceClient
	statsByIdClient       pb.CategoryStatsByIdServiceClient
	statsByMerchantClient pb.CategoryStatsByMerchantServiceClient
	logger                logger.LoggerInterface
	mapper                apimapper.CategoryStatsResponseMapper
	cache                 category_cache.CategoryMencache
	errors                apierror.ApiHandler
}

type categoryStatsHandleDeps struct {
	statsClient           pb.CategoryStatsServiceClient
	statsByIdClient       pb.CategoryStatsByIdServiceClient
	statsByMerchantClient pb.CategoryStatsByMerchantServiceClient
	router                chi.Router
	logger                logger.LoggerInterface
	mapper                apimapper.CategoryStatsResponseMapper
	cache                 category_cache.CategoryMencache
	apiHandler            apierror.ApiHandler
}

func NewCategoryStatsHandleApi(params *categoryStatsHandleDeps) *categoryStatsHandlerApi {
	handler := &categoryStatsHandlerApi{
		statsClient:           params.statsClient,
		statsByIdClient:       params.statsByIdClient,
		statsByMerchantClient: params.statsByMerchantClient,
		logger:                params.logger,
		mapper:                params.mapper,
		cache:                 params.cache,
		errors:                params.apiHandler,
	}

	params.router.Route("/api/category-stats", func(routerCategory chi.Router) {

		// Stats
		routerCategory.Get("/monthly-total-pricing", httpx.Handler(handler.FindMonthTotalPrice))
		routerCategory.Get("/yearly-total-pricing", httpx.Handler(handler.FindYearTotalPrice))
		routerCategory.Get("/merchant/monthly-total-pricing", httpx.Handler(handler.FindMonthTotalPriceByMerchant))
		routerCategory.Get("/merchant/yearly-total-pricing", httpx.Handler(handler.FindYearTotalPriceByMerchant))
		routerCategory.Get("/mycategory/monthly-total-pricing", httpx.Handler(handler.FindMonthTotalPriceById))
		routerCategory.Get("/mycategory/yearly-total-pricing", httpx.Handler(handler.FindYearTotalPriceById))

		routerCategory.Get("/monthly-pricing", httpx.Handler(handler.FindMonthPrice))
		routerCategory.Get("/yearly-pricing", httpx.Handler(handler.FindYearPrice))
		routerCategory.Get("/merchant/monthly-pricing", httpx.Handler(handler.FindMonthPriceByMerchant))
		routerCategory.Get("/merchant/yearly-pricing", httpx.Handler(handler.FindYearPriceByMerchant))
		routerCategory.Get("/mycategory/monthly-pricing", httpx.Handler(handler.FindMonthPriceById))
		routerCategory.Get("/mycategory/yearly-pricing", httpx.Handler(handler.FindYearPriceById))

	})
	return handler
}

// @Security Bearer
// @Summary Get monthly total pricing for all categories
// @Tags Category Stats
// @Description Retrieve monthly total revenue for all categories
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice "Monthly total pricing"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/monthly-total-pricing [get]
func (h *categoryStatsHandlerApi) FindMonthTotalPrice(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	ctx := r.Context()
	req := &requests.MonthTotalPrice{Year: year, Month: month}

	if cached, found := h.cache.GetCachedMonthTotalPriceCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsClient.FindMonthlyTotalPrices(ctx, &pb.FindYearMonthTotalPrices{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryMonthlyTotalPrice(res)
	h.cache.SetCachedMonthTotalPriceCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get yearly total pricing for all categories
// @Tags Category Stats
// @Description Retrieve yearly total revenue for all categories
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice "Yearly total pricing"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/yearly-total-pricing [get]
func (h *categoryStatsHandlerApi) FindYearTotalPrice(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))

	ctx := r.Context()
	if cached, found := h.cache.GetCachedYearTotalPriceCache(ctx, year); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsClient.FindYearlyTotalPrices(ctx, &pb.FindYearTotalPrices{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryYearlyTotalPrice(res)
	h.cache.SetCachedYearTotalPriceCache(ctx, year, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get monthly total pricing by merchant
// @Tags Category Stats
// @Description Retrieve monthly total revenue for categories by merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice "Monthly total pricing by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/merchant/monthly-total-pricing [get]
func (h *categoryStatsHandlerApi) FindMonthTotalPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	merchantId, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))

	ctx := r.Context()
	req := &requests.MonthTotalPriceMerchant{Year: year, Month: month, MerchantID: merchantId}

	if cached, found := h.cache.GetCachedMonthTotalPriceByMerchantCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByMerchantClient.FindMonthlyTotalPricesByMerchant(ctx, &pb.FindYearMonthTotalPriceByMerchant{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryMonthlyTotalPrice(res)
	h.cache.SetCachedMonthTotalPriceByMerchantCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get yearly total pricing by merchant
// @Tags Category Stats
// @Description Retrieve yearly total revenue for categories by merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice "Yearly total pricing by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/merchant/yearly-total-pricing [get]
func (h *categoryStatsHandlerApi) FindYearTotalPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantId, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))

	ctx := r.Context()
	req := &requests.YearTotalPriceMerchant{Year: year, MerchantID: merchantId}

	if cached, found := h.cache.GetCachedYearTotalPriceByMerchantCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByMerchantClient.FindYearlyTotalPricesByMerchant(ctx, &pb.FindYearTotalPriceByMerchant{
		Year: int32(year), MerchantId: int32(merchantId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryYearlyTotalPrice(res)
	h.cache.SetCachedYearTotalPriceByMerchantCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get monthly total pricing by category ID
// @Tags Category Stats
// @Description Retrieve monthly total revenue for a specific category
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice "Monthly total pricing by category ID"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/mycategory/monthly-total-pricing [get]
func (h *categoryStatsHandlerApi) FindMonthTotalPriceById(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	categoryId, _ := strconv.Atoi(r.URL.Query().Get("category_id"))

	ctx := r.Context()
	req := &requests.MonthTotalPriceCategory{Year: year, Month: month, CategoryID: categoryId}

	if cached, found := h.cache.GetCachedMonthTotalPriceByIdCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByIdClient.FindMonthlyTotalPricesById(ctx, &pb.FindYearMonthTotalPriceById{
		Year: int32(year), Month: int32(month), CategoryId: int32(categoryId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryMonthlyTotalPrice(res)
	h.cache.SetCachedMonthTotalPriceByIdCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get yearly total pricing by category ID
// @Tags Category Stats
// @Description Retrieve yearly total revenue for a specific category
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice "Yearly total pricing by category ID"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/mycategory/yearly-total-pricing [get]
func (h *categoryStatsHandlerApi) FindYearTotalPriceById(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	categoryId, _ := strconv.Atoi(r.URL.Query().Get("category_id"))

	ctx := r.Context()
	req := &requests.YearTotalPriceCategory{Year: year, CategoryID: categoryId}

	if cached, found := h.cache.GetCachedYearTotalPriceByIdCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByIdClient.FindYearlyTotalPricesById(ctx, &pb.FindYearTotalPriceById{
		Year: int32(year), CategoryId: int32(categoryId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryYearlyTotalPrice(res)
	h.cache.SetCachedYearTotalPriceByIdCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get monthly pricing stats
// @Tags Category Stats
// @Description Retrieve monthly pricing statistics for categories
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly pricing stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/monthly-pricing [get]
func (h *categoryStatsHandlerApi) FindMonthPrice(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	ctx := r.Context()
	if cached, found := h.cache.GetCachedMonthPriceCache(ctx, month, year); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsClient.FindMonthPrice(ctx, &pb.FindYearCategory{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryMonthPrice(res)
	h.cache.SetCachedMonthPriceCache(ctx, month, year, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get yearly pricing stats
// @Tags Category Stats
// @Description Retrieve yearly pricing statistics for categories
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly pricing stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/yearly-pricing [get]
func (h *categoryStatsHandlerApi) FindYearPrice(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))

	ctx := r.Context()
	if cached, found := h.cache.GetCachedYearPriceCache(ctx, year); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsClient.FindYearPrice(ctx, &pb.FindYearCategory{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryYearPrice(res)
	h.cache.SetCachedYearPriceCache(ctx, year, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get monthly pricing stats by merchant
// @Tags Category Stats
// @Description Retrieve monthly pricing statistics for categories by merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly pricing stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/merchant/monthly-pricing [get]
func (h *categoryStatsHandlerApi) FindMonthPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantId, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))

	ctx := r.Context()
	req := &requests.MonthPriceMerchant{Year: year, MerchantID: merchantId}

	if cached, found := h.cache.GetCachedMonthPriceByMerchantCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByMerchantClient.FindMonthPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year: int32(year), MerchantId: int32(merchantId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryMonthPrice(res)
	h.cache.SetCachedMonthPriceByMerchantCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get yearly pricing stats by merchant
// @Tags Category Stats
// @Description Retrieve yearly pricing statistics for categories by merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly pricing stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/merchant/yearly-pricing [get]
func (h *categoryStatsHandlerApi) FindYearPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantId, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))

	ctx := r.Context()
	req := &requests.YearPriceMerchant{Year: year, MerchantID: merchantId}

	if cached, found := h.cache.GetCachedYearPriceByMerchantCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByMerchantClient.FindYearPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year: int32(year), MerchantId: int32(merchantId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryYearPrice(res)
	h.cache.SetCachedYearPriceByMerchantCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get monthly pricing stats by category ID
// @Tags Category Stats
// @Description Retrieve monthly pricing statistics for a specific category
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly pricing stats by category ID"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/mycategory/monthly-pricing [get]
func (h *categoryStatsHandlerApi) FindMonthPriceById(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	categoryId, _ := strconv.Atoi(r.URL.Query().Get("category_id"))

	ctx := r.Context()
	req := &requests.MonthPriceId{Year: year, CategoryID: categoryId}

	if cached, found := h.cache.GetCachedMonthPriceByIdCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByIdClient.FindMonthPriceById(ctx, &pb.FindYearCategoryById{
		Year: int32(year), CategoryId: int32(categoryId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryMonthPrice(res)
	h.cache.SetCachedMonthPriceByIdCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}

// @Security Bearer
// @Summary Get yearly pricing stats by category ID
// @Tags Category Stats
// @Description Retrieve yearly pricing statistics for a specific category
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly pricing stats by category ID"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/category-stats/mycategory/yearly-pricing [get]
func (h *categoryStatsHandlerApi) FindYearPriceById(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	categoryId, _ := strconv.Atoi(r.URL.Query().Get("category_id"))

	ctx := r.Context()
	req := &requests.YearPriceId{Year: year, CategoryID: categoryId}

	if cached, found := h.cache.GetCachedYearPriceByIdCache(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.statsByIdClient.FindYearPriceById(ctx, &pb.FindYearCategoryById{
		Year: int32(year), CategoryId: int32(categoryId),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	response := h.mapper.ToApiResponseCategoryYearPrice(res)
	h.cache.SetCachedYearPriceByIdCache(ctx, req, response)

	return httpx.JSON(w, http.StatusOK, response)
}
