package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repositories struct {
	CategoryQuery           CategoryQueryRepository
	CategoryCommand         CategoryCommandRepository
	CategoryStats           CategoryStatsRepository
	CategoryStatsById       CategoryStatsByIdRepository
	CategoryStatsByMerchant CategoryStatsByMerchantRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	coll := db.Collection("categories")
	return &Repositories{
		CategoryQuery:           NewCategoryQueryRepository(coll),
		CategoryCommand:         NewCategoryCommandRepository(db, coll),
		CategoryStats:           NewCategoryStatsRepository(db),
		CategoryStatsById:       NewCategoryStatsByIdRepository(db),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantRepository(db),
	}
}
