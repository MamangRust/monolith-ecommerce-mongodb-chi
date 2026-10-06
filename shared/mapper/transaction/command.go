package transactionapimapper

import (
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type transactionCommandResponseMapper struct{}

func NewTransactionCommandResponseMapper() TransactionCommandResponseMapper {
	return &transactionCommandResponseMapper{}
}

func (t *transactionCommandResponseMapper) ToResponseTransaction(transaction *pbtransactions.TransactionResponse) *response.TransactionResponse {
	if transaction == nil {
		return nil
	}
	return &response.TransactionResponse{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionCommandResponseMapper) ToResponsesTransaction(transactions []*pbtransactions.TransactionResponse) []*response.TransactionResponse {
	var mappedTransactions []*response.TransactionResponse
	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransaction(transaction))
	}
	return mappedTransactions
}

func (t *transactionCommandResponseMapper) ToApiResponseTransaction(pbResponse *pbtransactions.ApiResponseTransaction) *response.ApiResponseTransaction {
	return &response.ApiResponseTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransaction(pbResponse.Data),
	}
}

func (t *transactionCommandResponseMapper) ToResponseTransactionDeleteAt(transaction *pbtransactions.TransactionResponseDeleteAt) *response.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var deletedAt string
	if transaction.DeletedAt != nil {
		deletedAt = transaction.DeletedAt.Value
	}

	return &response.TransactionResponseDeleteAt{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (t *transactionCommandResponseMapper) ToResponsesTransactionDeleteAt(transactions []*pbtransactions.TransactionResponseDeleteAt) []*response.TransactionResponseDeleteAt {
	var mappedTransactions []*response.TransactionResponseDeleteAt
	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransactionDeleteAt(transaction))
	}
	return mappedTransactions
}

func (t *transactionCommandResponseMapper) ToApiResponseTransactionDeleteAt(pbResponse *pbtransactions.ApiResponseTransactionDeleteAt) *response.ApiResponseTransactionDeleteAt {
	return &response.ApiResponseTransactionDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransactionDeleteAt(pbResponse.Data),
	}
}

func (t *transactionCommandResponseMapper) ToApiResponseTransactionDelete(pbResponse *pbtransactions.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete {
	return &response.ApiResponseTransactionDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionCommandResponseMapper) ToApiResponseTransactionAll(pbResponse *pbtransactions.ApiResponseTransactionAll) *response.ApiResponseTransactionAll {
	return &response.ApiResponseTransactionAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionCommandResponseMapper) ToApiResponsePaginationTransactionDeleteAt(pbResponse *pbtransactions.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt {
	return &response.ApiResponsePaginationTransactionDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransactionDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
