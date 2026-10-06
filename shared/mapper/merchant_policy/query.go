package merchantpolicyapimapper

import (
	pbmerchants_policy "github.com/MamangRust/monolith-ecommerce-pb/merchants_policy"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantPolicyQueryResponseMapper struct{}

func NewMerchantPolicyQueryResponseMapper() MerchantPolicyQueryResponseMapper {
	return &merchantPolicyQueryResponseMapper{}
}

func (m *merchantPolicyQueryResponseMapper) ToResponseMerchantPolicy(merchant *pbmerchants_policy.MerchantPoliciesResponse) *response.MerchantPoliciesResponse {
	if merchant == nil {
		return nil
	}
	return &response.MerchantPoliciesResponse{
		ID:           int(merchant.Id),
		MerchantID:   int(merchant.MerchantId),
		PolicyType:   merchant.PolicyType,
		Title:        merchant.Title,
		Description:  merchant.Description,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		MerchantName: merchant.MerchantName,
	}
}

func (m *merchantPolicyQueryResponseMapper) ToResponsesMerchantPolicy(merchants []*pbmerchants_policy.MerchantPoliciesResponse) []*response.MerchantPoliciesResponse {
	var mappedMerchants []*response.MerchantPoliciesResponse
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantPolicy(merchant))
	}
	return mappedMerchants
}

func (m *merchantPolicyQueryResponseMapper) ToApiResponseMerchantPolicies(pbResponse *pbmerchants_policy.ApiResponseMerchantPolicies) *response.ApiResponseMerchantPolicies {
	return &response.ApiResponseMerchantPolicies{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantPolicy(pbResponse.Data),
	}
}

func (m *merchantPolicyQueryResponseMapper) ToApiResponsesMerchantPolicies(pbResponse *pbmerchants_policy.ApiResponsesMerchantPolicies) *response.ApiResponsesMerchantPolicies {
	return &response.ApiResponsesMerchantPolicies{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantPolicy(pbResponse.Data),
	}
}

func (m *merchantPolicyQueryResponseMapper) ToApiResponsePaginationMerchantPolicies(pbResponse *pbmerchants_policy.ApiResponsePaginationMerchantPolicies) *response.ApiResponsePaginationMerchantPolicies {
	return &response.ApiResponsePaginationMerchantPolicies{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantPolicy(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantPolicyQueryResponseMapper) ToApiResponsePaginationMerchantPoliciesDeleteAt(pbResponse *pbmerchants_policy.ApiResponsePaginationMerchantPoliciesDeleteAt) *response.ApiResponsePaginationMerchantPoliciesDeleteAt {
	var data []*response.MerchantPoliciesResponseDeleteAt
	for _, b := range pbResponse.Data {
		var deletedAt string
		if b.DeletedAt != nil {
			deletedAt = b.DeletedAt.Value
		}
		data = append(data, &response.MerchantPoliciesResponseDeleteAt{
			ID:           int(b.Id),
			MerchantID:   int(b.MerchantId),
			PolicyType:   b.PolicyType,
			Title:        b.Title,
			Description:  b.Description,
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
			MerchantName: b.MerchantName,
			DeletedAt:    &deletedAt,
		})
	}
	return &response.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       data,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
