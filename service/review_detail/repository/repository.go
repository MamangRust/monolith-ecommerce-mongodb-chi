package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	ReviewDetailQuery   ReviewDetailQueryRepository
	ReviewDetailCommand ReviewDetailCommandRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		ReviewDetailQuery:   NewReviewDetailQueryRepository(db),
		ReviewDetailCommand: NewReviewDetailCommandRepository(db),
	}
}
