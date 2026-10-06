package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func (h *Handler) mapToCategoryResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *salemongo.CategoryRow:
		return &pbcategories.CategoryResponse{
			Id:            v.CategoryID,
			Name:          v.Name,
			Description:   derefString(v.Description),
			SlugCategory:  derefString(v.SlugCategory),
			ImageCategory: derefString(v.ImageCategory),
			CreatedAt:     v.CreatedAt.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Format("2006-01-02"),
		}
	case *salemongo.CategoryListRow:
		return &pbcategories.CategoryResponse{
			Id:            v.CategoryID,
			Name:          v.Name,
			Description:   derefString(v.Description),
			SlugCategory:  derefString(v.SlugCategory),
			ImageCategory: derefString(v.ImageCategory),
			CreatedAt:     v.CreatedAt.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Format("2006-01-02"),
		}
	case *repository.CategoryResult:
		return &pbcategories.CategoryResponse{
			Id:            v.CategoryID,
			Name:          v.Name,
			Description:   derefString(v.Description),
			SlugCategory:  derefString(v.SlugCategory),
			ImageCategory: derefString(v.ImageCategory),
			CreatedAt:     formatTimePtr(v.CreatedAt),
			UpdatedAt:     formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.CategoryActiveRow:
		var deletedAt string
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.Format("2006-01-02")
		}
		return &pbcategories.CategoryResponseDeleteAt{
			Id:            v.CategoryID,
			Name:          v.Name,
			Description:   derefString(v.Description),
			SlugCategory:  derefString(v.SlugCategory),
			ImageCategory: derefString(v.ImageCategory),
			CreatedAt:     v.CreatedAt.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Format("2006-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	case *salemongo.Category:
		return &pbcategories.CategoryResponse{
			Id:            v.CategoryID,
			Name:          v.Name,
			Description:   derefString(v.Description),
			SlugCategory:  derefString(v.SlugCategory),
			ImageCategory: derefString(v.ImageCategory),
			CreatedAt:     formatTimePtr(v.CreatedAt),
			UpdatedAt:     formatTimePtr(v.UpdatedAt),
		}
	case *salemongo.CategoryRevenueRow:
		return &pbcategories.CategoriesMonthlyTotalPriceResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *salemongo.CategoryYearlyRevenueRow:
		return &pbcategories.CategoriesYearlyTotalPriceResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *salemongo.CategoryStatsRow:
		return &pbcategories.CategoryMonthPriceResponse{
			Month:        v.Month,
			CategoryId:   v.CategoryID,
			CategoryName: v.CategoryName,
			OrderCount:   int32(v.OrderCount),
			ItemsSold:    int32(v.ItemsSold),
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *salemongo.CategoryYearlyStatsRow:
		return &pbcategories.CategoryYearPriceResponse{
			Year:               v.Year,
			CategoryId:         v.CategoryID,
			CategoryName:       v.CategoryName,
			OrderCount:         int32(v.OrderCount),
			ItemsSold:          int32(v.ItemsSold),
			TotalRevenue:       int32(v.TotalRevenue),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	default:
		log.Printf("Unknown type for mapping: %T", v)
		return nil
	}
}

func mapCategoryResultToDeleteAt(v *repository.CategoryResult) *pbcategories.CategoryResponseDeleteAt {
	var deletedAt string
	if v.DeletedAt != nil {
		deletedAt = v.DeletedAt.Format("2006-01-02")
	}
	return &pbcategories.CategoryResponseDeleteAt{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   derefString(v.Description),
		SlugCategory:  derefString(v.SlugCategory),
		ImageCategory: derefString(v.ImageCategory),
		CreatedAt:     formatTimePtr(v.CreatedAt),
		UpdatedAt:     formatTimePtr(v.UpdatedAt),
		DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
	}
}

func mapToCategoryDeleteAtResponse(v *salemongo.Category) *pbcategories.CategoryResponseDeleteAt {
	var deletedAt string
	if v.DeletedAt != nil {
		deletedAt = v.DeletedAt.Format("2006-01-02")
	}
	return &pbcategories.CategoryResponseDeleteAt{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   derefString(v.Description),
		SlugCategory:  derefString(v.SlugCategory),
		ImageCategory: derefString(v.ImageCategory),
		CreatedAt:     formatTimePtr(v.CreatedAt),
		UpdatedAt:     formatTimePtr(v.UpdatedAt),
		DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
	}
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
