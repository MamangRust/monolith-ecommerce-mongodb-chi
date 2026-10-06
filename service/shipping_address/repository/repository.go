package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	ShippingAddressQuery   ShippingAddressQueryRepository
	ShippingAddressCommand ShippingAddressCommandRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		ShippingAddressQuery:   NewShippingAddressQueryRepository(db),
		ShippingAddressCommand: NewShippingAddressCommandRepository(db),
	}
}
