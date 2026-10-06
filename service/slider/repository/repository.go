package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	SliderCommand SliderCommandRepository
	SliderQuery   SliderQueryRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	coll := db.Collection("sliders")
	return &Repositories{
		SliderCommand: NewSliderCommandRepository(db, coll),
		SliderQuery:   NewSliderQueryRepository(coll),
	}
}
