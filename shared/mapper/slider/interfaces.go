package sliderapimapper

import (
	pbsliders "github.com/MamangRust/monolith-ecommerce-pb/sliders"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderBaseResponseMapper interface {
	ToResponseSlider(pbResponse *pbsliders.SliderResponse) *response.SliderResponse
	ToResponsesSlider(pbResponses []*pbsliders.SliderResponse) []*response.SliderResponse
}

type SliderQueryResponseMapper interface {
	SliderBaseResponseMapper
	ToApiResponseSlider(pbResponse *pbsliders.ApiResponseSlider) *response.ApiResponseSlider
	ToApiResponsesSlider(pbResponse *pbsliders.ApiResponsesSlider) *response.ApiResponsesSlider
	ToApiResponsePaginationSlider(pbResponse *pbsliders.ApiResponsePaginationSlider) *response.ApiResponsePaginationSlider
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pbsliders.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
}

type SliderCommandResponseMapper interface {
	SliderBaseResponseMapper
	ToResponseSliderDeleteAt(pbResponse *pbsliders.SliderResponseDeleteAt) *response.SliderResponseDeleteAt
	ToResponsesSliderDeleteAt(pbResponses []*pbsliders.SliderResponseDeleteAt) []*response.SliderResponseDeleteAt
	ToApiResponseSliderDeleteAt(pbResponse *pbsliders.ApiResponseSliderDeleteAt) *response.ApiResponseSliderDeleteAt
	ToApiResponseSliderDelete(pbResponse *pbsliders.ApiResponseSliderDelete) *response.ApiResponseSliderDelete
	ToApiResponseSliderAll(pbResponse *pbsliders.ApiResponseSliderAll) *response.ApiResponseSliderAll
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pbsliders.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
}
