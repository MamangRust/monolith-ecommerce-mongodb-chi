package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/dto"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewQueryRepository struct {
	coll *mongo.Collection
}

func NewReviewQueryRepository(db *mongo.Database) ReviewQueryRepository {
	return &reviewQueryRepository{coll: db.Collection("reviews")}
}

func fmtReviewTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05.000")
	}
	return ""
}

func reviewToResult(v *salemongo.Review, totalCount int64) *ReviewResult {
	return &ReviewResult{
		ReviewID:   v.ReviewID,
		UserID:     v.UserID,
		ProductID:  v.ProductID,
		Name:       v.Name,
		Comment:    v.Comment,
		Rating:     v.Rating,
		CreatedAt:  strPtrRev(fmtReviewTime(v.CreatedAt)),
		UpdatedAt:  strPtrRev(fmtReviewTime(v.UpdatedAt)),
		DeletedAt:  strPtrRev(fmtReviewTime(v.DeletedAt)),
		TotalCount: totalCount,
	}
}

func strPtrRev(s string) *string {
	return &s
}

func (r *reviewQueryRepository) findReviews(ctx context.Context, filter bson.M, page, pageSize int) ([]*ReviewResult, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "review_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*salemongo.Review
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	var results []*ReviewResult
	for _, v := range docs {
		results = append(results, reviewToResult(v, total))
	}
	return results, nil
}

func reviewFilter(req *requests.FindAllReview, trashed bool) bson.M {
	filter := bson.M{}
	if trashed {
		filter["deleted_at"] = bson.M{"$ne": nil}
	} else {
		filter["deleted_at"] = nil
	}
	if req.Search != "" {
		filter["comment"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	return filter
}

func (r *reviewQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	results, err := r.findReviews(ctx, reviewFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, review_errors.ErrFindAllReviews.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*ReviewResult, error) {
	filter := bson.M{"product_id": int32(req.ProductID), "deleted_at": nil}
	if req.Search != "" {
		filter["comment"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	results, err := r.findReviews(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, review_errors.ErrFindReviewsByProduct.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*ReviewResult, error) {
	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["comment"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	results, err := r.findReviews(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, review_errors.ErrFindReviewsByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	results, err := r.findReviews(ctx, reviewFilter(req, false), req.Page, req.PageSize)
	if err != nil {
		return nil, review_errors.ErrFindActiveReviews.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	results, err := r.findReviews(ctx, reviewFilter(req, true), req.Page, req.PageSize)
	if err != nil {
		return nil, review_errors.ErrFindTrashedReviews.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindByID(ctx context.Context, id int) (*ReviewResult, error) {
	var v salemongo.Review
	if err := r.coll.FindOne(ctx, bson.M{"review_id": int32(id), "deleted_at": nil}).Decode(&v); err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	return reviewToResult(&v, 1), nil
}

var _ = dto.GetUserByIDRow{}
