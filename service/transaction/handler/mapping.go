package handler

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func getStringValue(v interface{}) *wrapperspb.StringValue {
	switch t := v.(type) {
	case time.Time:
		if !t.IsZero() {
			return wrapperspb.String(t.Format("2006-01-02 15:04:05"))
		}
	case *time.Time:
		if t != nil && !t.IsZero() {
			return wrapperspb.String(t.Format("2006-01-02 15:04:05"))
		}
	case string:
		if t != "" {
			return wrapperspb.String(t)
		}
	}
	return nil
}

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case time.Time:
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	case *time.Time:
		if t == nil || t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	case string:
		return t
	}
	return ""
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoTransactionResponse(v *repository.TransactionResult) *pbtransactions.TransactionResponse {
	if v == nil {
		return nil
	}
	return &pbtransactions.TransactionResponse{
		Id:            v.TransactionID,
		OrderId:       v.OrderID,
		MerchantId:    v.MerchantID,
		PaymentMethod: v.PaymentMethod,
		Amount:        v.Amount,
		PaymentStatus: v.PaymentStatus,
		CreatedAt:     formatTimestamp(v.CreatedAt),
		UpdatedAt:     formatTimestamp(v.UpdatedAt),
	}
}

func mapToProtoTransactionResponseDeleteAt(v *repository.TransactionResult) *pbtransactions.TransactionResponseDeleteAt {
	if v == nil {
		return nil
	}
	return &pbtransactions.TransactionResponseDeleteAt{
		Id:            v.TransactionID,
		OrderId:       v.OrderID,
		MerchantId:    v.MerchantID,
		PaymentMethod: v.PaymentMethod,
		Amount:        v.Amount,
		PaymentStatus: v.PaymentStatus,
		CreatedAt:     formatTimestamp(v.CreatedAt),
		UpdatedAt:     formatTimestamp(v.UpdatedAt),
		DeletedAt:     getStringValue(v.DeletedAt),
	}
}
