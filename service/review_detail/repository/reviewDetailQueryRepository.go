package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewDetailQueryRepository struct {
	db *mongo.Database
}

func NewReviewDetailQueryRepository(db *mongo.Database) ReviewDetailQueryRepository {
	return &reviewDetailQueryRepository{db: db}
}

func (r *reviewDetailQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("review_details")
}

func reviewDetailToListRow(d *salemongo.ReviewDetail, total int64) *salemongo.ReviewDetailListRow {
	return &salemongo.ReviewDetailListRow{
		ReviewDetailRow: *reviewDetailToRow(d),
		TotalCount:      total,
	}
}

func (r *reviewDetailQueryRepository) findDetails(ctx context.Context, filter bson.M, page, pageSize int) ([]*salemongo.ReviewDetailListRow, *int, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindAllReviewDetails
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "review_detail_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindAllReviewDetails
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.ReviewDetail
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, reviewdetail_errors.ErrFindAllReviewDetails
	}

	var results []*salemongo.ReviewDetailListRow
	for _, d := range docs {
		results = append(results, reviewDetailToListRow(d, total))
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	return results, &totalPages, nil
}

func reviewDetailFilter(req *requests.FindAllReview, trashed bool) bson.M {
	filter := bson.M{}
	if trashed {
		filter["deleted_at"] = bson.M{"$ne": nil}
	} else {
		filter["deleted_at"] = nil
	}
	if req.Search != "" {
		filter["type"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	return filter
}

func (r *reviewDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailListRow, *int, error) {
	return r.findDetails(ctx, reviewDetailFilter(req, false), req.Page, req.PageSize)
}

func (r *reviewDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailActiveRow, *int, error) {
	listRows, totalPages, err := r.findDetails(ctx, reviewDetailFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindActiveReviewDetails
	}
	var active []*salemongo.ReviewDetailActiveRow
	for _, row := range listRows {
		active = append(active, &salemongo.ReviewDetailActiveRow{
			ReviewDetailRow: row.ReviewDetailRow,
			TotalCount:      row.TotalCount,
		})
	}
	return active, totalPages, nil
}

func (r *reviewDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*salemongo.ReviewDetailTrashedRow, *int, error) {
	listRows, totalPages, err := r.findDetails(ctx, reviewDetailFilter(req, true), req.Page, req.PageSize)
	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindTrashedReviewDetails
	}
	var trashed []*salemongo.ReviewDetailTrashedRow
	for _, row := range listRows {
		trashed = append(trashed, &salemongo.ReviewDetailTrashedRow{
			ReviewDetailRow: row.ReviewDetailRow,
			TotalCount:      row.TotalCount,
		})
	}
	return trashed, totalPages, nil
}

func (r *reviewDetailQueryRepository) FindByID(ctx context.Context, userID int) (*salemongo.ReviewDetailRow, error) {
	var d salemongo.ReviewDetail
	if err := r.coll().FindOne(ctx, bson.M{"review_detail_id": int32(userID), "deleted_at": nil}).Decode(&d); err != nil {
		return nil, reviewdetail_errors.ErrReviewDetailNotFound
	}
	return reviewDetailToRow(&d), nil
}

func (r *reviewDetailQueryRepository) FindByIDTrashed(ctx context.Context, userID int) (*salemongo.ReviewDetail, error) {
	var d salemongo.ReviewDetail
	if err := r.coll().FindOne(ctx, bson.M{"review_detail_id": int32(userID), "deleted_at": bson.M{"$ne": nil}}).Decode(&d); err != nil {
		return nil, reviewdetail_errors.ErrReviewDetailNotFound
	}
	return &d, nil
}
