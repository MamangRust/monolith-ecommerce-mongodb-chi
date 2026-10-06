package roleapimapper

import (
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type roleCommandResponseMapper struct {
}

func NewRoleCommandResponseMapper() RoleCommandResponseMapper {
	return &roleCommandResponseMapper{}
}

func (s *roleCommandResponseMapper) ToApiResponseRole(pbResponse *pbroles.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRole(pbResponse.Data),
	}
}

func (s *roleCommandResponseMapper) ToApiResponseRoleDelete(pbResponse *pbroles.ApiResponseRoleDelete) *response.ApiResponseRoleDelete {
	return &response.ApiResponseRoleDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleCommandResponseMapper) ToApiResponseRoleAll(pbResponse *pbroles.ApiResponseRoleAll) *response.ApiResponseRoleAll {
	return &response.ApiResponseRoleAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleCommandResponseMapper) mapResponseRole(role *pbroles.RoleResponse) *response.RoleResponse {
	if role == nil {
		return nil
	}
	return &response.RoleResponse{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
