package handler

import (
	"math"
	"time"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
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

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case time.Time:
		return convert.FormatTimeRFC3339(&t)
	case *time.Time:
		return convert.FormatTimeRFC3339(t)
	default:
		return ""
	}
}

func deletedAtWrappers(v interface{}) *wrapperspb.StringValue {
	switch t := v.(type) {
	case *time.Time:
		return convert.TimeToWrappers(t)
	case time.Time:
		return convert.TimeToWrappers(&t)
	default:
		return &wrapperspb.StringValue{Value: ""}
	}
}

func mapToProtoUserResponse(m interface{}) *pbusers.UserResponse {
	switch v := m.(type) {
	case *salemongo.User:
		return &pbusers.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.UserRow:
		return &pbusers.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.UserByIDRow:
		return &pbusers.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoUserResponseDeleteAt(m interface{}) *pbusers.UserResponseDeleteAt {
	switch v := m.(type) {
	case *salemongo.User:
		return &pbusers.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
			DeletedAt: deletedAtWrappers(v.DeletedAt),
		}
	case *salemongo.UserActiveRow:
		return &pbusers.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
			DeletedAt: deletedAtWrappers(v.DeletedAt),
		}
	default:
		return nil
	}
}

func mapToProtoUserResponseWithPassword(m interface{}) *pbusers.UserResponseWithPassword {
	switch v := m.(type) {
	case *salemongo.User:
		return &pbusers.UserResponseWithPassword{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			Password:  v.Password,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *salemongo.UserByEmailWithPasswordRow:
		return &pbusers.UserResponseWithPassword{
			Id:       int32(v.UserID),
			Email:    v.Email,
			Password: v.Password,
		}
	default:
		return nil
	}
}
