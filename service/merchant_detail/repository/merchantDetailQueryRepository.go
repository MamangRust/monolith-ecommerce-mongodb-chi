package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantDetailQueryRepository struct {
	db *mongo.Database
}

func NewMerchantDetailQueryRepository(db *mongo.Database) MerchantDetailQueryRepository {
	return &merchantDetailQueryRepository{db: db}
}

func (r *merchantDetailQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_details")
}

func detailToListRow(d *salemongo.MerchantDetail, total int64, merchantName string) *salemongo.MerchantDetailListRow {
	return &salemongo.MerchantDetailListRow{
		MerchantDetailRow: *detailToRow(d),
		MerchantName:      merchantName,
		TotalCount:        total,
	}
}

func (r *merchantDetailQueryRepository) findDetails(ctx context.Context, filter bson.M, page, pageSize int) ([]*salemongo.MerchantDetailListRow, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "merchant_detail_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.MerchantDetail
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails
	}

	var results []*salemongo.MerchantDetailListRow
	for _, d := range docs {
		results = append(results, detailToListRow(d, total, ""))
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailListRow, error) {
	results, err := r.findDetails(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailActiveRow, error) {
	listRows, err := r.findDetails(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindActiveMerchantDetails
	}
	var active []*salemongo.MerchantDetailActiveRow
	for _, row := range listRows {
		active = append(active, &salemongo.MerchantDetailActiveRow{
			MerchantDetailRow: row.MerchantDetailRow,
			TotalCount:        row.TotalCount,
			MerchantName:      row.MerchantName,
		})
	}
	return active, nil
}

func (r *merchantDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantDetailTrashedRow, error) {
	listRows, err := r.findDetails(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Page, req.PageSize)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindTrashedMerchantDetails
	}
	var trashed []*salemongo.MerchantDetailTrashedRow
	for _, row := range listRows {
		trashed = append(trashed, &salemongo.MerchantDetailTrashedRow{
			MerchantDetailRow: row.MerchantDetailRow,
			TotalCount:        row.TotalCount,
			MerchantName:      row.MerchantName,
		})
	}
	return trashed, nil
}

func (r *merchantDetailQueryRepository) FindByID(ctx context.Context, userID int) (*salemongo.MerchantDetailWithJoinsRow, error) {
	// Look up merchant by user_id first, then the detail by merchant_id
	merchantColl := r.db.Collection("merchants")
	var merchant salemongo.Merchant
	if err := merchantColl.FindOne(ctx, bson.M{"user_id": int32(userID), "deleted_at": nil}).Decode(&merchant); err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}

	var d salemongo.MerchantDetail
	if err := r.coll().FindOne(ctx, bson.M{"merchant_id": merchant.MerchantID, "deleted_at": nil}).Decode(&d); err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}

	return &salemongo.MerchantDetailWithJoinsRow{
		MerchantDetailRow: *detailToRow(&d),
		MerchantName:      merchant.Name,
	}, nil
}
