package seeder

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type bannerSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewBannerSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *bannerSeeder {
	return &bannerSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *bannerSeeder) Seed() error {
	banners := []interface{}{
		salemongo.Banner{Name: "Banner 1", StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC), StartTime: parseTime("08:00").Format("15:04:05"), EndTime: parseTime("16:00").Format("15:04:05"), IsActive: toBoolPtr(true)},
		salemongo.Banner{Name: "Banner 2", StartDate: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC), StartTime: parseTime("09:00").Format("15:04:05"), EndTime: parseTime("17:00").Format("15:04:05"), IsActive: toBoolPtr(true)},
		salemongo.Banner{Name: "Banner 3", StartDate: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC), StartTime: parseTime("10:00").Format("15:04:05"), EndTime: parseTime("18:00").Format("15:04:05"), IsActive: toBoolPtr(false)},
		salemongo.Banner{Name: "Banner 4", StartDate: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC), StartTime: parseTime("07:00").Format("15:04:05"), EndTime: parseTime("15:00").Format("15:04:05"), IsActive: toBoolPtr(true)},
		salemongo.Banner{Name: "Banner 5", StartDate: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC), StartTime: parseTime("06:00").Format("15:04:05"), EndTime: parseTime("14:00").Format("15:04:05"), IsActive: toBoolPtr(false)},
		salemongo.Banner{Name: "Banner 6", StartDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC), StartTime: parseTime("12:00").Format("15:04:05"), EndTime: parseTime("20:00").Format("15:04:05"), IsActive: toBoolPtr(true)},
		salemongo.Banner{Name: "Banner 7", StartDate: time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC), StartTime: parseTime("08:30").Format("15:04:05"), EndTime: parseTime("16:30").Format("15:04:05"), IsActive: toBoolPtr(true)},
		salemongo.Banner{Name: "Banner 8", StartDate: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC), StartTime: parseTime("09:30").Format("15:04:05"), EndTime: parseTime("17:30").Format("15:04:05"), IsActive: toBoolPtr(false)},
	}

	return insertMany(r.ctx, r.db, "banners", banners, r.logger)
}

func parseTime(t string) time.Time {
	parsed, _ := time.Parse("15:04", t)
	return parsed
}
