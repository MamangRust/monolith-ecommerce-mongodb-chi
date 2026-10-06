package merchantawardapimapper

import (
	pbmerchants_award "github.com/MamangRust/monolith-ecommerce-pb/merchants_award"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantAwardCommandResponseMapper struct{}

func NewMerchantAwardCommandResponseMapper() MerchantAwardCommandResponseMapper {
	return &merchantAwardCommandResponseMapper{}
}

func (m *merchantAwardCommandResponseMapper) ToResponseMerchantAward(MerchantAward *pbmerchants_award.MerchantAwardResponse) *response.MerchantAwardResponse {
	if MerchantAward == nil {
		return nil
	}
	return &response.MerchantAwardResponse{
		ID:             int(MerchantAward.Id),
		MerchantID:     int(MerchantAward.MerchantId),
		Title:          MerchantAward.Title,
		Description:    MerchantAward.Description,
		IssuedBy:       MerchantAward.IssuedBy,
		IssueDate:      MerchantAward.IssueDate,
		ExpiryDate:     MerchantAward.ExpiryDate,
		CertificateUrl: MerchantAward.CertificateUrl,
		CreatedAt:      MerchantAward.CreatedAt,
		UpdatedAt:      MerchantAward.UpdatedAt,
	}
}

func (m *merchantAwardCommandResponseMapper) ToResponsesMerchantAward(MerchantAwards []*pbmerchants_award.MerchantAwardResponse) []*response.MerchantAwardResponse {
	var mappedMerchantAwards []*response.MerchantAwardResponse
	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAward(MerchantAward))
	}
	return mappedMerchantAwards
}

func (m *merchantAwardCommandResponseMapper) ToApiResponseMerchantAward(pbResponse *pbmerchants_award.ApiResponseMerchantAward) *response.ApiResponseMerchantAward {
	return &response.ApiResponseMerchantAward{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAward(pbResponse.Data),
	}
}

func (m *merchantAwardCommandResponseMapper) ToResponseMerchantAwardDeleteAt(MerchantAward *pbmerchants_award.MerchantAwardResponseDeleteAt) *response.MerchantAwardResponseDeleteAt {
	if MerchantAward == nil {
		return nil
	}
	var deletedAt string
	if MerchantAward.DeletedAt != nil {
		deletedAt = MerchantAward.DeletedAt.Value
	}

	return &response.MerchantAwardResponseDeleteAt{
		ID:             int(MerchantAward.Id),
		MerchantID:     int(MerchantAward.MerchantId),
		Title:          MerchantAward.Title,
		Description:    MerchantAward.Description,
		IssuedBy:       MerchantAward.IssuedBy,
		IssueDate:      MerchantAward.IssueDate,
		ExpiryDate:     MerchantAward.ExpiryDate,
		CertificateUrl: MerchantAward.CertificateUrl,
		CreatedAt:      MerchantAward.CreatedAt,
		UpdatedAt:      MerchantAward.UpdatedAt,
		DeletedAt:      &deletedAt,
	}
}

func (m *merchantAwardCommandResponseMapper) ToResponsesMerchantAwardDeleteAt(MerchantAwards []*pbmerchants_award.MerchantAwardResponseDeleteAt) []*response.MerchantAwardResponseDeleteAt {
	var mappedMerchantAwards []*response.MerchantAwardResponseDeleteAt
	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAwardDeleteAt(MerchantAward))
	}
	return mappedMerchantAwards
}

func (m *merchantAwardCommandResponseMapper) ToApiResponseMerchantAwardDeleteAt(pbResponse *pbmerchants_award.ApiResponseMerchantAwardDeleteAt) *response.ApiResponseMerchantAwardDeleteAt {
	return &response.ApiResponseMerchantAwardDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAwardDeleteAt(pbResponse.Data),
	}
}
