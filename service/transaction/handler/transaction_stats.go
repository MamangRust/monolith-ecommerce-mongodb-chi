package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type transactionStatsHandler struct {
	pbtransactions.UnimplementedTransactionStatsServiceServer
	service service.TransactionStatsService
	logger  logger.LoggerInterface
}

func NewTransactionStatsHandler(service service.TransactionStatsService, logger logger.LoggerInterface) *transactionStatsHandler {
	return &transactionStatsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionStatsHandler) GetMonthlyAmountSuccess(ctx context.Context, req *pbtransactions.MonthAmountTransactionRequest) (*pbtransactions.ApiResponseTransactionMonthAmountSuccess, error) {
	request := &requests.MonthAmountTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountSuccess(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionMonthlyAmountSuccess
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionMonthlyAmountSuccess{
			Year:         v.Year,
			Month:        v.Month,
			TotalSuccess: int32(v.TotalSuccess),
			TotalAmount:  int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly amount success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyAmountSuccess(ctx context.Context, req *pbtransactions.YearAmountTransactionRequest) (*pbtransactions.ApiResponseTransactionYearAmountSuccess, error) {
	data, err := h.service.FindYearlyAmountSuccess(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionYearlyAmountSuccess
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionYearlyAmountSuccess{
			Year:         v.Year,
			TotalSuccess: int32(v.TotalSuccess),
			TotalAmount:  int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly amount success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyAmountFailed(ctx context.Context, req *pbtransactions.MonthAmountTransactionRequest) (*pbtransactions.ApiResponseTransactionMonthAmountFailed, error) {
	request := &requests.MonthAmountTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountFailed(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionMonthlyAmountFailed
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionMonthlyAmountFailed{
			Year:        v.Year,
			Month:       v.Month,
			TotalFailed: int32(v.TotalFailed),
			TotalAmount: int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Successfully fetched monthly amount failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyAmountFailed(ctx context.Context, req *pbtransactions.YearAmountTransactionRequest) (*pbtransactions.ApiResponseTransactionYearAmountFailed, error) {
	data, err := h.service.FindYearlyAmountFailed(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionYearlyAmountFailed
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionYearlyAmountFailed{
			Year:        v.Year,
			TotalFailed: int32(v.TotalFailed),
			TotalAmount: int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Successfully fetched yearly amount failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *pbtransactions.MonthMethodTransactionRequest) (*pbtransactions.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodSuccess(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionMonthlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionMonthlyMethod{
			Month:             v.Month,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyTransactionMethodSuccess(ctx context.Context, req *pbtransactions.YearMethodTransactionRequest) (*pbtransactions.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.service.FindYearlyMethodSuccess(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionYearlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionYearlyMethod{
			Year:              v.Year,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyTransactionMethodFailed(ctx context.Context, req *pbtransactions.MonthMethodTransactionRequest) (*pbtransactions.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodFailed(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionMonthlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionMonthlyMethod{
			Month:             v.Month,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyTransactionMethodFailed(ctx context.Context, req *pbtransactions.YearMethodTransactionRequest) (*pbtransactions.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.service.FindYearlyMethodFailed(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransactions.TransactionYearlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransactions.TransactionYearlyMethod{
			Year:              v.Year,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransactions.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method failed stats",
		Data:    stats,
	}, nil
}
