package repository

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	product_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productQueryRepository struct {
	db *mongo.Database
}

func NewProductQueryRepository(db *mongo.Database) ProductQueryRepository {
	return &productQueryRepository{db: db}
}

func (r *productQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("products")
}

func fmtProductTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05.000")
	}
	return ""
}

func productToResult(p *salemongo.Product, totalCount int64) *ProductResult {
	return &ProductResult{
		ProductID:    p.ProductID,
		MerchantID:   p.MerchantID,
		CategoryID:   p.CategoryID,
		Name:         p.Name,
		Description:  p.Description,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		Brand:        p.Brand,
		Weight:       p.Weight,
		Rating:       p.Rating,
		SlugProduct:  p.SlugProduct,
		ImageProduct: p.ImageProduct,
		CreatedAt:    stringPtr(fmtProductTime(p.CreatedAt)),
		UpdatedAt:    stringPtr(fmtProductTime(p.UpdatedAt)),
		DeletedAt:    stringPtr(fmtProductTime(p.DeletedAt)),
		TotalCount:   totalCount,
	}
}

func (r *productQueryRepository) findProducts(ctx context.Context, filter bson.M, page, pageSize int) ([]*ProductResult, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "product_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.Product
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	var results []*ProductResult
	for _, p := range docs {
		results = append(results, productToResult(p, total))
	}
	return results, nil
}

func productFilter(req *requests.FindAllProduct, trashed bool) bson.M {
	filter := bson.M{}
	if trashed {
		filter["deleted_at"] = bson.M{"$ne": nil}
	} else {
		filter["deleted_at"] = nil
	}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	return filter
}

func (r *productQueryRepository) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	results, err := r.findProducts(ctx, productFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, product_errors.ErrFindAllProducts.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	results, err := r.findProducts(ctx, productFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, product_errors.ErrFindActiveProducts.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	results, err := r.findProducts(ctx, productFilter(req, true), req.Page, req.PageSize)
	if err != nil {
		return nil, product_errors.ErrFindTrashedProducts.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*ProductResult, error) {
	filter := bson.M{"merchant_id": int32(req.MerchantID), "deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	if req.MinPrice != nil {
		filter["price"] = bson.M{"$gte": int32(*req.MinPrice)}
	}
	if req.MaxPrice != nil {
		if _, ok := filter["price"]; ok {
			filter["price"].(bson.M)["$lte"] = int32(*req.MaxPrice)
		} else {
			filter["price"] = bson.M{"$lte": int32(*req.MaxPrice)}
		}
	}
	results, err := r.findProducts(ctx, filter, 1, 100)
	if err != nil {
		return nil, product_errors.ErrFindProductsByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*ProductResult, error) {
	// Look up category by name first
	catColl := r.db.Collection("categories")
	var cat salemongo.Category
	if err := catColl.FindOne(ctx, bson.M{"name": bson.M{"$regex": req.CategoryName, "$options": "i"}, "deleted_at": nil}).Decode(&cat); err != nil {
		return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
	}

	filter := bson.M{"category_id": cat.CategoryID, "deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	results, err := r.findProducts(ctx, filter, 1, 100)
	if err != nil {
		return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindByID(ctx context.Context, productID int) (*ProductResult, error) {
	var p salemongo.Product
	if err := r.coll().FindOne(ctx, bson.M{"product_id": int32(productID), "deleted_at": nil}).Decode(&p); err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	return productToResult(&p, 1), nil
}

func (r *productQueryRepository) FindByIDTrashed(ctx context.Context, productID int) (*ProductResult, error) {
	var p salemongo.Product
	if err := r.coll().FindOne(ctx, bson.M{"product_id": int32(productID), "deleted_at": bson.M{"$ne": nil}}).Decode(&p); err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	return productToResult(&p, 1), nil
}

var _ = fmt.Sprintf
var _ = time.Now
