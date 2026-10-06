package categoryapimapper

import (
	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type categoryCommandResponseMapper struct{}

func NewCategoryCommandResponseMapper() CategoryCommandResponseMapper {
	return &categoryCommandResponseMapper{}
}

func (c *categoryCommandResponseMapper) ToResponseCategory(category *pbcategories.CategoryResponse) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:            int(category.Id),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (c *categoryCommandResponseMapper) ToResponsesCategory(categories []*pbcategories.CategoryResponse) []*response.CategoryResponse {
	var mappedCategories []*response.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategory(category))
	}
	return mappedCategories
}

func (c *categoryCommandResponseMapper) ToResponseCategoryDelete(category *pbcategories.CategoryResponseDeleteAt) *response.CategoryResponseDeleteAt {
	var deletedAt string
	if category.DeletedAt != nil {
		deletedAt = category.DeletedAt.Value
	}

	return &response.CategoryResponseDeleteAt{
		ID:            int(category.Id),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (c *categoryCommandResponseMapper) ToResponsesCategoryDeleteAt(categories []*pbcategories.CategoryResponseDeleteAt) []*response.CategoryResponseDeleteAt {
	var mappedCategories []*response.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategoryDelete(category))
	}
	return mappedCategories
}

func (c *categoryCommandResponseMapper) ToApiResponseCategory(pbResponse *pbcategories.ApiResponseCategory) *response.ApiResponseCategory {
	return &response.ApiResponseCategory{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategory(pbResponse.Data),
	}
}

func (c *categoryCommandResponseMapper) ToApiResponseCategoryDeleteAt(pbResponse *pbcategories.ApiResponseCategoryDeleteAt) *response.ApiResponseCategoryDeleteAt {
	return &response.ApiResponseCategoryDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryDelete(pbResponse.Data),
	}
}

func (c *categoryCommandResponseMapper) ToApiResponseCategoryDelete(pbResponse *pbcategories.ApiResponseCategoryDelete) *response.ApiResponseCategoryDelete {
	return &response.ApiResponseCategoryDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *categoryCommandResponseMapper) ToApiResponseCategoryAll(pbResponse *pbcategories.ApiResponseCategoryAll) *response.ApiResponseCategoryAll {
	return &response.ApiResponseCategoryAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *categoryCommandResponseMapper) ToApiResponsePaginationCategoryDeleteAt(pbResponse *pbcategories.ApiResponsePaginationCategoryDeleteAt) *response.ApiResponsePaginationCategoryDeleteAt {
	return &response.ApiResponsePaginationCategoryDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCategoryDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
