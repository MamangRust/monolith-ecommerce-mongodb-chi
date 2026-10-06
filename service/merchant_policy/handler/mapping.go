package handler

import (
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchants_policy "github.com/MamangRust/monolith-ecommerce-pb/merchants_policy"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapToSingleResponse(data interface{}) *pbmerchants_policy.ApiResponseMerchantPolicies {
	return &pbmerchants_policy.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data:    mapToMerchantPolicyResponse(data).(*pbmerchants_policy.MerchantPoliciesResponse),
	}
}

func mapToPaginationResponse(data []*salemongo.MerchantPolicyListRow, total *int) *pbmerchants_policy.ApiResponsePaginationMerchantPolicies {
	var policies []*pbmerchants_policy.MerchantPoliciesResponse
	for _, v := range data {
		policies = append(policies, mapToMerchantPolicyResponse(v).(*pbmerchants_policy.MerchantPoliciesResponse))
	}

	return &pbmerchants_policy.ApiResponsePaginationMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pbcommon.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToPaginationDeleteAtResponse(data interface{}, total *int) *pbmerchants_policy.ApiResponsePaginationMerchantPoliciesDeleteAt {
	var policies []*pbmerchants_policy.MerchantPoliciesResponseDeleteAt

	switch v := data.(type) {
	case []*salemongo.MerchantPolicyActiveRow:
		for _, item := range v {
			policies = append(policies, mapToMerchantPolicyResponse(item).(*pbmerchants_policy.MerchantPoliciesResponseDeleteAt))
		}
	}

	return &pbmerchants_policy.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pbcommon.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToSingleDeleteAtResponse(data *salemongo.MerchantPolicy) *pbmerchants_policy.ApiResponseMerchantPoliciesDeleteAt {
	return &pbmerchants_policy.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully processed merchant policy",
		Data:    mapToMerchantPolicyResponse(data).(*pbmerchants_policy.MerchantPoliciesResponseDeleteAt),
	}
}

func mapToMerchantPolicyResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *salemongo.MerchantPolicyRow:
		return &pbmerchants_policy.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.String(),
			UpdatedAt:   v.UpdatedAt.String(),
		}
	case *salemongo.MerchantPolicyListRow:
		return &pbmerchants_policy.MerchantPoliciesResponse{
			Id:           int32(v.MerchantPolicyID),
			MerchantId:   int32(v.MerchantID),
			PolicyType:   v.PolicyType,
			Title:        v.Title,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt.String(),
			UpdatedAt:    v.UpdatedAt.String(),
			MerchantName: v.MerchantName,
		}
	case *salemongo.MerchantPolicyActiveRow:
		res := &pbmerchants_policy.MerchantPoliciesResponseDeleteAt{
			Id:           int32(v.MerchantPolicyID),
			MerchantId:   int32(v.MerchantID),
			PolicyType:   v.PolicyType,
			Title:        v.Title,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt.String(),
			UpdatedAt:    v.UpdatedAt.String(),
			MerchantName: v.MerchantName,
		}
		if v.DeletedAt != nil {
			res.DeletedAt = &wrapperspb.StringValue{Value: v.DeletedAt.Format("2006-01-02 15:04:05.000")}
		}
		return res
	case *salemongo.MerchantPolicy:
		res := &pbmerchants_policy.MerchantPoliciesResponseDeleteAt{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.String(),
			UpdatedAt:   v.UpdatedAt.String(),
		}
		if v.DeletedAt != nil {
			res.DeletedAt = &wrapperspb.StringValue{Value: v.DeletedAt.Format("2006-01-02 15:04:05.000")}
		}
		return res
	default:
		return nil
	}
}
