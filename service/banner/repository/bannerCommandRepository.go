package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bannerCommandRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewBannerCommandRepository(db *mongo.Database, coll *mongo.Collection) *bannerCommandRepository {
	return &bannerCommandRepository{db: db, coll: coll}
}

func timePtr(t time.Time) *time.Time { return &t }

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func parseDatePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}
	return &t, nil
}

func parseTimeStr(s string) (*string, error) {
	if s == "" {
		return nil, nil
	}
	_, err := time.Parse("15:04:05", s)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}
	return &s, nil
}

func (r *bannerCommandRepository) Create(ctx context.Context, request *requests.CreateBannerRequest) (*salemongo.Banner, error) {
	startDate, err := parseDatePtr(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}
	endDate, err := parseDatePtr(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}
	startTime, err := parseTimeStr(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}
	endTime, err := parseTimeStr(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	bannerID, err := salemongo.NextID(ctx, r.db, "banners")
	if err != nil {
		return nil, banner_errors.ErrCreateBanner.WithInternal(err)
	}

	now := time.Now()
	banner := &salemongo.Banner{
		BannerID:  bannerID,
		Name:      request.Name,
		StartDate: derefTime(startDate),
		EndDate:   derefTime(endDate),
		StartTime: derefStr(startTime),
		EndTime:   derefStr(endTime),
		IsActive:  &request.IsActive,
		CreatedAt: timePtr(now),
		UpdatedAt: timePtr(now),
	}

	if _, err := r.coll.InsertOne(ctx, banner); err != nil {
		return nil, banner_errors.ErrCreateBanner.WithInternal(err)
	}
	return banner, nil
}

func (r *bannerCommandRepository) Update(ctx context.Context, request *requests.UpdateBannerRequest) (*salemongo.Banner, error) {
	startDate, err := parseDatePtr(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}
	endDate, err := parseDatePtr(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}
	startTime, err := parseTimeStr(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}
	endTime, err := parseTimeStr(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"banner_id": *request.BannerID, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"name":       request.Name,
			"start_date": startDate,
			"end_date":   endDate,
			"start_time": startTime,
			"end_time":   endTime,
			"is_active":  request.IsActive,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var banner salemongo.Banner
	if err := result.Decode(&banner); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, banner_errors.ErrBannerNotFound
		}
		return nil, banner_errors.ErrUpdateBanner.WithInternal(err)
	}
	return &banner, nil
}

func (r *bannerCommandRepository) Trash(ctx context.Context, bannerID int) (*salemongo.Banner, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"banner_id": bannerID, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var banner salemongo.Banner
	if err := result.Decode(&banner); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, banner_errors.ErrBannerNotFound
		}
		return nil, banner_errors.ErrTrashedBanner.WithInternal(err)
	}
	return &banner, nil
}

func (r *bannerCommandRepository) Restore(ctx context.Context, bannerID int) (*salemongo.Banner, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"banner_id": bannerID, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var banner salemongo.Banner
	if err := result.Decode(&banner); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, banner_errors.ErrBannerNotFound
		}
		return nil, banner_errors.ErrRestoreBanner.WithInternal(err)
	}
	return &banner, nil
}

func (r *bannerCommandRepository) DeletePermanent(ctx context.Context, bannerID int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"banner_id": bannerID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, banner_errors.ErrDeleteBannerPermanent.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, banner_errors.ErrBannerNotFound
	}
	return true, nil
}

func (r *bannerCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, banner_errors.ErrRestoreAllBanners.WithInternal(err)
	}
	return true, nil
}

func (r *bannerCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, banner_errors.ErrDeleteAllBanners.WithInternal(err)
	}
	return true, nil
}
