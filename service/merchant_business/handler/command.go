package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_business "github.com/MamangRust/monolith-ecommerce-pb/merchants_business"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantBusinessCommandHandler struct {
	pbmerchants_business.UnimplementedMerchantBusinessCommandServiceServer
	merchantBusinessCommand service.MerchantBusinessCommandService
	logger                  logger.LoggerInterface
}

func NewMerchantBusinessCommandHandler(svc service.MerchantBusinessCommandService, logger logger.LoggerInterface) MerchantBusinessCommandHandler {
	return &merchantBusinessCommandHandler{
		merchantBusinessCommand: svc,
		logger:                  logger,
	}
}

func (s *merchantBusinessCommandHandler) Create(ctx context.Context, request *pbmerchants_business.CreateMerchantBusinessRequest) (*pbmerchants_business.ApiResponseMerchantBusiness, error) {
	req := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        int(request.GetMerchantId()),
		BusinessType:      request.GetBusinessType(),
		TaxID:             request.GetTaxId(),
		EstablishedYear:   int(request.GetEstablishedYear()),
		NumberOfEmployees: int(request.GetNumberOfEmployees()),
		WebsiteUrl:        request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantbusiness_errors.ErrGrpcValidateCreateMerchantBusiness
	}

	merchant, err := s.merchantBusinessCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_business.ApiResponseMerchantBusiness{
		Status:  "success",
		Message: "Successfully created merchant business",
		Data:    mapToProtoMerchantBusinessResponse(merchant),
	}, nil
}

func (s *merchantBusinessCommandHandler) Update(ctx context.Context, request *pbmerchants_business.UpdateMerchantBusinessRequest) (*pbmerchants_business.ApiResponseMerchantBusiness, error) {
	id := int(request.GetMerchantBusinessInfoId())
	req := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &id,
		BusinessType:           request.GetBusinessType(),
		TaxID:                  request.GetTaxId(),
		EstablishedYear:        int(request.GetEstablishedYear()),
		NumberOfEmployees:      int(request.GetNumberOfEmployees()),
		WebsiteUrl:             request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantbusiness_errors.ErrGrpcValidateUpdateMerchantBusiness
	}

	merchant, err := s.merchantBusinessCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_business.ApiResponseMerchantBusiness{
		Status:  "success",
		Message: "Successfully updated merchant business",
		Data:    mapToProtoMerchantBusinessResponse(merchant),
	}, nil
}

func (s *merchantBusinessCommandHandler) TrashedMerchantBusiness(ctx context.Context, request *pbmerchants_business.FindByIdMerchantBusinessRequest) (*pbmerchants_business.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_business.ApiResponseMerchantBusinessDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant business",
		Data:    mapToProtoMerchantBusinessResponseDeleteAt(merchant),
	}, nil
}

func (s *merchantBusinessCommandHandler) RestoreMerchantBusiness(ctx context.Context, request *pbmerchants_business.FindByIdMerchantBusinessRequest) (*pbmerchants_business.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_business.ApiResponseMerchantBusinessDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant business",
		Data:    mapToProtoMerchantBusinessResponseDeleteAt(merchant),
	}, nil
}

func (s *merchantBusinessCommandHandler) DeleteMerchantBusinessPermanent(ctx context.Context, request *pbmerchants_business.FindByIdMerchantBusinessRequest) (*pbmerchants.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	_, err := s.merchantBusinessCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant business permanently",
	}, nil
}

func (s *merchantBusinessCommandHandler) RestoreAllMerchantBusiness(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	_, err := s.merchantBusinessCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all trashed merchant businesses",
	}, nil
}

func (s *merchantBusinessCommandHandler) DeleteAllMerchantBusinessPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	_, err := s.merchantBusinessCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant businesses permanently",
	}, nil
}
