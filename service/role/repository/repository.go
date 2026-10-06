package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
	UserRole    UserRoleRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	roleColl := db.Collection("roles")
	userRoleColl := db.Collection("user_roles")

	return &Repositories{
		RoleCommand: NewRoleCommandRepository(db, roleColl),
		RoleQuery:   NewRoleQueryRepository(db, roleColl, userRoleColl),
		UserRole:    NewUserRoleRepository(db, userRoleColl),
	}
}
