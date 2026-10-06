package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_award "github.com/MamangRust/monolith-ecommerce-pb/merchants_award"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantAwardCommandHandler struct {
	pbmerchants_award.UnimplementedMerchantAwardCommandServiceServer
	merchantAwardCommand service.MerchantAwardCommandService
	logger               logger.LoggerInterface
}

func NewMerchantAwardCommandHandler(svc service.MerchantAwardCommandService, logger logger.LoggerInterface) MerchantAwardCommandHandler {
	return &merchantAwardCommandHandler{
		merchantAwardCommand: svc,
		logger:               logger,
	}
}

func (s *merchantAwardCommandHandler) Create(ctx context.Context, request *pbmerchants_award.CreateMerchantAwardRequest) (*pbmerchants_award.ApiResponseMerchantAward, error) {
	req := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     int(request.GetMerchantId()),
		Title:          request.GetTitle(),
		Description:    request.GetDescription(),
		IssuedBy:       request.GetIssuedBy(),
		CertificateUrl: request.GetCertificateUrl(),
		IssueDate:      request.GetIssueDate(),
		ExpiryDate:     request.GetExpiryDate(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantaward_errors.ErrGrpcValidateCreateMerchantAward
	}

	merchant, err := s.merchantAwardCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_award.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully created merchant award",
		Data:    mapToProtoMerchantAwardResponse(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) Update(ctx context.Context, request *pbmerchants_award.UpdateMerchantAwardRequest) (*pbmerchants_award.ApiResponseMerchantAward, error) {
	id := int(request.GetMerchantCertificationId())
	req := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &id,
		Title:                   request.GetTitle(),
		Description:             request.GetDescription(),
		IssuedBy:                request.GetIssuedBy(),
		CertificateUrl:          request.GetCertificateUrl(),
		IssueDate:               request.GetIssueDate(),
		ExpiryDate:              request.GetExpiryDate(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantaward_errors.ErrGrpcValidateUpdateMerchantAward
	}

	merchant, err := s.merchantAwardCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_award.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully updated merchant award",
		Data:    mapToProtoMerchantAwardResponse(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) TrashedMerchantAward(ctx context.Context, request *pbmerchants_award.FindByIdMerchantAwardRequest) (*pbmerchants_award.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_award.ApiResponseMerchantAwardDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant award",
		Data:    mapToProtoMerchantAwardResponseDeleteAt(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) RestoreMerchantAward(ctx context.Context, request *pbmerchants_award.FindByIdMerchantAwardRequest) (*pbmerchants_award.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants_award.ApiResponseMerchantAwardDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant award",
		Data:    mapToProtoMerchantAwardResponseDeleteAt(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) DeleteMerchantAwardPermanent(ctx context.Context, request *pbmerchants_award.FindByIdMerchantAwardRequest) (*pbmerchants.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	_, err := s.merchantAwardCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant award permanently",
	}, nil
}

func (s *merchantAwardCommandHandler) RestoreAllMerchantAward(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all trashed merchant awards",
	}, nil
}

func (s *merchantAwardCommandHandler) DeleteAllMerchantAwardPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant awards permanently",
	}, nil
}
