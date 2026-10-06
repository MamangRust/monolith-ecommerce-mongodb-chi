package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_detail "github.com/MamangRust/monolith-ecommerce-pb/merchants_detail"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
)

type merchantDetailQueryHandler struct {
	pbmerchants_detail.UnimplementedMerchantDetailQueryServiceServer
	MerchantDetailQuery service.MerchantDetailQueryService
	logger              logger.LoggerInterface
}

func NewMerchantDetailQueryHandler(svc service.MerchantDetailQueryService, logger logger.LoggerInterface) MerchantDetailQueryHandler {
	return &merchantDetailQueryHandler{
		MerchantDetailQuery: svc,
		logger:              logger,
	}
}

func (s *merchantDetailQueryHandler) FindAll(ctx context.Context, request *pbmerchants.FindAllMerchantRequest) (*pbmerchants_detail.ApiResponsePaginationMerchantDetail, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.MerchantDetailQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoDetails := make([]*pbmerchants_detail.MerchantDetailResponse, len(details))
	for i, detail := range details {
		protoDetails[i] = mapToProtoMerchantDetailResponse(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchants_detail.ApiResponsePaginationMerchantDetail{
		Status:     "success",
		Message:    "Successfully fetched merchant details",
		Data:       protoDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailQueryHandler) FindById(ctx context.Context, request *pbmerchants_detail.FindByIdMerchantDetailRequest) (*pbmerchants_detail.ApiResponseMerchantDetail, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	detail, err := s.MerchantDetailQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_detail.ApiResponseMerchantDetail{
		Status:  "success",
		Message: "Successfully fetched merchant detail",
		Data:    mapToProtoMerchantDetailResponse(detail),
	}, nil
}

func (s *merchantDetailQueryHandler) FindByActive(ctx context.Context, request *pbmerchants.FindAllMerchantRequest) (*pbmerchants_detail.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.MerchantDetailQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoDetails := make([]*pbmerchants_detail.MerchantDetailResponseDeleteAt, len(details))
	for i, detail := range details {
		protoDetails[i] = mapToProtoMerchantDetailResponseDeleteAt(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchants_detail.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant details",
		Data:       protoDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailQueryHandler) FindByTrashed(ctx context.Context, request *pbmerchants.FindAllMerchantRequest) (*pbmerchants_detail.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.MerchantDetailQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoDetails := make([]*pbmerchants_detail.MerchantDetailResponseDeleteAt, len(details))
	for i, detail := range details {
		protoDetails[i] = mapToProtoMerchantDetailResponseDeleteAt(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchants_detail.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant details",
		Data:       protoDetails,
		Pagination: paginationMeta,
	}, nil
}
