package handler

import (
	"math"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ProductResult aliases repository.ProductResult so the proto mapping type switch
// sees the concrete type returned by the service layer.
type ProductResult = repository.ProductResult

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
		if t == nil {
			return ""
		}
		return convert.FormatTimeRFC3339(t)
	case string:
		return t
	}
	return ""
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func float64PtrToFloat32(f *float64) float32 {
	if f == nil {
		return 0
	}
	return float32(*f)
}

func weightPtrToInt32(w interface{}) int32 {
	switch v := w.(type) {
	case *int32:
		if v != nil {
			return *v
		}
	case int32:
		return v
	}
	return 0
}

func productResultToResponse(v *ProductResult) *pbproducts.ProductResponse {
	return &pbproducts.ProductResponse{
		Id:           v.ProductID,
		MerchantId:   v.MerchantID,
		CategoryId:   v.CategoryID,
		Name:         v.Name,
		Description:  stringPtrToString(v.Description),
		Price:        v.Price,
		CountInStock: v.CountInStock,
		Brand:        stringPtrToString(v.Brand),
		Weight:       weightPtrToInt32(v.Weight),
		Rating:       float64PtrToFloat32(v.Rating),
		SlugProduct:  stringPtrToString(v.SlugProduct),
		ImageProduct: stringPtrToString(v.ImageProduct),
		CreatedAt:    formatTimestamp(v.CreatedAt),
		UpdatedAt:    formatTimestamp(v.UpdatedAt),
	}
}

func mapToProtoProductResponse(item interface{}) *pbproducts.ProductResponse {
	switch v := item.(type) {
	case *ProductResult:
		return productResultToResponse(v)
	case *salemongo.Product:
		return &pbproducts.ProductResponse{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoProductResponseDeleteAt(item interface{}) *pbproducts.ProductResponseDeleteAt {
	var res *pbproducts.ProductResponseDeleteAt
	var deletedAt interface{}

	switch v := item.(type) {
	case *ProductResult:
		res = &pbproducts.ProductResponseDeleteAt{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *salemongo.Product:
		res = &pbproducts.ProductResponseDeleteAt{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
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
