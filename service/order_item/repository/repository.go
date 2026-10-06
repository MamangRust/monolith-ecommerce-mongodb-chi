package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repositories struct {
	OrderItemQuery   OrderItemQueryRepository
	OrderItemCommand OrderItemCommandRepository
}

func NewRepositories(DB *mongo.Database) *Repositories {
	return &Repositories{
		OrderItemQuery:   NewOrderItemQueryRepository(DB),
		OrderItemCommand: NewOrderItemCommandRepository(DB),
	}
}
