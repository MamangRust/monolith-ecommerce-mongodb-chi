package handler

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbshipping_addresses "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
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
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func formatTime(t time.Time) string {
	return convert.FormatTimeRFC3339(&t)
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return formatDate(*t)
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapToProtoShippingResponse(shipping interface{}) *pbshipping_addresses.ShippingResponse {
	switch s := shipping.(type) {
	case *salemongo.ShippingAddress:
		return &pbshipping_addresses.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      formatTimePtr(s.CreatedAt),
			UpdatedAt:      formatTimePtr(s.UpdatedAt),
			Courier:        s.Courier,
		}
	case *salemongo.ShippingAddressRow:
		return &pbshipping_addresses.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      formatDate(s.CreatedAt),
			UpdatedAt:      formatDate(s.UpdatedAt),
			Courier:        s.Courier,
		}
	case *salemongo.ShippingAddressListRow:
		return &pbshipping_addresses.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      formatDate(s.CreatedAt),
			UpdatedAt:      formatDate(s.UpdatedAt),
			Courier:        s.Courier,
		}
	case *repository.ShippingAddressResult:
		return &pbshipping_addresses.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      getString(s.CreatedAt),
			UpdatedAt:      getString(s.UpdatedAt),
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}

func mapToProtoShippingResponseDeleteAt(shipping interface{}) *pbshipping_addresses.ShippingResponseDeleteAt {
	switch s := shipping.(type) {
	case *salemongo.ShippingAddress:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt != nil {
			deletedAt = wrapperspb.String(s.DeletedAt.Format("2006-01-02"))
		}
		return &pbshipping_addresses.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      formatTimePtr(s.CreatedAt),
			UpdatedAt:      formatTimePtr(s.UpdatedAt),
			DeletedAt:      deletedAt,
			Courier:        s.Courier,
		}
	case *salemongo.ShippingAddressActiveRow:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt != nil {
			deletedAt = wrapperspb.String(s.DeletedAt.Format("2006-01-02"))
		}
		return &pbshipping_addresses.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      formatDate(s.CreatedAt),
			UpdatedAt:      formatDate(s.UpdatedAt),
			DeletedAt:      deletedAt,
			Courier:        s.Courier,
		}
	case *repository.ShippingAddressResult:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt != nil {
			deletedAt = wrapperspb.String(*s.DeletedAt)
		}
		return &pbshipping_addresses.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      getString(s.CreatedAt),
			UpdatedAt:      getString(s.UpdatedAt),
			DeletedAt:      deletedAt,
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}
