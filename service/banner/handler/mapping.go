package handler

import (
	"math"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	pbbanners "github.com/MamangRust/monolith-ecommerce-pb/banners"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
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

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return convert.FormatTimeRFC3339(&t)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return formatTime(*t)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func formatTimeOfDay(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case *string:
		if t != nil {
			return *t
		}
	}
	return ""
}

func deletedAtWrappersPtr(t *time.Time) *wrapperspb.StringValue {
	if t == nil {
		return nil
	}
	return convert.TimeToWrappers(t)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func mapToProtoBannerResponse(m interface{}) *pbbanners.BannerResponse {
	switch v := m.(type) {
	case *salemongo.Banner:
		return &pbbanners.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: formatDate(v.StartDate),
			EndDate:   formatDate(v.EndDate),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			IsActive:  *v.IsActive,
			CreatedAt: formatTimePtr(v.CreatedAt),
			UpdatedAt: formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.BannerRow:
		return &pbbanners.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: formatDate(v.StartDate),
			EndDate:   formatDate(v.EndDate),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			IsActive:  *v.IsActive,
			CreatedAt: formatTime(v.CreatedAt),
			UpdatedAt: formatTime(v.UpdatedAt),
		}
	case *repository.BannerResult:
		return &pbbanners.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: derefString(v.StartDate),
			EndDate:   derefString(v.EndDate),
			StartTime: derefString(v.StartTime),
			EndTime:   derefString(v.EndTime),
			IsActive:  derefBool(v.IsActive),
			CreatedAt: derefString(v.CreatedAt),
			UpdatedAt: derefString(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoBannerResponseDeleteAt(m interface{}) *pbbanners.BannerResponseDeleteAt {
	switch v := m.(type) {
	case *salemongo.Banner:
		res := &pbbanners.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: formatDate(v.StartDate),
			EndDate:   formatDate(v.EndDate),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			IsActive:  *v.IsActive,
			CreatedAt: formatTimePtr(v.CreatedAt),
			UpdatedAt: formatTimePtr(v.UpdatedAt),
		}
		if v.DeletedAt != nil {
			res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
		}
		return res
	case *salemongo.BannerActiveRow:
		res := &pbbanners.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: formatDate(v.StartDate),
			EndDate:   formatDate(v.EndDate),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			IsActive:  *v.IsActive,
			CreatedAt: formatTime(v.CreatedAt),
			UpdatedAt: formatTime(v.UpdatedAt),
		}
		if v.DeletedAt != nil {
			res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
		}
		return res
	case *repository.BannerResult:
		res := &pbbanners.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: derefString(v.StartDate),
			EndDate:   derefString(v.EndDate),
			StartTime: derefString(v.StartTime),
			EndTime:   derefString(v.EndTime),
			IsActive:  derefBool(v.IsActive),
			CreatedAt: derefString(v.CreatedAt),
			UpdatedAt: derefString(v.UpdatedAt),
		}
		if v.DeletedAt != nil && *v.DeletedAt != "" {
			res.DeletedAt = &wrapperspb.StringValue{Value: *v.DeletedAt}
		}
		return res
	default:
		return nil
	}
}
