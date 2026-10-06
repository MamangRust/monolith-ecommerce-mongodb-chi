package repository

import (
	"context"
	"regexp"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderQueryRepository struct {
	db *mongo.Database
}

func NewOrderQueryRepository(db *mongo.Database) OrderQueryRepository {
	return &orderQueryRepository{db: db}
}

func buildOrderSearchFilter(search string, extra bson.M) bson.M {
	filter := bson.M{"deleted_at": nil}
	if search != "" {
		escaped := regexp.QuoteMeta(search)
		filter["order_id_str"] = bson.M{"$regex": escaped, "$options": "i"}
	}
	for k, v := range extra {
		filter[k] = v
	}
	return filter
}

func (r *orderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	return r.findOrders(ctx, req.Search, bson.M{"deleted_at": nil}, &req.Page, &req.PageSize, false)
}

func (r *orderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	return r.findOrders(ctx, req.Search, bson.M{"deleted_at": nil}, &req.Page, &req.PageSize, true)
}

func (r *orderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	return r.findOrders(ctx, req.Search, bson.M{"deleted_at": bson.M{"$ne": nil}}, &req.Page, &req.PageSize, true)
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*OrderResult, error) {
	return r.findOrders(ctx, req.Search, bson.M{"merchant_id": int32(req.MerchantID), "deleted_at": nil}, &req.Page, &req.PageSize, false)
}

func (r *orderQueryRepository) findOrders(ctx context.Context, search string, base bson.M, page, pageSize *int, withDeleted bool) ([]*OrderResult, error) {
	filter := bson.M{}
	for k, v := range base {
		filter[k] = v
	}
	if search != "" {
		// Match order_id numerically when the search is a number; otherwise
		// cast to string is emulated by matching against nothing.
		if id, err := parseOrderID(search); err == nil {
			filter["order_id"] = id
		} else {
			escaped := regexp.QuoteMeta(search)
			filter["order_id_str"] = bson.M{"$regex": escaped, "$options": "i"}
		}
	}

	p := 1
	if page != nil && *page > 0 {
		p = *page
	}
	ps := 10
	if pageSize != nil && *pageSize > 0 {
		ps = *pageSize
	}
	offset := (p - 1) * ps

	coll := r.db.Collection("orders")

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}
	totalCount := int(total)

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(ps))
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var orders []*salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}

	results := make([]*OrderResult, 0, len(orders))
	for _, o := range orders {
		results = append(results, &OrderResult{
			OrderID:    o.OrderID,
			UserID:     o.UserID,
			MerchantID: o.MerchantID,
			TotalPrice: o.TotalPrice,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
			DeletedAt:  o.DeletedAt,
			TotalCount: int64(totalCount),
		})
	}
	return results, nil
}

func (r *orderQueryRepository) FindByID(ctx context.Context, id int) (*salemongo.Order, error) {
	var order salemongo.Order
	err := r.db.Collection("orders").FindOne(ctx, bson.M{"order_id": id, "deleted_at": nil}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order_errors.ErrOrderNotFound.WithInternal(err)
		}
		return nil, order_errors.ErrFindById.WithInternal(err)
	}
	return &order, nil
}

func parseOrderID(s string) (int32, error) {
	var id int
	if err := parseIntStrict(s, &id); err != nil {
		return 0, err
	}
	return int32(id), nil
}

func parseIntStrict(s string, out *int) error {
	if s == "" {
		return errInvalidInt
	}
	n := 0
	neg := false
	for i, c := range s {
		if i == 0 && (c == '-' || c == '+') {
			neg = c == '-'
			continue
		}
		if c < '0' || c > '9' {
			return errInvalidInt
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		n = -n
	}
	*out = n
	return nil
}

var errInvalidInt = errorIntInvalid{}

type errorIntInvalid struct{}

func (errorIntInvalid) Error() string { return "invalid integer" }
