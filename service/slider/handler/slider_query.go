package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/service"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbsliders "github.com/MamangRust/monolith-ecommerce-pb/sliders"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type sliderQueryHandler struct {
	pbsliders.UnimplementedSliderQueryServiceServer
	sliderQuery service.SliderQueryService
	logger      logger.LoggerInterface
}

func NewSliderQueryHandler(sliderQuery service.SliderQueryService, logger logger.LoggerInterface) *sliderQueryHandler {
	return &sliderQueryHandler{
		sliderQuery: sliderQuery,
		logger:      logger,
	}
}

func (s *sliderQueryHandler) FindAll(ctx context.Context, request *pbsliders.FindAllSliderRequest) (*pbsliders.ApiResponsePaginationSlider, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pbsliders.SliderResponse, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = MapToSliderResponseGetSlidersRow(slider)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pbsliders.ApiResponsePaginationSlider{
		Status:     "success",
		Message:    "Successfully fetched slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderQueryHandler) FindByActive(ctx context.Context, request *pbsliders.FindAllSliderRequest) (*pbsliders.ApiResponsePaginationSliderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pbsliders.SliderResponseDeleteAt, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = MapToSliderResponseDeleteAtGetSlidersActiveRow(slider)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pbsliders.ApiResponsePaginationSliderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderQueryHandler) FindByTrashed(ctx context.Context, request *pbsliders.FindAllSliderRequest) (*pbsliders.ApiResponsePaginationSliderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pbsliders.SliderResponseDeleteAt, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = MapToSliderResponseDeleteAtGetSlidersTrashedRow(slider)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pbsliders.ApiResponsePaginationSliderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderQueryHandler) FindById(ctx context.Context, request *pbsliders.FindByIdSliderRequest) (*pbsliders.ApiResponseSlider, error) {
	id := int(request.GetId())

	slider, err := s.sliderQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbsliders.ApiResponseSlider{
		Status:  "success",
		Message: "Successfully fetched slider by ID",
		Data:    MapToSliderResponseGetSliderByIDRow(slider),
	}, nil
}
