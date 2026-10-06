package handler

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-ecommerce-pb/merchants_document"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func getString(s *string) string { return convert.StrVal(s) }

func getStringFromTime(t *time.Time) *wrapperspb.StringValue {
	return convert.TimeToWrappers(t)
}

func formatTime(t time.Time) string {
	return convert.FormatTimeRFC3339(&t)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return convert.FormatTimeRFC3339(t)
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

func mapToProtoMerchantResponse(m interface{}) *pbmerchants.MerchantResponse {
	switch v := m.(type) {
	case *salemongo.Merchant:
		return &pbmerchants.MerchantResponse{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.MerchantRow:
		return &pbmerchants.MerchantResponse{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	case *repository.MerchantResult:
		return &pbmerchants.MerchantResponse{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    getString(v.CreatedAt),
			UpdatedAt:    getString(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantResponseDeleteAt(m interface{}) *pbmerchants.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *salemongo.Merchant:
		return &pbmerchants.MerchantResponseDeleteAt{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.MerchantActiveRow:
		return &pbmerchants.MerchantResponseDeleteAt{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	case *repository.MerchantResult:
		return &pbmerchants.MerchantResponseDeleteAt{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    getString(v.CreatedAt),
			UpdatedAt:    getString(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantResponseTrashed(m interface{}) *pbmerchants.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *salemongo.MerchantTrashedRow:
		return &pbmerchants.MerchantResponseDeleteAt{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	case *repository.MerchantResult:
		return &pbmerchants.MerchantResponseDeleteAt{
			Id:           v.MerchantID,
			UserId:       v.UserID,
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    getString(v.CreatedAt),
			UpdatedAt:    getString(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponse(m interface{}) *pbmerchants_document.MerchantDocument {
	switch v := m.(type) {
	case *salemongo.MerchantDocument:
		return &pbmerchants_document.MerchantDocument{
			DocumentId:   v.DocumentID,
			MerchantId:   v.MerchantID,
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.MerchantDocumentRow:
		return &pbmerchants_document.MerchantDocument{
			DocumentId:   v.DocumentID,
			MerchantId:   v.MerchantID,
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	case *salemongo.MerchantDocumentsRow:
		return &pbmerchants_document.MerchantDocument{
			DocumentId:   v.DocumentID,
			MerchantId:   v.MerchantID,
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	case *salemongo.MerchantDocumentsActiveRow:
		return &pbmerchants_document.MerchantDocument{
			DocumentId:   v.DocumentID,
			MerchantId:   v.MerchantID,
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponseAt(m interface{}) *pbmerchants_document.MerchantDocumentDeleteAt {
	switch v := m.(type) {
	case *salemongo.MerchantDocumentsTrashedRow:
		return &pbmerchants_document.MerchantDocumentDeleteAt{
			DocumentId:   v.DocumentID,
			MerchantId:   v.MerchantID,
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTime(v.CreatedAt),
			UpdatedAt:    formatTime(v.UpdatedAt),
		}
	default:
		return nil
	}
}
