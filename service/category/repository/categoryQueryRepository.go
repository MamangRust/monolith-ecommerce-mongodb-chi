package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type categoryQueryRepository struct {
	coll *mongo.Collection
}

func NewCategoryQueryRepository(coll *mongo.Collection) *categoryQueryRepository {
	return &categoryQueryRepository{coll: coll}
}

func (r *categoryQueryRepository) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := int64((req.Page - 1) * req.PageSize)
	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(offset).
		SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var cats []salemongo.Category
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}

	var results []*CategoryResult
	for _, c := range cats {
		cr := &CategoryResult{
			CategoryID:   c.CategoryID,
			Name:         c.Name,
			Description:  c.Description,
			SlugCategory: c.SlugCategory,
			TotalCount:   totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := int64((req.Page - 1) * req.PageSize)
	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(offset).
		SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var cats []salemongo.Category
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}

	var results []*CategoryResult
	for _, c := range cats {
		cr := &CategoryResult{
			CategoryID:   c.CategoryID,
			Name:         c.Name,
			Description:  c.Description,
			SlugCategory: c.SlugCategory,
			TotalCount:   totalCount,
			DeletedAt:    c.DeletedAt,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := int64((req.Page - 1) * req.PageSize)
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(offset).
		SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var cats []salemongo.Category
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}

	var results []*CategoryResult
	for _, c := range cats {
		now := time.Now()
		cr := &CategoryResult{
			CategoryID:   c.CategoryID,
			Name:         c.Name,
			Description:  c.Description,
			SlugCategory: c.SlugCategory,
			TotalCount:   totalCount,
			DeletedAt:    &now,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindByID(ctx context.Context, categoryID int) (*salemongo.Category, error) {
	var category salemongo.Category
	err := r.coll.FindOne(ctx, bson.M{"category_id": categoryID, "deleted_at": nil}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrFindCategoryById.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryQueryRepository) FindByIDTrashed(ctx context.Context, categoryID int) (*salemongo.Category, error) {
	var category salemongo.Category
	err := r.coll.FindOne(ctx, bson.M{"category_id": categoryID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrFindCategoryByIdTrashed.WithInternal(err)
	}
	return &category, nil
}
