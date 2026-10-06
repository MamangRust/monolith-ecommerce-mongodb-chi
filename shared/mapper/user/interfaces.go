package userapimapper

import (
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserBaseResponseMapper interface {
	// Converts a single user response into an API response.
	ToApiResponseUser(pbResponse *pbusers.ApiResponseUser) *response.ApiResponseUser
}

type UserQueryResponseMapper interface {
	UserBaseResponseMapper

	// Converts paginated user records into an API response.
	ToApiResponsePaginationUser(pbResponse *pbusers.ApiResponsePaginationUser) *response.ApiResponsePaginationUser

	// Converts paginated soft-deleted users into an API response.
	ToApiResponsePaginationUserDeleteAt(pbResponse *pbusers.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt
}

type UserCommandResponseMapper interface {
	UserBaseResponseMapper

	// Converts a soft-deleted user response into an API response.
	ToApiResponseUserDeleteAt(pbResponse *pbusers.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt

	// Converts a permanently deleted user response into an API response.
	ToApiResponseUserDelete(pbResponse *pbusers.ApiResponseUserDelete) *response.ApiResponseUserDelete

	// Converts all user records into an API response.
	ToApiResponseUserAll(pbResponse *pbusers.ApiResponseUserAll) *response.ApiResponseUserAll
}
