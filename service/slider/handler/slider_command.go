package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/service"
	pbsliders "github.com/MamangRust/monolith-ecommerce-pb/sliders"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderCommandHandler struct {
	pbsliders.UnimplementedSliderCommandServiceServer
	sliderCommand service.SliderCommandService
	logger        logger.LoggerInterface
}

func NewSliderCommandHandler(sliderCommand service.SliderCommandService, logger logger.LoggerInterface) *sliderCommandHandler {
	return &sliderCommandHandler{
		sliderCommand: sliderCommand,
		logger:        logger,
	}
}

func (s *sliderCommandHandler) Create(ctx context.Context, request *pbsliders.CreateSliderRequest) (*pbsliders.ApiResponseSlider, error) {
	req := &requests.CreateSliderRequest{
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, slider_errors.ErrGrpcValidateCreateSlider
	}

	slider, err := s.sliderCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSlider{
		Status:  "success",
		Message: "Successfully created slider",
		Data:    MapToSliderResponseCreateSliderRow(slider),
	}, nil
}

func (s *sliderCommandHandler) Update(ctx context.Context, request *pbsliders.UpdateSliderRequest) (*pbsliders.ApiResponseSlider, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateSliderRequest{
		ID:       &id,
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, slider_errors.ErrGrpcValidateUpdateSlider
	}

	slider, err := s.sliderCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSlider{
		Status:  "success",
		Message: "Successfully updated slider",
		Data:    MapToSliderResponseUpdateSliderRow(slider),
	}, nil
}

func (s *sliderCommandHandler) TrashedSlider(ctx context.Context, request *pbsliders.FindByIdSliderRequest) (*pbsliders.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	slider, err := s.sliderCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSliderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed slider",
		Data:    MapToSliderResponseDeleteAt(slider),
	}, nil
}

func (s *sliderCommandHandler) RestoreSlider(ctx context.Context, request *pbsliders.FindByIdSliderRequest) (*pbsliders.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	slider, err := s.sliderCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSliderDeleteAt{
		Status:  "success",
		Message: "Successfully restored slider",
		Data:    MapToSliderResponseDeleteAt(slider),
	}, nil
}

func (s *sliderCommandHandler) DeleteSliderPermanent(ctx context.Context, request *pbsliders.FindByIdSliderRequest) (*pbsliders.ApiResponseSliderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	_, err := s.sliderCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSliderDelete{
		Status:  "success",
		Message: "Successfully deleted slider permanently",
	}, nil
}

func (s *sliderCommandHandler) RestoreAllSlider(ctx context.Context, _ *emptypb.Empty) (*pbsliders.ApiResponseSliderAll, error) {
	_, err := s.sliderCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSliderAll{
		Status:  "success",
		Message: "Successfully restored all sliders",
	}, nil
}

func (s *sliderCommandHandler) DeleteAllSliderPermanent(ctx context.Context, _ *emptypb.Empty) (*pbsliders.ApiResponseSliderAll, error) {
	_, err := s.sliderCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSliderAll{
		Status:  "success",
		Message: "Successfully deleted all sliders permanently",
	}, nil
}
