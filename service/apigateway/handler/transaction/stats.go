package transactionhandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/transaction"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/httpx"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/go-chi/chi/v5"
)

type transactionStatsHandlerApi struct {
	statsClient           pb.TransactionStatsServiceClient
	statsByMerchantClient pb.TransactionStatsByMerchantServiceClient
	logger                logger.LoggerInterface
	statsMapper           apimapper.TransactionStatsResponseMapper
	statsCache            transaction_cache.TransactionStatsCache
	statsByMerchantCache  transaction_cache.TransactionStatsByMerchantCache
	apiHandler            apierror.ApiHandler
}

type transactionStatsHandleDeps struct {
	statsClient           pb.TransactionStatsServiceClient
	statsByMerchantClient pb.TransactionStatsByMerchantServiceClient
	router                chi.Router
	logger                logger.LoggerInterface
	statsMapper           apimapper.TransactionStatsResponseMapper
	statsCache            transaction_cache.TransactionStatsCache
	statsByMerchantCache  transaction_cache.TransactionStatsByMerchantCache
	apiHandler            apierror.ApiHandler
}

func NewTransactionStatsHandleApi(params *transactionStatsHandleDeps) *transactionStatsHandlerApi {
	handler := &transactionStatsHandlerApi{
		statsClient:           params.statsClient,
		statsByMerchantClient: params.statsByMerchantClient,
		logger:                params.logger,
		statsMapper:           params.statsMapper,
		statsCache:            params.statsCache,
		statsByMerchantCache:  params.statsByMerchantCache,
		apiHandler:            params.apiHandler,
	}

	params.router.Route("/api/transaction-stats", func(routerTransaction chi.Router) {

		// Stats
		routerTransaction.Get("/monthly-success", params.apiHandler.Handle("GetMonthlyAmountSuccess", handler.FindMonthStatusSuccess))
		routerTransaction.Get("/yearly-success", params.apiHandler.Handle("GetYearlyAmountSuccess", handler.FindYearStatusSuccess))
		routerTransaction.Get("/monthly-failed", params.apiHandler.Handle("GetMonthlyAmountFailed", handler.FindMonthStatusFailed))
		routerTransaction.Get("/yearly-failed", params.apiHandler.Handle("GetYearlyAmountFailed", handler.FindYearStatusFailed))

		routerTransaction.Get("/merchant/monthly-success", params.apiHandler.Handle("GetMonthlyAmountSuccessByMerchant", handler.FindMonthStatusSuccessByMerchant))
		routerTransaction.Get("/merchant/yearly-success", params.apiHandler.Handle("GetYearlyAmountSuccessByMerchant", handler.FindYearStatusSuccessByMerchant))
		routerTransaction.Get("/merchant/monthly-failed", params.apiHandler.Handle("GetMonthlyAmountFailedByMerchant", handler.FindMonthStatusFailedByMerchant))
		routerTransaction.Get("/merchant/yearly-failed", params.apiHandler.Handle("GetYearlyAmountFailedByMerchant", handler.FindYearStatusFailedByMerchant))

		routerTransaction.Get("/monthly-method-success", params.apiHandler.Handle("GetMonthlyTransactionMethodSuccess", handler.FindMonthMethodSuccess))
		routerTransaction.Get("/yearly-method-success", params.apiHandler.Handle("GetYearlyTransactionMethodSuccess", handler.FindYearMethodSuccess))
		routerTransaction.Get("/merchant/monthly-method-success", params.apiHandler.Handle("GetMonthlyTransactionMethodByMerchantSuccess", handler.FindMonthMethodByMerchantSuccess))
		routerTransaction.Get("/merchant/yearly-method-success", params.apiHandler.Handle("GetYearlyTransactionMethodByMerchantSuccess", handler.FindYearMethodByMerchantSuccess))

		routerTransaction.Get("/monthly-method-failed", params.apiHandler.Handle("GetMonthlyTransactionMethodFailed", handler.FindMonthMethodFailed))
		routerTransaction.Get("/yearly-method-failed", params.apiHandler.Handle("GetYearlyTransactionMethodFailed", handler.FindYearMethodFailed))
		routerTransaction.Get("/merchant/monthly-method-failed", params.apiHandler.Handle("GetMonthlyTransactionMethodByMerchantFailed", handler.FindMonthMethodByMerchantFailed))
		routerTransaction.Get("/merchant/yearly-method-failed", params.apiHandler.Handle("GetYearlyTransactionMethodByMerchantFailed", handler.FindYearMethodByMerchantFailed))

	})
	return handler
}

// @Security Bearer
// @Summary Get monthly successful transaction amount
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess "Monthly successful transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-success [get]
func (h *transactionStatsHandlerApi) FindMonthStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	req := &requests.MonthAmountTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthAmountSuccessCached(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyAmountSuccess(ctx, &pb.MonthAmountTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountSuccess(res)
	h.statsCache.SetCachedMonthAmountSuccessCached(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction amount
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess "Yearly successful transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-success [get]
func (h *transactionStatsHandlerApi) FindYearStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()

	if cachedData, found := h.statsCache.GetCachedYearAmountSuccessCached(ctx, year); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyAmountSuccess(ctx, &pb.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountSuccess(res)
	h.statsCache.SetCachedYearAmountSuccessCached(ctx, year, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction amount
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed "Monthly failed transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-failed [get]
func (h *transactionStatsHandlerApi) FindMonthStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	req := &requests.MonthAmountTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthAmountFailedCached(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyAmountFailed(ctx, &pb.MonthAmountTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountFailed(res)
	h.statsCache.SetCachedMonthAmountFailedCached(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction amount
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed "Yearly failed transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-failed [get]
func (h *transactionStatsHandlerApi) FindYearStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()

	if cachedData, found := h.statsCache.GetCachedYearAmountFailedCached(ctx, year); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyAmountFailed(ctx, &pb.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountFailed(res)
	h.statsCache.SetCachedYearAmountFailedCached(ctx, year, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly successful transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess "Monthly successful transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-success [get]
func (h *transactionStatsHandlerApi) FindMonthStatusSuccessByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.MonthAmountTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthAmountSuccessByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyAmountSuccessByMerchant(ctx, &pb.MonthAmountTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountSuccess(res)
	h.statsByMerchantCache.SetCachedMonthAmountSuccessByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess "Yearly successful transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-success [get]
func (h *transactionStatsHandlerApi) FindYearStatusSuccessByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.YearAmountTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearAmountSuccessByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyAmountSuccessByMerchant(ctx, &pb.YearAmountTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountSuccess(res)
	h.statsByMerchantCache.SetCachedYearAmountSuccessByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed "Monthly failed transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-failed [get]
func (h *transactionStatsHandlerApi) FindMonthStatusFailedByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.MonthAmountTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthAmountFailedByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyAmountFailedByMerchant(ctx, &pb.MonthAmountTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountFailed(res)
	h.statsByMerchantCache.SetCachedMonthAmountFailedByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed "Yearly failed transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-failed [get]
func (h *transactionStatsHandlerApi) FindYearStatusFailedByMerchant(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.YearAmountTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearAmountFailedByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyAmountFailedByMerchant(ctx, &pb.YearAmountTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountFailed(res)
	h.statsByMerchantCache.SetCachedYearAmountFailedByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly successful transaction method stats
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly successful transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-method-success [get]
func (h *transactionStatsHandlerApi) FindMonthMethodSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	req := &requests.MonthMethodTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthMethodSuccessCached(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyTransactionMethodSuccess(ctx, &pb.MonthMethodTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsCache.SetCachedMonthMethodSuccessCached(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction method stats
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly successful transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-method-success [get]
func (h *transactionStatsHandlerApi) FindYearMethodSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()

	if cachedData, found := h.statsCache.GetCachedYearMethodSuccessCached(ctx, year); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyTransactionMethodSuccess(ctx, &pb.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsCache.SetCachedYearMethodSuccessCached(ctx, year, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly successful transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly successful transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-method-success [get]
func (h *transactionStatsHandlerApi) FindMonthMethodByMerchantSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.MonthMethodTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthMethodSuccessByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyTransactionMethodByMerchantSuccess(ctx, &pb.MonthMethodTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsByMerchantCache.SetCachedMonthMethodSuccessByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly successful transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-method-success [get]
func (h *transactionStatsHandlerApi) FindYearMethodByMerchantSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.YearMethodTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearMethodSuccessByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyTransactionMethodByMerchantSuccess(ctx, &pb.YearMethodTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsByMerchantCache.SetCachedYearMethodSuccessByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction method stats
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly failed transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-method-failed [get]
func (h *transactionStatsHandlerApi) FindMonthMethodFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	req := &requests.MonthMethodTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthMethodFailedCached(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyTransactionMethodFailed(ctx, &pb.MonthMethodTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsCache.SetCachedMonthMethodFailedCached(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction method stats
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly failed transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-method-failed [get]
func (h *transactionStatsHandlerApi) FindYearMethodFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()

	if cachedData, found := h.statsCache.GetCachedYearMethodFailedCached(ctx, year); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyTransactionMethodFailed(ctx, &pb.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsCache.SetCachedYearMethodFailedCached(ctx, year, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly failed transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-method-failed [get]
func (h *transactionStatsHandlerApi) FindMonthMethodByMerchantFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.MonthMethodTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthMethodFailedByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyTransactionMethodByMerchantFailed(ctx, &pb.MonthMethodTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsByMerchantCache.SetCachedMonthMethodFailedByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly failed transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-method-failed [get]
func (h *transactionStatsHandlerApi) FindYearMethodByMerchantFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	merchantID, _ := strconv.Atoi(r.URL.Query().Get("merchant_id"))
	ctx := r.Context()
	req := &requests.YearMethodTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearMethodFailedByMerchant(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyTransactionMethodByMerchantFailed(ctx, &pb.YearMethodTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsByMerchantCache.SetCachedYearMethodFailedByMerchant(ctx, req, apiResponse)
	return httpx.JSON(w, http.StatusOK, apiResponse)
}
