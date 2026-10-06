package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	// Removed unused db import
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type transactionStatsByMerchantHandler struct {
	pbtransactions.UnimplementedTransactionStatsByMerchantServiceServer
	service service.TransactionStatsByMerchantService
	logger  logger.LoggerInterface
}

func NewTransactionStatsByMerchantHandler(service service.TransactionStatsByMerchantService, logger logger.LoggerInterface) *transactionStatsByMerchantHandler {
	return &transactionStatsByMerchantHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionStatsByMerchantHandler) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *pbtransactions.MonthAmountTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionMonthAmountSuccess, error) {
	request := &requests.MonthAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountSuccessByMerchant(ctx, request)
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
		Message: "Successfully fetched monthly amount success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *pbtransactions.YearAmountTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionYearAmountSuccess, error) {
	request := &requests.YearAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyAmountSuccessByMerchant(ctx, request)
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
		Message: "Successfully fetched yearly amount success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *pbtransactions.MonthAmountTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionMonthAmountFailed, error) {
	request := &requests.MonthAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountFailedByMerchant(ctx, request)
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
		Message: "Successfully fetched monthly amount failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyAmountFailedByMerchant(ctx context.Context, req *pbtransactions.YearAmountTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionYearAmountFailed, error) {
	request := &requests.YearAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyAmountFailedByMerchant(ctx, request)
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
		Message: "Successfully fetched yearly amount failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pbtransactions.MonthMethodTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodByMerchantSuccess(ctx, request)
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
		Message: "Successfully fetched monthly transaction method success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pbtransactions.YearMethodTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionYearPaymentmethod, error) {
	request := &requests.YearMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyMethodByMerchantSuccess(ctx, request)
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
		Message: "Successfully fetched yearly transaction method success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *pbtransactions.MonthMethodTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodByMerchantFailed(ctx, request)
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
		Message: "Successfully fetched monthly transaction method failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *pbtransactions.YearMethodTransactionMerchantRequest) (*pbtransactions.ApiResponseTransactionYearPaymentmethod, error) {
	request := &requests.YearMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyMethodByMerchantFailed(ctx, request)
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
		Message: "Successfully fetched yearly transaction method failed stats by merchant",
		Data:    stats,
	}, nil
}
