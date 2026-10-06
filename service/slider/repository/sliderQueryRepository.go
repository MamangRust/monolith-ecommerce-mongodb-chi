package repository

import (
	"context"
	"errors"
	"regexp"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sliderQueryRepository struct {
	coll *mongo.Collection
}

func NewSliderQueryRepository(coll *mongo.Collection) *sliderQueryRepository {
	return &sliderQueryRepository{coll: coll}
}

func sliderToResult(s salemongo.Slider) *SliderResult {
	createdAt, updatedAt, deletedAt := "", "", ""
	if s.CreatedAt != nil {
		createdAt = s.CreatedAt.Format("2006-01-02 15:04:05.000")
	}
	if s.UpdatedAt != nil {
		updatedAt = s.UpdatedAt.Format("2006-01-02 15:04:05.000")
	}
	if s.DeletedAt != nil {
		deletedAt = s.DeletedAt.Format("2006-01-02 15:04:05.000")
	}
	return &SliderResult{
		SliderID:  s.SliderID,
		Name:      s.Name,
		Image:     s.Image,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
		DeletedAt: &deletedAt,
	}
}

func (r *sliderQueryRepository) findPage(ctx context.Context, filter bson.M, search string, page, pageSize int) ([]*SliderResult, error) {
	if search != "" {
		filter["name"] = bson.M{"$regex": regexp.QuoteMeta(search), "$options": "i"}
	}
	total, _ := r.coll.CountDocuments(ctx, filter)
	opts := options.Find().SetSort(bson.D{{Key: "slider_id", Value: -1}}).SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, slider_errors.ErrFindAllSliders.WithInternal(err)
	}
	defer cursor.Close(ctx)
	var results []*SliderResult
	for cursor.Next(ctx) {
		var s salemongo.Slider
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		r := sliderToResult(s)
		r.TotalCount = total
		results = append(results, r)
	}
	return results, nil
}

func (r *sliderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	return r.findPage(ctx, bson.M{"deleted_at": nil}, req.Search, req.Page, req.PageSize)
}
func (r *sliderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	return r.FindAll(ctx, req)
}
func (r *sliderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	return r.findPage(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Search, req.Page, req.PageSize)
}
func (r *sliderQueryRepository) FindByID(ctx context.Context, id int) (*SliderResult, error) {
	var s salemongo.Slider
	err := r.coll.FindOne(ctx, bson.M{"slider_id": id, "deleted_at": nil}).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, slider_errors.ErrSliderNotFound
		}
		return nil, slider_errors.ErrFindAllSliders.WithInternal(err)
	}
	return sliderToResult(s), nil
}

var _ = time.Now
