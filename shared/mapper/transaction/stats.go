package transactionapimapper

import (
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type transactionStatsResponseMapper struct{}

func NewTransactionStatsResponseMapper() TransactionStatsResponseMapper {
	return &transactionStatsResponseMapper{}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthAmountSuccess(row *pbtransactions.TransactionMonthlyAmountSuccess) *response.TransactionMonthlyAmountSuccessResponse {
	return &response.TransactionMonthlyAmountSuccessResponse{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthlyAmountSuccess(rows []*pbtransactions.TransactionMonthlyAmountSuccess) []*response.TransactionMonthlyAmountSuccessResponse {
	var mapped []*response.TransactionMonthlyAmountSuccessResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionMonthAmountSuccess(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionYearAmountSuccess(row *pbtransactions.TransactionYearlyAmountSuccess) *response.TransactionYearlyAmountSuccessResponse {
	return &response.TransactionYearlyAmountSuccessResponse{
		Year:         row.Year,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionYearlyAmountSuccess(rows []*pbtransactions.TransactionYearlyAmountSuccess) []*response.TransactionYearlyAmountSuccessResponse {
	var mapped []*response.TransactionYearlyAmountSuccessResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionYearAmountSuccess(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionMonthAmountFailed(row *pbtransactions.TransactionMonthlyAmountFailed) *response.TransactionMonthlyAmountFailedResponse {
	return &response.TransactionMonthlyAmountFailedResponse{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthlyAmountFailed(rows []*pbtransactions.TransactionMonthlyAmountFailed) []*response.TransactionMonthlyAmountFailedResponse {
	var mapped []*response.TransactionMonthlyAmountFailedResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionMonthAmountFailed(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionYearAmountFailed(row *pbtransactions.TransactionYearlyAmountFailed) *response.TransactionYearlyAmountFailedResponse {
	return &response.TransactionYearlyAmountFailedResponse{
		Year:        row.Year,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionYearlyAmountFailed(rows []*pbtransactions.TransactionYearlyAmountFailed) []*response.TransactionYearlyAmountFailedResponse {
	var mapped []*response.TransactionYearlyAmountFailedResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionYearAmountFailed(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionMonthMethod(row *pbtransactions.TransactionMonthlyMethod) *response.TransactionMonthlyMethodResponse {
	return &response.TransactionMonthlyMethodResponse{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthlyMethod(rows []*pbtransactions.TransactionMonthlyMethod) []*response.TransactionMonthlyMethodResponse {
	var mapped []*response.TransactionMonthlyMethodResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionMonthMethod(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionYearMethod(row *pbtransactions.TransactionYearlyMethod) *response.TransactionYearlyMethodResponse {
	return &response.TransactionYearlyMethodResponse{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionYearlyMethod(rows []*pbtransactions.TransactionYearlyMethod) []*response.TransactionYearlyMethodResponse {
	var mapped []*response.TransactionYearlyMethodResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionYearMethod(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionMonthAmountSuccess(pbResponse *pbtransactions.ApiResponseTransactionMonthAmountSuccess) *response.ApiResponsesTransactionMonthSuccess {
	return &response.ApiResponsesTransactionMonthSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionMonthlyAmountSuccess(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionMonthAmountFailed(pbResponse *pbtransactions.ApiResponseTransactionMonthAmountFailed) *response.ApiResponsesTransactionMonthFailed {
	return &response.ApiResponsesTransactionMonthFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionMonthlyAmountFailed(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionYearAmountSuccess(pbResponse *pbtransactions.ApiResponseTransactionYearAmountSuccess) *response.ApiResponsesTransactionYearSuccess {
	return &response.ApiResponsesTransactionYearSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionYearlyAmountSuccess(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionYearAmountFailed(pbResponse *pbtransactions.ApiResponseTransactionYearAmountFailed) *response.ApiResponsesTransactionYearFailed {
	return &response.ApiResponsesTransactionYearFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionYearlyAmountFailed(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionMonthMethod(pbResponse *pbtransactions.ApiResponseTransactionMonthPaymentMethod) *response.ApiResponsesTransactionMonthMethod {
	return &response.ApiResponsesTransactionMonthMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionMonthlyMethod(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionYearMethod(pbResponse *pbtransactions.ApiResponseTransactionYearPaymentmethod) *response.ApiResponsesTransactionYearMethod {
	return &response.ApiResponsesTransactionYearMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionYearlyMethod(pbResponse.Data),
	}
}
