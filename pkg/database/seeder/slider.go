package seeder

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type sliderSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewSliderSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *sliderSeeder {
	return &sliderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *sliderSeeder) Seed() error {
	sliders := []salemongo.Slider{
		{Name: "Promo Akhir Tahun", Image: "slider1.jpg"},
		{Name: "Diskon Elektronik", Image: "slider2.jpg"},
		{Name: "Flash Sale Mingguan", Image: "slider3.jpg"},
		{Name: "Produk Terbaru", Image: "slider4.jpg"},
		{Name: "Promo Kesehatan", Image: "slider5.jpg"},
		{Name: "Belanja Hemat", Image: "slider6.jpg"},
		{Name: "Gaming Gear Diskon", Image: "slider7.jpg"},
		{Name: "Gratis Ongkir", Image: "slider8.jpg"},
		{Name: "Ramadhan Sale", Image: "slider9.jpg"},
		{Name: "Perlengkapan Bayi", Image: "slider10.jpg"},
	}

	for _, slider := range sliders {
		if _, err := r.db.Collection("sliders").InsertOne(r.ctx, slider); err != nil {
			r.logger.Error("failed to seed slider", zap.Error(err))
			return err
		}
	}

	r.logger.Info("slider successfully seeded")

	return nil
}
