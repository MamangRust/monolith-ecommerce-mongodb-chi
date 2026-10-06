package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandHandler struct {
	pbtransactions.UnimplementedTransactionCommandServiceServer
	service service.TransactionCommandService
	logger  logger.LoggerInterface
}

func NewTransactionCommandHandler(service service.TransactionCommandService, logger logger.LoggerInterface) *transactionCommandHandler {
	return &transactionCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionCommandHandler) Create(ctx context.Context, req *pbtransactions.CreateTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	request := &requests.CreateTransactionRequest{
		OrderID:       int(req.GetOrderId()),
		MerchantID:    int(req.GetMerchantId()),
		UserID:        int(req.GetUserId()),
		PaymentMethod: req.GetPaymentMethod(),
		Amount:        int(req.GetAmount()),
		PaymentStatus: &[]string{req.GetPaymentStatus()}[0],
	}

	data, err := h.service.Create(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    h.ToTransactionResponseCreate(data),
	}, nil
}

func (h *transactionCommandHandler) Update(ctx context.Context, req *pbtransactions.UpdateTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	transactionID := int(req.GetTransactionId())
	request := &requests.UpdateTransactionRequest{
		TransactionID: &transactionID,
		MerchantID:    int(req.GetMerchantId()),
		OrderID:       int(req.GetOrderId()),
		PaymentMethod: req.GetPaymentMethod(),
		Amount:        int(req.GetAmount()),
		PaymentStatus: &[]string{req.GetPaymentStatus()}[0],
	}

	data, err := h.service.Update(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    h.ToTransactionResponseUpdate(data),
	}, nil
}

func (h *transactionCommandHandler) TrashedTransaction(ctx context.Context, req *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDeleteAt, error) {
	data, err := h.service.Trash(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    h.ToTransactionResponseDeleteAt(data),
	}, nil
}

func (h *transactionCommandHandler) RestoreTransaction(ctx context.Context, req *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDeleteAt, error) {
	data, err := h.service.Restore(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    h.ToTransactionResponseDeleteAt(data),
	}, nil
}

func (h *transactionCommandHandler) DeleteTransactionPermanent(ctx context.Context, req *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDelete, error) {
	_, err := h.service.DeletePermanent(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction permanently",
	}, nil
}

func (h *transactionCommandHandler) RestoreAllTransaction(ctx context.Context, req *emptypb.Empty) (*pbtransactions.ApiResponseTransactionAll, error) {
	_, err := h.service.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restored all transactions",
	}, nil
}

func (h *transactionCommandHandler) DeleteTransactionByOrderPermanent(ctx context.Context, req *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDelete, error) {
	_, err := h.service.DeleteByOrderIDPermanent(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transactions by order permanently",
	}, nil
}

func (h *transactionCommandHandler) DeleteAllTransactionPermanent(ctx context.Context, req *emptypb.Empty) (*pbtransactions.ApiResponseTransactionAll, error) {
	_, err := h.service.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransactions.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully deleted all transactions permanently",
	}, nil
}

// Manual Mappings

func (h *transactionCommandHandler) ToTransactionResponseCreate(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionCommandHandler) ToTransactionResponseUpdate(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionCommandHandler) ToTransactionResponseDeleteAt(v *repository.TransactionResult) *pbtransactions.TransactionResponseDeleteAt {
	return mapToProtoTransactionResponseDeleteAt(v)
}
