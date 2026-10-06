package handler

import (
	"math"
	"time"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
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

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case time.Time:
		return t.Format("2006-01-02 15:04:05.000")
	case *time.Time:
		if t != nil {
			return t.Format("2006-01-02 15:04:05.000")
		}
	}
	return ""
}

func mapToProtoRoleResponse(m interface{}) *pbroles.RoleResponse {
	switch v := m.(type) {
	case *salemongo.Role:
		return &pbroles.RoleResponse{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.RoleRow:
		return &pbroles.RoleResponse{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoRoleResponseDeleteAt(m interface{}) *pbroles.RoleResponseDeleteAt {
	var res *pbroles.RoleResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *salemongo.Role:
		res = &pbroles.RoleResponseDeleteAt{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.RoleActiveRow:
		res = &pbroles.RoleResponseDeleteAt{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
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
