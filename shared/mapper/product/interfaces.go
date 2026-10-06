package productapimapper

import (
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductBaseResponseMapper interface {
	ToResponseProduct(product *pbproducts.ProductResponse) *response.ProductResponse
	ToResponsesProduct(products []*pbproducts.ProductResponse) []*response.ProductResponse
	ToResponseProductDeleteAt(product *pbproducts.ProductResponseDeleteAt) *response.ProductResponseDeleteAt
	ToResponsesProductDeleteAt(products []*pbproducts.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt
	ToApiResponseProduct(pbResponse *pbproducts.ApiResponseProduct) *response.ApiResponseProduct
	ToApiResponsePaginationProductDeleteAt(pbResponse *pbproducts.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt
}

type ProductQueryResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProduct(pbResponse *pbproducts.ApiResponsesProduct) *response.ApiResponsesProduct
	ToApiResponsePaginationProduct(pbResponse *pbproducts.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct
}

type ProductCommandResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProductDeleteAt(pbResponse *pbproducts.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt
	ToApiResponseProductDelete(pbResponse *pbproducts.ApiResponseProductDelete) *response.ApiResponseProductDelete
	ToApiResponseProductAll(pbResponse *pbproducts.ApiResponseProductAll) *response.ApiResponseProductAll
}
