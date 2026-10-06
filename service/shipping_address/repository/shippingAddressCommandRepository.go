package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type shippingAddressCommandRepository struct {
	db *mongo.Database
}

func NewShippingAddressCommandRepository(db *mongo.Database) ShippingAddressCommandRepository {
	return &shippingAddressCommandRepository{db: db}
}

func (r *shippingAddressCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("shipping_addresses")
}

func (r *shippingAddressCommandRepository) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*salemongo.ShippingAddress, error) {
	now := time.Now()

	id, err := salemongo.NextID(ctx, r.db, "shipping_addresses")
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}

	var orderID int32
	if request.OrderID != nil {
		orderID = int32(*request.OrderID)
	}

	addr := &salemongo.ShippingAddress{
		ShippingAddressID: id,
		OrderID:           orderID,
		Alamat:            request.Alamat,
		Provinsi:          request.Provinsi,
		Kota:              request.Kota,
		Negara:            request.Negara,
		Courier:           request.Courier,
		ShippingMethod:    request.ShippingMethod,
		ShippingCost:      float64(request.ShippingCost),
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}

	if _, err := r.coll().InsertOne(ctx, addr); err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}
	return addr, nil
}

func (r *shippingAddressCommandRepository) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*salemongo.ShippingAddress, error) {
	now := time.Now()
	filter := bson.M{"shipping_address_id": int32(*request.ShippingID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"alamat":          request.Alamat,
		"provinsi":        request.Provinsi,
		"kota":            request.Kota,
		"negara":          request.Negara,
		"courier":         request.Courier,
		"shipping_method": request.ShippingMethod,
		"shipping_cost":   float64(request.ShippingCost),
		"updated_at":      now,
	}}

	var addr salemongo.ShippingAddress
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}
	return &addr, nil
}

func (r *shippingAddressCommandRepository) Trash(ctx context.Context, shippingID int) (*salemongo.ShippingAddress, error) {
	now := time.Now()
	filter := bson.M{"shipping_address_id": int32(shippingID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	var addr salemongo.ShippingAddress
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrTrashShippingAddress
	}
	return &addr, nil
}

func (r *shippingAddressCommandRepository) Restore(ctx context.Context, shippingID int) (*salemongo.ShippingAddress, error) {
	filter := bson.M{"shipping_address_id": int32(shippingID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}

	var addr salemongo.ShippingAddress
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}
	return &addr, nil
}

func (r *shippingAddressCommandRepository) DeletePermanent(ctx context.Context, shippingID int) (bool, error) {
	res, err := r.coll().DeleteOne(ctx, bson.M{"shipping_address_id": int32(shippingID)})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}
	if res.DeletedCount == 0 {
		return false, shippingaddress_errors.ErrShippingAddressNotFound
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"order_id": int32(orderID)}); err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	return err == nil, err
}

func (r *shippingAddressCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	return err == nil, err
}
