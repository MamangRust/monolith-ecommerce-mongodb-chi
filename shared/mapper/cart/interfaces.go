package cartapimapper

import (
	pbcarts "github.com/MamangRust/monolith-ecommerce-pb/carts"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartBaseResponseMapper interface {
	ToResponseCart(pbResponse *pbcarts.CartResponse) *response.CartResponse
	ToResponseCarts(pbResponse []*pbcarts.CartResponse) []*response.CartResponse
	ToApiResponseCart(pbResponse *pbcarts.ApiResponseCart) *response.ApiResponseCart
}

type CartQueryResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartPagination(pbResponse *pbcarts.ApiResponsePaginationCart) *response.ApiResponseCartPagination
}

type CartCommandResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartDelete(pbResponse *pbcarts.ApiResponseCartDelete) *response.ApiResponseCartDelete
	ToApiResponseCartAll(pbResponse *pbcarts.ApiResponseCartAll) *response.ApiResponseCartAll
}
