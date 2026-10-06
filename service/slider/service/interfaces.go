package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type SliderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, error)
	FindByID(ctx context.Context, slider_id int) (*repository.SliderResult, error)
}

type SliderCommandService interface {
	Create(ctx context.Context, request *requests.CreateSliderRequest) (*salemongo.Slider, error)
	Update(ctx context.Context, request *requests.UpdateSliderRequest) (*salemongo.Slider, error)
	Trash(ctx context.Context, slider_id int) (*salemongo.Slider, error)
	Restore(ctx context.Context, slider_id int) (*salemongo.Slider, error)
	DeletePermanent(ctx context.Context, slider_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
