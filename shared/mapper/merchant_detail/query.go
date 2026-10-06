package merchantdetailapimapper

import (
	pbmerchants_detail "github.com/MamangRust/monolith-ecommerce-pb/merchants_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantDetailQueryResponseMapper struct {
	MerchantDetailCommandResponseMapper
}

func NewMerchantDetailQueryResponseMapper() MerchantDetailQueryResponseMapper {
	return &merchantDetailQueryResponseMapper{
		MerchantDetailCommandResponseMapper: NewMerchantDetailCommandResponseMapper(),
	}
}

func (m *merchantDetailQueryResponseMapper) ToResponseMerchantDetail(merchant *pbmerchants_detail.MerchantDetailResponse) *response.MerchantDetailResponse {
	return m.MerchantDetailCommandResponseMapper.ToResponseMerchantDetail(merchant)
}

func (m *merchantDetailQueryResponseMapper) ToResponseMerchantDetailRelation(merchant *pbmerchants_detail.MerchantDetailResponse) *response.MerchantDetailResponse {
	return m.MerchantDetailCommandResponseMapper.ToResponseMerchantDetailRelation(merchant)
}

func (m *merchantDetailQueryResponseMapper) ToResponsesMerchantDetail(merchants []*pbmerchants_detail.MerchantDetailResponse) []*response.MerchantDetailResponse {
	return m.MerchantDetailCommandResponseMapper.ToResponsesMerchantDetail(merchants)
}

func (m *merchantDetailQueryResponseMapper) ToApiResponseMerchantDetail(pbResponse *pbmerchants_detail.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetail {
	return m.MerchantDetailCommandResponseMapper.ToApiResponseMerchantDetail(pbResponse)
}

func (m *merchantDetailQueryResponseMapper) ToApiResponseMerchantDetailRelation(pbResponse *pbmerchants_detail.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetailRelation {
	return &response.ApiResponseMerchantDetailRelation{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetailRelation(pbResponse.Data),
	}
}

func (m *merchantDetailQueryResponseMapper) ToApiResponsesMerchantDetail(pbResponse *pbmerchants_detail.ApiResponsesMerchantDetail) *response.ApiResponsesMerchantDetail {
	return &response.ApiResponsesMerchantDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantDetail(pbResponse.Data),
	}
}

func (m *merchantDetailQueryResponseMapper) ToApiResponsePaginationMerchantDetail(pbResponse *pbmerchants_detail.ApiResponsePaginationMerchantDetail) *response.ApiResponsePaginationMerchantDetail {
	return &response.ApiResponsePaginationMerchantDetail{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantDetail(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantDetailQueryResponseMapper) ToApiResponsePaginationMerchantDetailDeleteAt(pbResponse *pbmerchants_detail.ApiResponsePaginationMerchantDetailDeleteAt) *response.ApiResponsePaginationMerchantDetailDeleteAt {
	return m.MerchantDetailCommandResponseMapper.ToApiResponsePaginationMerchantDetailDeleteAt(pbResponse)
}
