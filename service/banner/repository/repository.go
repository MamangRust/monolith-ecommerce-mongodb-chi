package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	BannerCommand BannerCommandRepository
	BannerQuery   BannerQueryRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	bannerColl := db.Collection("banners")

	return &Repositories{
		BannerCommand: NewBannerCommandRepository(db, bannerColl),
		BannerQuery:   NewBannerQueryRepository(bannerColl),
	}
}
