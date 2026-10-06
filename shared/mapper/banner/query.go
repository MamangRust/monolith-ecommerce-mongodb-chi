package bannerapimapper

import (
	pbbanners "github.com/MamangRust/monolith-ecommerce-pb/banners"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type bannerQueryResponseMapper struct{}

func NewBannerQueryResponseMapper() BannerQueryResponseMapper {
	return &bannerQueryResponseMapper{}
}

func (m *bannerQueryResponseMapper) ToResponseBanner(banner *pbbanners.BannerResponse) *response.BannerResponse {
	if banner == nil {
		return nil
	}
	return &response.BannerResponse{
		ID:        banner.BannerId,
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
	}
}

func (m *bannerQueryResponseMapper) ToResponsesBanner(banners []*pbbanners.BannerResponse) []*response.BannerResponse {
	var mappedBanners []*response.BannerResponse
	for _, banner := range banners {
		mappedBanners = append(mappedBanners, m.ToResponseBanner(banner))
	}
	return mappedBanners
}

func (m *bannerQueryResponseMapper) ToApiResponseBanner(pbResponse *pbbanners.ApiResponseBanner) *response.ApiResponseBanner {
	return &response.ApiResponseBanner{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBanner(pbResponse.Data),
	}
}

func (m *bannerQueryResponseMapper) ToApiResponsesBanner(pbResponse *pbbanners.ApiResponsesBanner) *response.ApiResponsesBanner {
	return &response.ApiResponsesBanner{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesBanner(pbResponse.Data),
	}
}

func (m *bannerQueryResponseMapper) ToApiResponsePaginationBanner(pbResponse *pbbanners.ApiResponsePaginationBanner) *response.ApiResponsePaginationBanner {
	return &response.ApiResponsePaginationBanner{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesBanner(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *bannerQueryResponseMapper) ToApiResponsePaginationBannerDeleteAt(pbResponse *pbbanners.ApiResponsePaginationBannerDeleteAt) *response.ApiResponsePaginationBannerDeleteAt {
	var data []*response.BannerResponseDeleteAt
	for _, b := range pbResponse.Data {
		var deletedAt string
		if b.DeletedAt != nil {
			deletedAt = b.DeletedAt.Value
		}
		data = append(data, &response.BannerResponseDeleteAt{
			ID:        b.BannerId,
			Name:      b.Name,
			StartDate: b.StartDate,
			EndDate:   b.EndDate,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
			IsActive:  b.IsActive,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
			DeletedAt: &deletedAt,
		})
	}
	return &response.ApiResponsePaginationBannerDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       data,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
