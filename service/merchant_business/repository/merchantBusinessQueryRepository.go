package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantBusinessQueryRepository struct {
	db *mongo.Database
}

func NewMerchantBusinessQueryRepository(db *mongo.Database) MerchantBusinessQueryRepository {
	return &merchantBusinessQueryRepository{db: db}
}

func (r *merchantBusinessQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_business_information")
}

func bizToListRow(b *salemongo.MerchantBusinessInformation, total int64) *salemongo.MerchantBusinessInfoListRow {
	return &salemongo.MerchantBusinessInfoListRow{
		MerchantBusinessInfoRow: *bizToRow(b),
		TotalCount:              total,
	}
}

func (r *merchantBusinessQueryRepository) findBiz(ctx context.Context, filter bson.M, page, pageSize int) ([]*salemongo.MerchantBusinessInfoListRow, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "merchant_business_info_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.MerchantBusinessInformation
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}

	var results []*salemongo.MerchantBusinessInfoListRow
	for _, b := range docs {
		results = append(results, bizToListRow(b, total))
	}
	return results, nil
}

func (r *merchantBusinessQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoListRow, error) {
	results, err := r.findBiz(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}
	return results, nil
}

func (r *merchantBusinessQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoActiveRow, error) {
	listRows, err := r.findBiz(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, merchantbusiness_errors.ErrFindActiveMerchantBusinesses
	}
	var active []*salemongo.MerchantBusinessInfoActiveRow
	for _, row := range listRows {
		active = append(active, &salemongo.MerchantBusinessInfoActiveRow{
			MerchantBusinessInfoRow: row.MerchantBusinessInfoRow,
			TotalCount:              row.TotalCount,
		})
	}
	return active, nil
}

func (r *merchantBusinessQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantBusinessInfoTrashedRow, error) {
	listRows, err := r.findBiz(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Page, req.PageSize)
	if err != nil {
		return nil, merchantbusiness_errors.ErrFindTrashedMerchantBusinesses
	}
	var trashed []*salemongo.MerchantBusinessInfoTrashedRow
	for _, row := range listRows {
		trashed = append(trashed, &salemongo.MerchantBusinessInfoTrashedRow{
			MerchantBusinessInfoRow: row.MerchantBusinessInfoRow,
			TotalCount:              row.TotalCount,
		})
	}
	return trashed, nil
}

func (r *merchantBusinessQueryRepository) FindByID(ctx context.Context, merchantID int) (*salemongo.MerchantBusinessInfoRow, error) {
	var b salemongo.MerchantBusinessInformation
	if err := r.coll().FindOne(ctx, bson.M{"merchant_id": int32(merchantID), "deleted_at": nil}).Decode(&b); err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}
	return bizToRow(&b), nil
}
