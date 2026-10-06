package repository

import (
	"context"
	"errors"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sliderCommandRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewSliderCommandRepository(db *mongo.Database, coll *mongo.Collection) *sliderCommandRepository {
	return &sliderCommandRepository{db: db, coll: coll}
}

func (r *sliderCommandRepository) Create(ctx context.Context, request *requests.CreateSliderRequest) (*salemongo.Slider, error) {
	sliderID, err := salemongo.NextID(ctx, r.db, "sliders")
	if err != nil {
		return nil, slider_errors.ErrCreateSlider.WithInternal(err)
	}
	now := time.Now()
	slider := &salemongo.Slider{SliderID: sliderID, Name: request.Nama, Image: request.FilePath, CreatedAt: &now, UpdatedAt: &now}
	if _, err := r.coll.InsertOne(ctx, slider); err != nil {
		return nil, slider_errors.ErrCreateSlider.WithInternal(err)
	}
	return slider, nil
}

func (r *sliderCommandRepository) Update(ctx context.Context, request *requests.UpdateSliderRequest) (*salemongo.Slider, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(ctx, bson.M{"slider_id": *request.ID, "deleted_at": nil}, bson.M{"$set": bson.M{"name": request.Nama, "image": request.FilePath, "updated_at": now}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	var slider salemongo.Slider
	if err := result.Decode(&slider); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, slider_errors.ErrSliderNotFound
		}
		return nil, slider_errors.ErrUpdateSlider.WithInternal(err)
	}
	return &slider, nil
}

func (r *sliderCommandRepository) Trash(ctx context.Context, id int) (*salemongo.Slider, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(ctx, bson.M{"slider_id": id, "deleted_at": nil}, bson.M{"$set": bson.M{"deleted_at": now}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	var slider salemongo.Slider
	if err := result.Decode(&slider); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, slider_errors.ErrSliderNotFound
		}
		return nil, slider_errors.ErrTrashSlider.WithInternal(err)
	}
	return &slider, nil
}

func (r *sliderCommandRepository) Restore(ctx context.Context, id int) (*salemongo.Slider, error) {
	result := r.coll.FindOneAndUpdate(ctx, bson.M{"slider_id": id, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	var slider salemongo.Slider
	if err := result.Decode(&slider); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, slider_errors.ErrSliderNotFound
		}
		return nil, slider_errors.ErrRestoreSlider.WithInternal(err)
	}
	return &slider, nil
}

func (r *sliderCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"slider_id": id, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, slider_errors.ErrSliderNotFound
	}
	return true, nil
}

func (r *sliderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider.WithInternal(err)
	}
	return true, nil
}

func (r *sliderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider.WithInternal(err)
	}
	return true, nil
}
