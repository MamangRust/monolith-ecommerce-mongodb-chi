package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cartCommandRepository struct {
	coll *mongo.Collection
}

func NewCartCommandRepository(coll *mongo.Collection) CartCommandRepository {
	return &cartCommandRepository{coll: coll}
}

func (r *cartCommandRepository) CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*CartCreateResult, error) {
	now := time.Now()

	var maxResult struct {
		MaxID int32 `bson:"cart_id"`
	}
	_ = r.coll.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{Key: "cart_id", Value: -1}})).Decode(&maxResult)
	cartID := maxResult.MaxID + 1

	doc := bson.M{
		"cart_id":    cartID,
		"user_id":    int32(req.UserID),
		"product_id": int32(req.ProductID),
		"name":       req.Name,
		"price":      int32(req.Price),
		"image":      req.ImageProduct,
		"quantity":   int32(req.Quantity),
		"weight":     int32(req.Weight),
		"created_at": now,
		"updated_at": now,
	}

	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, cart_errors.ErrCreateCart.WithInternal(err)
	}

	return &CartCreateResult{
		CartID: cartID, UserID: int32(req.UserID), ProductID: int32(req.ProductID),
		Name: req.Name, Price: int32(req.Price), Image: req.ImageProduct,
		Quantity: int32(req.Quantity), Weight: int32(req.Weight),
		CreatedAt: now.Format("2006-01-02 15:04:05"), UpdatedAt: now.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *cartCommandRepository) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"cart_id": req.CartID, "user_id": req.UserID})
	if err != nil {
		return false, cart_errors.ErrDeleteCartPermanent.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, cart_errors.ErrDeleteCartPermanent
	}
	return true, nil
}

func (r *cartCommandRepository) DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error) {
	result, err := r.coll.DeleteMany(ctx, bson.M{"cart_id": bson.M{"$in": req.CartIds}, "user_id": req.UserID})
	if err != nil {
		return false, cart_errors.ErrDeleteAllCarts.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, cart_errors.ErrDeleteAllCarts
	}
	return true, nil
}
