package handler

import "time"

import (
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"math"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchants_detail "github.com/MamangRust/monolith-ecommerce-pb/merchants_detail"
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

func mapToProtoMerchantDetailResponse(m interface{}) *pbmerchants_detail.MerchantDetailResponse {
	switch v := m.(type) {
	case *salemongo.MerchantDetail:
		return &pbmerchants_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.MerchantDetailWithJoinsRow:
		return &pbmerchants_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.MerchantDetailRow:
		return &pbmerchants_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.MerchantDetailListRow:
		return &pbmerchants_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDetailResponseDeleteAt(m interface{}) *pbmerchants_detail.MerchantDetailResponseDeleteAt {
	var res *pbmerchants_detail.MerchantDetailResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *salemongo.MerchantDetail:
		res = &pbmerchants_detail.MerchantDetailResponseDeleteAt{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.MerchantDetailActiveRow:
		res = &pbmerchants_detail.MerchantDetailResponseDeleteAt{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
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

func mapToProtoMerchantSocialLinkResponse(m interface{}) *pbmerchants_detail.MerchantSocialMediaLinkResponse {
	switch v := m.(type) {
	case *salemongo.MerchantSocialMediaLink:
		return &pbmerchants_detail.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.MerchantSocialMediaLinkRow:
		return &pbmerchants_detail.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}
