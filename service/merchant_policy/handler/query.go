package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_policy "github.com/MamangRust/monolith-ecommerce-pb/merchants_policy"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type merchantPolicyQueryHandler struct {
	pbmerchants_policy.UnimplementedMerchantPolicyQueryServiceServer
	merchantPolicyService service.MerchantPoliciesQueryService
	logger                logger.LoggerInterface
}

func NewMerchantPolicyQueryHandler(
	merchantPolicyService service.MerchantPoliciesQueryService,
	logger logger.LoggerInterface,
) pbmerchants_policy.MerchantPolicyQueryServiceServer {
	return &merchantPolicyQueryHandler{
		merchantPolicyService: merchantPolicyService,
		logger:                logger,
	}
}

func (h *merchantPolicyQueryHandler) FindAll(ctx context.Context, req *pbmerchants.FindAllMerchantRequest) (*pbmerchants_policy.ApiResponsePaginationMerchantPolicies, error) {
	merchants, total, err := h.merchantPolicyService.FindAll(ctx, &requests.FindAllMerchant{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	})

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToPaginationResponse(merchants, total), nil
}

func (h *merchantPolicyQueryHandler) FindById(ctx context.Context, req *pbmerchants_policy.FindByIdMerchantPoliciesRequest) (*pbmerchants_policy.ApiResponseMerchantPolicies, error) {
	merchant, err := h.merchantPolicyService.FindByID(ctx, int(req.GetId()))

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToSingleResponse(merchant), nil
}

func (h *merchantPolicyQueryHandler) FindByActive(ctx context.Context, req *pbmerchants.FindAllMerchantRequest) (*pbmerchants_policy.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
	merchants, total, err := h.merchantPolicyService.FindActive(ctx, &requests.FindAllMerchant{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	})

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToPaginationDeleteAtResponse(merchants, total), nil
}

func (h *merchantPolicyQueryHandler) FindByTrashed(ctx context.Context, req *pbmerchants.FindAllMerchantRequest) (*pbmerchants_policy.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
	merchants, total, err := h.merchantPolicyService.FindTrashed(ctx, &requests.FindAllMerchant{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	})

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToPaginationDeleteAtResponse(merchants, total), nil
}
