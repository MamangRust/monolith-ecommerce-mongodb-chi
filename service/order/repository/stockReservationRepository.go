package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type stockReservationRepository struct {
	db *mongo.Database
}

func NewStockReservationRepository(db *mongo.Database) StockReservationRepository {
	return &stockReservationRepository{db: db}
}

func (r *stockReservationRepository) GetByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderStockReservation, error) {
	cursor, err := r.db.Collection("order_stock_reservations").Find(ctx, bson.M{"order_id": int32(orderID)})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reservations []*salemongo.OrderStockReservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *stockReservationRepository) Upsert(ctx context.Context, orderID, productID, quantity int) (*salemongo.OrderStockReservation, error) {
	now := time.Now()

	var existing salemongo.OrderStockReservation
	err := r.db.Collection("order_stock_reservations").FindOne(ctx, bson.M{
		"order_id":   int32(orderID),
		"product_id": int32(productID),
	}).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		reservationID, err := salemongo.NextID(ctx, r.db, "order_stock_reservations")
		if err != nil {
			return nil, err
		}
		reservation := &salemongo.OrderStockReservation{
			ReservationID: reservationID,
			OrderID:       int32(orderID),
			ProductID:     int32(productID),
			Quantity:      int32(quantity),
			Status:        "reserved",
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}
		if _, err := r.db.Collection("order_stock_reservations").InsertOne(ctx, reservation); err != nil {
			return nil, err
		}
		return reservation, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = r.db.Collection("order_stock_reservations").UpdateOne(ctx,
		bson.M{"_id": existing.ID},
		bson.M{"$set": bson.M{"quantity": int32(quantity), "status": "reserved", "updated_at": now}})
	if err != nil {
		return nil, err
	}
	existing.Quantity = int32(quantity)
	existing.Status = "reserved"
	existing.UpdatedAt = &now
	return &existing, nil
}

func (r *stockReservationRepository) UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*salemongo.OrderStockReservation, error) {
	now := time.Now()
	var reservation salemongo.OrderStockReservation
	err := r.db.Collection("order_stock_reservations").FindOneAndUpdate(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(productID)},
		bson.M{"$set": bson.M{"quantity": int32(quantity), "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&reservation)
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *stockReservationRepository) Release(ctx context.Context, orderID, productID int) (*salemongo.OrderStockReservation, error) {
	now := time.Now()
	var reservation salemongo.OrderStockReservation
	err := r.db.Collection("order_stock_reservations").FindOneAndUpdate(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(productID), "status": "reserved"},
		bson.M{"$set": bson.M{"status": "released", "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&reservation)
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *stockReservationRepository) Reserve(ctx context.Context, orderID, productID int) (*salemongo.OrderStockReservation, error) {
	now := time.Now()
	var reservation salemongo.OrderStockReservation
	err := r.db.Collection("order_stock_reservations").FindOneAndUpdate(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(productID), "status": "released"},
		bson.M{"$set": bson.M{"status": "reserved", "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&reservation)
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// trashedOrderIDs returns order ids whose deleted_at is set.
func (r *stockReservationRepository) trashedOrderIDs(ctx context.Context) ([]int32, error) {
	var orders []*salemongo.Order
	cursor, err := r.db.Collection("orders").Find(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, options.Find().SetProjection(bson.M{"order_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	ids := make([]int32, 0, len(orders))
	for _, o := range orders {
		ids = append(ids, o.OrderID)
	}
	return ids, nil
}

// activeOrderIDs returns order ids whose deleted_at is nil.
func (r *stockReservationRepository) activeOrderIDs(ctx context.Context) ([]int32, error) {
	var orders []*salemongo.Order
	cursor, err := r.db.Collection("orders").Find(ctx, bson.M{"deleted_at": nil}, options.Find().SetProjection(bson.M{"order_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	ids := make([]int32, 0, len(orders))
	for _, o := range orders {
		ids = append(ids, o.OrderID)
	}
	return ids, nil
}

func inFilter(ids []int32) bson.M {
	return bson.M{"order_id": bson.M{"$in": ids}}
}

func (r *stockReservationRepository) GetReservedForTrashedOrders(ctx context.Context) ([]*salemongo.OrderStockReservation, error) {
	ids, err := r.trashedOrderIDs(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	cursor, err := r.db.Collection("order_stock_reservations").Find(ctx, bson.M{"status": "reserved", "order_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reservations []*salemongo.OrderStockReservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *stockReservationRepository) GetReleasedForTrashedOrders(ctx context.Context) ([]*salemongo.OrderStockReservation, error) {
	ids, err := r.trashedOrderIDs(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	cursor, err := r.db.Collection("order_stock_reservations").Find(ctx, bson.M{"status": "released", "order_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reservations []*salemongo.OrderStockReservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *stockReservationRepository) DeleteByOrder(ctx context.Context, orderID int) error {
	_, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": int32(orderID)})
	return err
}

func (r *stockReservationRepository) DeleteByOrderProduct(ctx context.Context, orderID, productID int) error {
	_, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{
		"order_id":   int32(orderID),
		"product_id": int32(productID),
	})
	return err
}

func (r *stockReservationRepository) DeleteAllForTrashedOrders(ctx context.Context) error {
	ids, err := r.trashedOrderIDs(ctx)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	_, err = r.db.Collection("order_stock_reservations").DeleteMany(ctx, inFilter(ids))
	return err
}

func (r *stockReservationRepository) GetReleasedForActiveOrders(ctx context.Context) ([]*salemongo.OrderStockReservation, error) {
	ids, err := r.activeOrderIDs(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	cursor, err := r.db.Collection("order_stock_reservations").Find(ctx, bson.M{"status": "released", "order_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reservations []*salemongo.OrderStockReservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *stockReservationRepository) DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error) {
	ids, err := r.trashedOrderIDs(ctx)
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{
		"status":     "released",
		"updated_at": bson.M{"$lt": cutoff},
		"order_id":   bson.M{"$in": ids},
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *stockReservationRepository) DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.Collection("product_stock_adjustments").DeleteMany(ctx, bson.M{"created_at": bson.M{"$lt": cutoff}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
