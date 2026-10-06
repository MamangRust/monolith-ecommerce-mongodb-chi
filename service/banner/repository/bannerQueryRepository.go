package repository

import (
	"context"
	"regexp"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bannerQueryRepository struct {
	coll *mongo.Collection
}

func NewBannerQueryRepository(coll *mongo.Collection) *bannerQueryRepository {
	return &bannerQueryRepository{coll: coll}
}

func bannerSearchFilter(search string) bson.M {
	rx := bson.M{"$regex": regexp.QuoteMeta(search), "$options": "i"}
	return bson.M{"name": rx}
}

func bannerToResult(b salemongo.Banner) *BannerResult {
	startDate := b.StartDate.Format("2006-01-02")
	endDate := b.EndDate.Format("2006-01-02")
	startTime := b.StartTime
	endTime := b.EndTime
	createdAt, updatedAt, deletedAt := "", "", ""
	if b.CreatedAt != nil {
		createdAt = b.CreatedAt.Format("2006-01-02 15:04:05.000")
	}
	if b.UpdatedAt != nil {
		updatedAt = b.UpdatedAt.Format("2006-01-02 15:04:05.000")
	}
	if b.DeletedAt != nil {
		deletedAt = b.DeletedAt.Format("2006-01-02 15:04:05.000")
	}
	return &BannerResult{
		BannerID:  b.BannerID,
		Name:      b.Name,
		StartDate: &startDate,
		EndDate:   &endDate,
		StartTime: &startTime,
		EndTime:   &endTime,
		IsActive:  b.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
		DeletedAt: &deletedAt,
	}
}

func (r *bannerQueryRepository) findPage(ctx context.Context, filter bson.M, search string, page, pageSize int) ([]*BannerResult, error) {
	if search != "" {
		filter["name"] = bson.M{"$regex": regexp.QuoteMeta(search), "$options": "i"}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "banner_id", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []*BannerResult
	for cursor.Next(ctx) {
		var b salemongo.Banner
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		r := bannerToResult(b)
		r.TotalCount = total
		results = append(results, r)
	}
	return results, nil
}

func (r *bannerQueryRepository) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	return r.findPage(ctx, bson.M{"deleted_at": nil}, req.Search, req.Page, req.PageSize)
}

func (r *bannerQueryRepository) FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	filter := bson.M{"deleted_at": nil, "is_active": true}
	return r.findPage(ctx, filter, req.Search, req.Page, req.PageSize)
}

func (r *bannerQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	return r.findPage(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Search, req.Page, req.PageSize)
}

func (r *bannerQueryRepository) FindByID(ctx context.Context, bannerID int) (*BannerResult, error) {
	var b salemongo.Banner
	err := r.coll.FindOne(ctx, bson.M{"banner_id": bannerID, "deleted_at": nil}).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, banner_errors.ErrBannerNotFound
		}
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}
	return bannerToResult(b), nil
}

// keep time import
var _ = time.Now
