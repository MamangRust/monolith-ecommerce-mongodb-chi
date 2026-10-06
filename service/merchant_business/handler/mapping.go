package handler

import (
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"math"
	"time"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchants_business "github.com/MamangRust/monolith-ecommerce-pb/merchants_business"
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

func getInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case time.Time:
		return convert.FormatTimeRFC3339(&t)
	case string:
		return t
	}
	return ""
}

func mapToProtoMerchantBusinessResponse(m interface{}) *pbmerchants_business.MerchantBusinessResponse {
	switch v := m.(type) {
	case *salemongo.MerchantBusinessInformation:
		return &pbmerchants_business.MerchantBusinessResponse{
			Id:                v.MerchantBusinessInfoID,
			MerchantId:        v.MerchantID,
			BusinessType:      getString(v.BusinessType),
			TaxId:             getString(v.TaxID),
			EstablishedYear:   getInt32(v.EstablishedYear),
			NumberOfEmployees: getInt32(v.NumberOfEmployees),
			WebsiteUrl:        getString(v.WebsiteUrl),
			CreatedAt:         formatTimestamp(v.CreatedAt),
			UpdatedAt:         formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.MerchantBusinessInfoListRow:
		return &pbmerchants_business.MerchantBusinessResponse{
			Id:                v.MerchantBusinessInfoID,
			MerchantId:        v.MerchantID,
			BusinessType:      getString(v.BusinessType),
			TaxId:             getString(v.TaxID),
			EstablishedYear:   getInt32(v.EstablishedYear),
			NumberOfEmployees: getInt32(v.NumberOfEmployees),
			WebsiteUrl:        getString(v.WebsiteUrl),
			CreatedAt:         formatTimestamp(v.CreatedAt),
			UpdatedAt:         formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.MerchantBusinessInfoRow:
		return &pbmerchants_business.MerchantBusinessResponse{
			Id:                v.MerchantBusinessInfoID,
			MerchantId:        v.MerchantID,
			BusinessType:      getString(v.BusinessType),
			TaxId:             getString(v.TaxID),
			EstablishedYear:   getInt32(v.EstablishedYear),
			NumberOfEmployees: getInt32(v.NumberOfEmployees),
			WebsiteUrl:        getString(v.WebsiteUrl),
			CreatedAt:         formatTimestamp(v.CreatedAt),
			UpdatedAt:         formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantBusinessResponseDeleteAt(m interface{}) *pbmerchants_business.MerchantBusinessResponseDeleteAt {
	var res *pbmerchants_business.MerchantBusinessResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *salemongo.MerchantBusinessInformation:
		res = &pbmerchants_business.MerchantBusinessResponseDeleteAt{
			Id:                v.MerchantBusinessInfoID,
			MerchantId:        v.MerchantID,
			BusinessType:      getString(v.BusinessType),
			TaxId:             getString(v.TaxID),
			EstablishedYear:   getInt32(v.EstablishedYear),
			NumberOfEmployees: getInt32(v.NumberOfEmployees),
			WebsiteUrl:        getString(v.WebsiteUrl),
			CreatedAt:         formatTimestamp(v.CreatedAt),
			UpdatedAt:         formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.MerchantBusinessInfoActiveRow:
		res = &pbmerchants_business.MerchantBusinessResponseDeleteAt{
			Id:                v.MerchantBusinessInfoID,
			MerchantId:        v.MerchantID,
			BusinessType:      getString(v.BusinessType),
			TaxId:             getString(v.TaxID),
			EstablishedYear:   getInt32(v.EstablishedYear),
			NumberOfEmployees: getInt32(v.NumberOfEmployees),
			WebsiteUrl:        getString(v.WebsiteUrl),
			CreatedAt:         formatTimestamp(v.CreatedAt),
			UpdatedAt:         formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if val := formatTimestamp(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}
