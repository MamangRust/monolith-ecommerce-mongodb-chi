package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderCommandRepository struct {
	db *mongo.Database
}

func NewOrderCommandRepository(db *mongo.Database) OrderCommandRepository {
	return &orderCommandRepository{db: db}
}

func (r *orderCommandRepository) Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*salemongo.Order, error) {
	now := time.Now()

	orderID, err := salemongo.NextID(ctx, r.db, "orders")
	if err != nil {
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}

	order := &salemongo.Order{
		OrderID:    orderID,
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	if _, err := r.db.Collection("orders").InsertOne(ctx, order); err != nil {
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}
	return order, nil
}

func (r *orderCommandRepository) Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*salemongo.Order, error) {
	now := time.Now()

	var order salemongo.Order
	err := r.db.Collection("orders").FindOne(ctx, bson.M{"order_id": request.OrderID, "deleted_at": nil}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order_errors.ErrOrderNotFound
		}
		return nil, order_errors.ErrOrderNotFound.WithInternal(err)
	}

	_, err = r.db.Collection("orders").UpdateOne(ctx,
		bson.M{"order_id": request.OrderID},
		bson.M{"$set": bson.M{"total_price": int32(request.TotalPrice), "updated_at": now}})
	if err != nil {
		return nil, order_errors.ErrUpdateOrder.WithInternal(err)
	}

	order.TotalPrice = int32(request.TotalPrice)
	order.UpdatedAt = &now
	return &order, nil
}

func (r *orderCommandRepository) Trash(ctx context.Context, order_id int) (*salemongo.Order, error) {
	now := time.Now()

	result, err := r.db.Collection("orders").UpdateOne(ctx,
		bson.M{"order_id": order_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}})
	if err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, order_errors.ErrOrderNotFound
	}

	var order salemongo.Order
	if err := r.db.Collection("orders").FindOne(ctx, bson.M{"order_id": order_id}).Decode(&order); err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) Restore(ctx context.Context, order_id int) (*salemongo.Order, error) {
	result, err := r.db.Collection("orders").UpdateOne(ctx,
		bson.M{"order_id": order_id, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}
	if result.MatchedCount == 0 {
		return nil, order_errors.ErrOrderNotFound
	}

	var order salemongo.Order
	if err := r.db.Collection("orders").FindOne(ctx, bson.M{"order_id": order_id}).Decode(&order); err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) FindTrashedByID(ctx context.Context, order_id int) (*salemongo.Order, error) {
	var order salemongo.Order
	err := r.db.Collection("orders").FindOne(ctx, bson.M{"order_id": order_id, "deleted_at": bson.M{"$ne": nil}}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order_errors.ErrOrderNotFound
		}
		return nil, order_errors.ErrOrderNotFound.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) FindTrashed(ctx context.Context) ([]*salemongo.Order, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.db.Collection("orders").Find(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, opts)
	if err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var orders []*salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	return orders, nil
}

func (r *orderCommandRepository) DeletePermanent(ctx context.Context, order_id int) (bool, error) {
	// Delete child reservations first so orphaned reservations never linger.
	if _, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": order_id}); err != nil {
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}

	result, err := r.db.Collection("orders").DeleteOne(ctx, bson.M{"order_id": order_id, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, order_errors.ErrOrderNotFound
	}
	return true, nil
}

func (r *orderCommandRepository) DeletePermanentWithChildren(ctx context.Context, order_id int) (bool, error) {
	return r.DeletePermanent(ctx, order_id)
}

func (r *orderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("orders").UpdateMany(ctx,
		bson.M{"deleted_at": bson.M{"$ne": nil}},
		bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, order_errors.ErrRestoreAllOrder.WithInternal(err)
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	trashed, err := r.db.Collection("orders").Distinct(ctx, "order_id", bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	if len(trashed) > 0 {
		for _, id := range trashed {
			if orderID, ok := id.(int32); ok {
				if _, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": orderID}); err != nil {
					return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
				}
			}
		}
	}

	if _, err := r.db.Collection("orders").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}); err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}
