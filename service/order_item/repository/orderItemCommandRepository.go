package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderItemCommandRepository struct {
	db *mongo.Database
}

func NewOrderItemCommandRepository(db *mongo.Database) OrderItemCommandRepository {
	return &orderItemCommandRepository{db: db}
}

func (r *orderItemCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("order_items")
}

func orderItemToRow(o *salemongo.OrderItem) *salemongo.OrderItemRow {
	return &salemongo.OrderItemRow{
		OrderItemID: o.OrderItemID,
		OrderID:     o.OrderID,
		ProductID:   o.ProductID,
		Quantity:    o.Quantity,
		Price:       o.Price,
		CreatedAt:   derefOrderItemTime(o.CreatedAt),
		UpdatedAt:   derefOrderItemTime(o.UpdatedAt),
	}
}

func derefOrderItemTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func (r *orderItemCommandRepository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItemRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "order_items")
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	doc := &salemongo.OrderItem{
		OrderItemID: id,
		OrderID:     int32(req.OrderID),
		ProductID:   int32(req.ProductID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}
	return orderItemToRow(doc), nil
}

func (r *orderItemCommandRepository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItemRow, error) {
	now := time.Now()
	filter := bson.M{"order_item_id": int32(req.OrderItemID), "deleted_at": nil}

	var doc salemongo.OrderItem
	if err := r.coll().FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}

	update := bson.M{"$set": bson.M{
		"quantity":   int32(req.Quantity),
		"price":      int32(req.Price),
		"updated_at": now,
	}}

	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}
	return orderItemToRow(&doc), nil
}

func (r *orderItemCommandRepository) Trash(ctx context.Context, orderItemID int) (*salemongo.OrderItem, error) {
	now := time.Now()
	filter := bson.M{"order_item_id": int32(orderItemID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var doc salemongo.OrderItem
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}
	return &doc, nil
}

func (r *orderItemCommandRepository) Restore(ctx context.Context, orderItemID int) (*salemongo.OrderItem, error) {
	filter := bson.M{"order_item_id": int32(orderItemID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var doc salemongo.OrderItem
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&doc); err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}
	return &doc, nil
}

func (r *orderItemCommandRepository) DeletePermanent(ctx context.Context, orderItemID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"order_item_id": int32(orderItemID)}); err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"order_id": int32(orderID)}); err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) CalculateTotalPrice(ctx context.Context, orderID int) (int, error) {
	cursor, err := r.coll().Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"order_id": int32(orderID), "deleted_at": nil}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": bson.M{"$multiply": bson.A{"$quantity", "$price"}}}}}},
	})
	if err != nil {
		return 0, orderitem_errors.ErrCalculateTotalPrice
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int32 `bson:"total"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, orderitem_errors.ErrCalculateTotalPrice
		}
	}
	return int(result.Total), nil
}
