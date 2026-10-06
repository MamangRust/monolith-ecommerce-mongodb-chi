package handler

import (
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"math"
	"time"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchants_award "github.com/MamangRust/monolith-ecommerce-pb/merchants_award"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func getString(s *string) string { return convert.StrVal(s) }

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case time.Time:
		return convert.FormatTimeRFC3339(&t)
	case string:
		return t
	}
	return ""
}

func mapToProtoMerchantAwardResponse(m interface{}) *pbmerchants_award.MerchantAwardResponse {
	var id int32
	var merchantID int32
	var title, description, issuedBy, certificateUrl string
	var issueDate, expiryDate, createdAt, updatedAt string

	switch v := m.(type) {
	case *salemongo.MerchantCertificationsAndAward:
		id, merchantID = v.MerchantCertificationID, v.MerchantID
		title, description, issuedBy, certificateUrl = v.Title, getString(v.Description), getString(v.IssuedBy), getString(v.CertificateUrl)
		issueDate, expiryDate = formatTimestamp(v.IssueDate), formatTimestamp(v.ExpiryDate)
		createdAt, updatedAt = formatTimestamp(v.CreatedAt), formatTimestamp(v.UpdatedAt)
	case *salemongo.MerchantCertificationListRow:
		id, merchantID = v.MerchantCertificationID, v.MerchantID
		title, description, issuedBy, certificateUrl = v.Title, getString(v.Description), getString(v.IssuedBy), getString(v.CertificateUrl)
		issueDate, expiryDate = formatTimestamp(v.IssueDate), formatTimestamp(v.ExpiryDate)
		createdAt, updatedAt = formatTimestamp(v.CreatedAt), formatTimestamp(v.UpdatedAt)
	case *salemongo.MerchantCertificationRow:
		id, merchantID = v.MerchantCertificationID, v.MerchantID
		title, description, issuedBy, certificateUrl = v.Title, getString(v.Description), getString(v.IssuedBy), getString(v.CertificateUrl)
		issueDate, expiryDate = formatTimestamp(v.IssueDate), formatTimestamp(v.ExpiryDate)
		createdAt, updatedAt = formatTimestamp(v.CreatedAt), formatTimestamp(v.UpdatedAt)
	default:
		return nil
	}

	return &pbmerchants_award.MerchantAwardResponse{
		Id:             id,
		MerchantId:     merchantID,
		Title:          title,
		Description:    description,
		IssuedBy:       issuedBy,
		CertificateUrl: certificateUrl,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func mapToProtoMerchantAwardResponseDeleteAt(m interface{}) *pbmerchants_award.MerchantAwardResponseDeleteAt {
	var res *pbmerchants_award.MerchantAwardResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *salemongo.MerchantCertificationsAndAward:
		res = &pbmerchants_award.MerchantAwardResponseDeleteAt{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.MerchantCertificationActiveRow:
		res = &pbmerchants_award.MerchantAwardResponseDeleteAt{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if val := formatTimestamp(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}
