package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	product_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productCommandRepository struct {
	db *mongo.Database
}

func NewProductCommandRepository(db *mongo.Database) ProductCommandRepository {
	return &productCommandRepository{db: db}
}

func (r *productCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("products")
}

func (r *productCommandRepository) Create(ctx context.Context, req *requests.CreateProductRequest) (*salemongo.Product, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "products")
	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	description := req.Description
	brand := req.Brand
	weight := int32(req.Weight)
	slug := ""
	if req.SlugProduct != nil {
		slug = *req.SlugProduct
	}
	image := req.ImageProduct

	doc := &salemongo.Product{
		ProductID:    id,
		MerchantID:   int32(req.MerchantID),
		CategoryID:   int32(req.CategoryID),
		Name:         req.Name,
		Description:  &description,
		Price:        int32(req.Price),
		CountInStock: int32(req.CountInStock),
		Brand:        &brand,
		Weight:       &weight,
		SlugProduct:  &slug,
		ImageProduct: &image,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	if req.Rating != nil {
		rating := float64(*req.Rating)
		doc.Rating = &rating
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, product_errors.ErrCreateProduct
	}
	return doc, nil
}

func (r *productCommandRepository) Update(ctx context.Context, req *requests.UpdateProductRequest) (*salemongo.Product, error) {
	now := time.Now()
	productID := int32(*req.ProductID)
	filter := bson.M{"product_id": productID, "deleted_at": nil}

	var doc salemongo.Product
	if err := r.coll().FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, product_errors.ErrProductNotFound
	}

	description := req.Description
	brand := req.Brand
	weight := int32(req.Weight)
	slug := ""
	if req.SlugProduct != nil {
		slug = *req.SlugProduct
	}
	image := req.ImageProduct

	update := bson.M{"$set": bson.M{
		"merchant_id":    int32(req.MerchantID),
		"category_id":    int32(req.CategoryID),
		"name":           req.Name,
		"description":    description,
		"price":          int32(req.Price),
		"count_in_stock": int32(req.CountInStock),
		"brand":          brand,
		"weight":         weight,
		"slug_product":   slug,
		"image_product":  image,
		"updated_at":     now,
	}}
	if req.Rating != nil {
		update["$set"].(bson.M)["rating"] = float64(*req.Rating)
	}

	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, product_errors.ErrUpdateProduct
	}
	return &doc, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*salemongo.Product, error) {
	now := time.Now()
	filter := bson.M{"product_id": int32(productID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"count_in_stock": int32(stock),
		"updated_at":     now,
	}}

	var doc salemongo.Product
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}
	return &doc, nil
}

func (r *productCommandRepository) AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*salemongo.Product, error) {
	now := time.Now()
	filter := bson.M{"product_id": int32(productID), "deleted_at": nil}

	// Record the adjustment operation
	adjColl := r.db.Collection("product_stock_adjustments")
	adjDoc := &salemongo.ProductStockAdjustment{
		OperationID: operationID,
		ProductID:   int32(productID),
		Delta:       int32(delta),
		CreatedAt:   &now,
	}
	if _, err := adjColl.InsertOne(ctx, adjDoc); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	update := bson.M{"$inc": bson.M{"count_in_stock": int32(delta)}, "$set": bson.M{"updated_at": now}}

	var doc salemongo.Product
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}
	return &doc, nil
}

func (r *productCommandRepository) Trash(ctx context.Context, productID int) (*salemongo.Product, error) {
	now := time.Now()
	filter := bson.M{"product_id": int32(productID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var doc salemongo.Product
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, product_errors.ErrTrashedProduct
	}
	return &doc, nil
}

func (r *productCommandRepository) Restore(ctx context.Context, productID int) (*salemongo.Product, error) {
	filter := bson.M{"product_id": int32(productID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var doc salemongo.Product
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, product_errors.ErrRestoreProduct
	}
	return &doc, nil
}

func (r *productCommandRepository) DeletePermanent(ctx context.Context, productID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"product_id": int32(productID)}); err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}
	return true, nil
}

func (r *productCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}
	return true, nil
}

func (r *productCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, product_errors.ErrDeleteAllProducts
	}
	return true, nil
}
