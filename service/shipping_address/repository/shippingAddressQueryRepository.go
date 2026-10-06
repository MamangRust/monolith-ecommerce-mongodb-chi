package repository

import (
	"context"
	"regexp"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type shippingAddressQueryRepository struct {
	db *mongo.Database
}

func NewShippingAddressQueryRepository(db *mongo.Database) ShippingAddressQueryRepository {
	return &shippingAddressQueryRepository{db: db}
}

func (r *shippingAddressQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("shipping_addresses")
}

func fmtTimePtrS(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05.000")
}

func addrToResult(a *salemongo.ShippingAddress, total int64) *ShippingAddressResult {
	return &ShippingAddressResult{
		ShippingAddressID: a.ShippingAddressID,
		OrderID:           a.OrderID,
		Alamat:            a.Alamat,
		Provinsi:          a.Provinsi,
		Kota:              a.Kota,
		Negara:            a.Negara,
		Courier:           a.Courier,
		ShippingMethod:    a.ShippingMethod,
		ShippingCost:      a.ShippingCost,
		CreatedAt:         strPtr(fmtTimePtrS(a.CreatedAt)),
		UpdatedAt:         strPtr(fmtTimePtrS(a.UpdatedAt)),
		DeletedAt:         strPtr(fmtTimePtrS(a.DeletedAt)),
		TotalCount:        total,
	}
}

func strPtr(s string) *string { return &s }

func (r *shippingAddressQueryRepository) findAddresses(ctx context.Context, search string, base bson.M, page, pageSize int) ([]*ShippingAddressResult, error) {
	filter := bson.M{}
	for k, v := range base {
		filter[k] = v
	}
	if search != "" {
		rx := regexp.QuoteMeta(search)
		filter["$or"] = []bson.M{
			{"alamat": bson.M{"$regex": rx, "$options": "i"}},
			{"provinsi": bson.M{"$regex": rx, "$options": "i"}},
		}
	}

	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, shippingaddress_errors.ErrFindAllShippingAddress
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "shipping_address_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, shippingaddress_errors.ErrFindAllShippingAddress
	}
	defer cursor.Close(ctx)

	var addrs []*salemongo.ShippingAddress
	if err := cursor.All(ctx, &addrs); err != nil {
		return nil, shippingaddress_errors.ErrFindAllShippingAddress
	}

	var results []*ShippingAddressResult
	for _, a := range addrs {
		results = append(results, addrToResult(a, total))
	}
	return results, nil
}

func (r *shippingAddressQueryRepository) FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	return r.findAddresses(ctx, req.Search, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
}

func (r *shippingAddressQueryRepository) FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	return r.findAddresses(ctx, req.Search, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
}

func (r *shippingAddressQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	return r.findAddresses(ctx, req.Search, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Page, req.PageSize)
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, shippingID int) (*ShippingAddressResult, error) {
	var a salemongo.ShippingAddress
	if err := r.coll().FindOne(ctx, bson.M{"shipping_address_id": int32(shippingID), "deleted_at": nil}).Decode(&a); err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}
	return addrToResult(&a, 0), nil
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, orderID int) (*ShippingAddressResult, error) {
	var a salemongo.ShippingAddress
	if err := r.coll().FindOne(ctx, bson.M{"order_id": int32(orderID)}).Decode(&a); err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}
	return addrToResult(&a, 0), nil
}
