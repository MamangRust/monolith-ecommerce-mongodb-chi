package merchantapimapper

import (
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBaseResponseMapper interface {
	ToResponseMerchant(merchant *pbmerchants.MerchantResponse) *response.MerchantResponse
	ToResponsesMerchant(merchants []*pbmerchants.MerchantResponse) []*response.MerchantResponse
}

type MerchantQueryResponseMapper interface {
	MerchantBaseResponseMapper
	ToApiResponseMerchant(pbResponse *pbmerchants.ApiResponseMerchant) *response.ApiResponseMerchant
	ToApiResponsesMerchant(pbResponse *pbmerchants.ApiResponsesMerchant) *response.ApiResponsesMerchant
	ToApiResponsePaginationMerchant(pbResponse *pbmerchants.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pbmerchants.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
}

type MerchantCommandResponseMapper interface {
	MerchantBaseResponseMapper
	ToApiResponseMerchant(pbResponse *pbmerchants.ApiResponseMerchant) *response.ApiResponseMerchant
	ToResponseMerchantDeleteAt(merchant *pbmerchants.MerchantResponseDeleteAt) *response.MerchantResponseDeleteAt
	ToResponsesMerchantDeleteAt(merchants []*pbmerchants.MerchantResponseDeleteAt) []*response.MerchantResponseDeleteAt
	ToApiResponseMerchantDeleteAt(pbResponse *pbmerchants.ApiResponseMerchantDeleteAt) *response.ApiResponseMerchantDeleteAt
	ToApiResponseMerchantDelete(pbResponse *pbmerchants.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete
	ToApiResponseMerchantAll(pbResponse *pbmerchants.ApiResponseMerchantAll) *response.ApiResponseMerchantAll
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pbmerchants.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
}
