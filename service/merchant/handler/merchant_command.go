package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantCommandHandler struct {
	pbmerchants.UnimplementedMerchantCommandServiceServer
	merchantCommand service.MerchantCommandService
	logger          logger.LoggerInterface
}

func NewMerchantCommandHandler(svc service.MerchantCommandService, logger logger.LoggerInterface) pbmerchants.MerchantCommandServiceServer {
	return &merchantCommandHandler{
		merchantCommand: svc,
		logger:          logger,
	}
}

func (s *merchantCommandHandler) Create(ctx context.Context, request *pbmerchants.CreateMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	req := &requests.CreateMerchantRequest{
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    mapToProtoMerchantResponse(merchant),
	}, nil
}

func (s *merchantCommandHandler) Update(ctx context.Context, request *pbmerchants.UpdateMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	id := int(request.GetMerchantId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &id,
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    mapToProtoMerchantResponse(merchant),
	}, nil
}

func (s *merchantCommandHandler) UpdateStatus(ctx context.Context, request *pbmerchants.UpdateMerchantStatusRequest) (*pbmerchants.ApiResponseMerchant, error) {
	id := int(request.GetMerchantId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	req := &requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant status",
		Data:    mapToProtoMerchantResponse(merchant),
	}, nil
}

func (s *merchantCommandHandler) TrashedMerchant(ctx context.Context, request *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    mapToProtoMerchantResponseDeleteAt(merchant),
	}, nil
}

func (s *merchantCommandHandler) RestoreMerchant(ctx context.Context, req *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	id := int(req.GetId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    mapToProtoMerchantResponse(merchant),
	}, nil
}

func (s *merchantCommandHandler) DeleteMerchantPermanent(ctx context.Context, request *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	_, err := s.merchantCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantCommandHandler) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	_, err := s.merchantCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all merchants",
	}, nil
}

func (s *merchantCommandHandler) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	_, err := s.merchantCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchants permanently",
	}, nil
}
