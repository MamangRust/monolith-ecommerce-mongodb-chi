package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

type categoryCommandRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewCategoryCommandRepository(db *mongo.Database, coll *mongo.Collection) *categoryCommandRepository {
	return &categoryCommandRepository{db: db, coll: coll}
}

func (r *categoryCommandRepository) Create(ctx context.Context, request *requests.CreateCategoryRequest) (*salemongo.Category, error) {
	categoryID, err := salemongo.NextID(ctx, r.db, "categories")
	if err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}

	now := time.Now()
	category := &salemongo.Category{
		CategoryID:    categoryID,
		Name:          request.Name,
		Description:   &request.Description,
		SlugCategory:  request.SlugCategory,
		ImageCategory: &request.ImageCategory,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	if _, err := r.coll.InsertOne(ctx, category); err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}
	return category, nil
}

func (r *categoryCommandRepository) Update(ctx context.Context, request *requests.UpdateCategoryRequest) (*salemongo.Category, error) {
	now := time.Now()
	filter := bson.M{"category_id": *request.CategoryID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"name":           request.Name,
		"description":    &request.Description,
		"slug_category":  request.SlugCategory,
		"image_category": &request.ImageCategory,
		"updated_at":     now,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var category salemongo.Category
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) Trash(ctx context.Context, categoryID int) (*salemongo.Category, error) {
	now := time.Now()
	filter := bson.M{"category_id": categoryID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, category_errors.ErrCategoryNotFound
	}

	var category salemongo.Category
	if err := r.coll.FindOne(ctx, bson.M{"category_id": categoryID}).Decode(&category); err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) Restore(ctx context.Context, categoryID int) (*salemongo.Category, error) {
	filter := bson.M{"category_id": categoryID, "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil, "updated_at": time.Now()}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, category_errors.ErrCategoryNotFound
	}

	var category salemongo.Category
	if err := r.coll.FindOne(ctx, bson.M{"category_id": categoryID}).Decode(&category); err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) DeletePermanent(ctx context.Context, categoryID int) (bool, error) {
	filter := bson.M{"category_id": categoryID, "deleted_at": bson.M{"$ne": nil}}
	result, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, category_errors.ErrCategoryNotFound
	}
	return true, nil
}

func (r *categoryCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, category_errors.ErrRestoreAllCategories.WithInternal(err)
	}
	return true, nil
}

func (r *categoryCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories.WithInternal(err)
	}
	return true, nil
}
