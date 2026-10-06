package handler

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type transactionQueryHandler struct {
	pbtransactions.UnimplementedTransactionQueryServiceServer
	service service.TransactionQueryService
	logger  logger.LoggerInterface
}

func NewTransactionQueryHandler(service service.TransactionQueryService, logger logger.LoggerInterface) *transactionQueryHandler {
	return &transactionQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionQueryHandler) FindAllTransactions(ctx context.Context, req *pbtransactions.FindAllTransactionRequest) (*pbtransactions.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindAll(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransactions.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponse(v))
	}

	return &pbtransactions.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByActive(ctx context.Context, req *pbtransactions.FindAllTransactionRequest) (*pbtransactions.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindActive(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransactions.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseActive(v))
	}

	return &pbtransactions.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched active transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByTrashed(ctx context.Context, req *pbtransactions.FindAllTransactionRequest) (*pbtransactions.ApiResponsePaginationTransactionDeleteAt, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindTrashed(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransactions.TransactionResponseDeleteAt
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseDeleteAt(v))
	}

	return &pbtransactions.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindById(ctx context.Context, req *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	data, err := h.service.FindByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    h.ToTransactionResponseId(data),
	}, nil
}

func (h *transactionQueryHandler) FindByOrderId(ctx context.Context, req *pbtransactions.FindByOrderIdTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	data, err := h.service.FindByOrderID(ctx, int(req.GetOrderId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction by order id",
		Data:    h.ToTransactionResponseOrderId(data),
	}, nil
}

func (h *transactionQueryHandler) FindByMerchant(ctx context.Context, req *pbtransactions.FindAllTransactionByMerchantRequest) (*pbtransactions.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransactionByMerchant{
		MerchantID: int(req.GetMerchantId()),
		Page:       int(req.GetPage()),
		PageSize:   int(req.GetPageSize()),
		Search:     req.GetSearch(),
	}

	data, total, err := h.service.FindByMerchant(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransactions.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseMerchant(v))
	}

	return &pbtransactions.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions by merchant",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

// Manual Mappings

func (h *transactionQueryHandler) ToTransactionResponse(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseActive(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseDeleteAt(v *repository.TransactionResult) *pbtransactions.TransactionResponseDeleteAt {
	return mapToProtoTransactionResponseDeleteAt(v)
}

func (h *transactionQueryHandler) ToTransactionResponseId(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseOrderId(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseMerchant(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}
