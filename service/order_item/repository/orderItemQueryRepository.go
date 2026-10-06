package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderItemQueryRepository struct {
	db *mongo.Database
}

func NewOrderItemQueryRepository(db *mongo.Database) OrderItemQueryRepository {
	return &orderItemQueryRepository{db: db}
}

func (r *orderItemQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("order_items")
}

func orderItemToListRow(o *salemongo.OrderItem, total int64) *salemongo.OrderItemListRow {
	return &salemongo.OrderItemListRow{
		OrderItemRow: *orderItemToRow(o),
		TotalCount:   total,
	}
}

func (r *orderItemQueryRepository) findItems(ctx context.Context, filter bson.M, page, pageSize int) ([]*salemongo.OrderItemListRow, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "order_item_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.OrderItem
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems
	}

	var results []*salemongo.OrderItemListRow
	for _, o := range docs {
		results = append(results, orderItemToListRow(o, total))
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemListRow, error) {
	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["product_id"] = req.Search
	}
	results, err := r.findItems(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemActiveRow, error) {
	listRows, err := r.findItems(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, orderitem_errors.ErrFindByActive
	}
	var active []*salemongo.OrderItemActiveRow
	for _, row := range listRows {
		active = append(active, &salemongo.OrderItemActiveRow{
			OrderItemRow: row.OrderItemRow,
			TotalCount:   row.TotalCount,
		})
	}
	return active, nil
}

func (r *orderItemQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*salemongo.OrderItemTrashedRow, error) {
	listRows, err := r.findItems(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Page, req.PageSize)
	if err != nil {
		return nil, orderitem_errors.ErrFindByTrashed
	}
	var trashed []*salemongo.OrderItemTrashedRow
	for _, row := range listRows {
		trashed = append(trashed, &salemongo.OrderItemTrashedRow{
			OrderItemRow: row.OrderItemRow,
			TotalCount:   row.TotalCount,
		})
	}
	return trashed, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderItemRow, error) {
	cursor, err := r.coll().Find(ctx, bson.M{"order_id": int32(orderID), "deleted_at": nil})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.OrderItem
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	var results []*salemongo.OrderItemRow
	for _, o := range docs {
		results = append(results, orderItemToRow(o))
	}
	return results, nil
}
