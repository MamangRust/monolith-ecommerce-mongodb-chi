package handler

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	pbsliders "github.com/MamangRust/monolith-ecommerce-pb/sliders"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func fmtTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func fmtTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func MapToSliderResponse(slider *salemongo.Slider) *pbsliders.SliderResponse {
	return &pbsliders.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: fmtTimePtr(slider.CreatedAt),
		UpdatedAt: fmtTimePtr(slider.UpdatedAt),
	}
}

func MapToSliderResponseGetSlidersRow(slider *repository.SliderResult) *pbsliders.SliderResponse {
	createdAt, updatedAt := "", ""
	if slider.CreatedAt != nil {
		createdAt = *slider.CreatedAt
	}
	if slider.UpdatedAt != nil {
		updatedAt = *slider.UpdatedAt
	}
	return &pbsliders.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func MapToSliderResponseGetSliderByIDRow(slider *repository.SliderResult) *pbsliders.SliderResponse {
	return &pbsliders.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: derefStr(slider.CreatedAt),
		UpdatedAt: derefStr(slider.UpdatedAt),
	}
}

func MapToSliderResponseCreateSliderRow(slider *salemongo.Slider) *pbsliders.SliderResponse {
	return &pbsliders.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: fmtTimePtr(slider.CreatedAt),
		UpdatedAt: fmtTimePtr(slider.UpdatedAt),
	}
}

func MapToSliderResponseUpdateSliderRow(slider *salemongo.Slider) *pbsliders.SliderResponse {
	return &pbsliders.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: fmtTimePtr(slider.CreatedAt),
		UpdatedAt: fmtTimePtr(slider.UpdatedAt),
	}
}

func MapToSliderResponseDeleteAt(slider *salemongo.Slider) *pbsliders.SliderResponseDeleteAt {
	return &pbsliders.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: fmtTimePtr(slider.CreatedAt),
		UpdatedAt: fmtTimePtr(slider.UpdatedAt),
		DeletedAt: deletedAtWrappers(slider.DeletedAt),
	}
}

func MapToSliderResponseDeleteAtGetSlidersActiveRow(slider *repository.SliderResult) *pbsliders.SliderResponseDeleteAt {
	createdAt, updatedAt, deletedAt := "", "", ""
	if slider.CreatedAt != nil {
		createdAt = *slider.CreatedAt
	}
	if slider.UpdatedAt != nil {
		updatedAt = *slider.UpdatedAt
	}
	if slider.DeletedAt != nil {
		deletedAt = *slider.DeletedAt
	}
	return &pbsliders.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}
}

func MapToSliderResponseDeleteAtGetSlidersTrashedRow(slider *repository.SliderResult) *pbsliders.SliderResponseDeleteAt {
	return MapToSliderResponseDeleteAtGetSlidersActiveRow(slider)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func deletedAtWrappers(t *time.Time) *wrapperspb.StringValue {
	if t == nil {
		return &wrapperspb.StringValue{Value: ""}
	}
	return &wrapperspb.StringValue{Value: t.Format("2006-01-02")}
}
