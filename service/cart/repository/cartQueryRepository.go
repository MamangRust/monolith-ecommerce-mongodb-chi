package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cartQueryRepository struct {
	coll *mongo.Collection
}

func NewCartQueryRepository(coll *mongo.Collection) CartQueryRepository {
	return &cartQueryRepository{coll: coll}
}

func (r *cartQueryRepository) FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*CartResult, error) {
	filter := bson.M{"user_id": req.UserID, "deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var carts []salemongo.Cart
	if err := cursor.All(ctx, &carts); err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}

	var out []*CartResult
	for _, c := range carts {
		cart := c
		out = append(out, &CartResult{
			CartID: cart.CartID, UserID: cart.UserID, ProductID: cart.ProductID,
			Name: cart.Name, Price: cart.Price, Image: cart.Image,
			Quantity: cart.Quantity, Weight: cart.Weight,
			CreatedAt:  formatCartTime(cart.CreatedAt),
			UpdatedAt:  formatCartTime(cart.UpdatedAt),
			TotalCount: totalCount,
		})
	}
	return out, nil
}

func formatCartTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
