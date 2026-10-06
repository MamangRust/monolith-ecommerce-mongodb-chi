package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	RoleQuery   pbroles.RoleQueryServiceServer
	RoleCommand pbroles.RoleCommandServiceServer
	UserRole    pbuserroles.UserRoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		RoleQuery:   NewRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		RoleCommand: NewRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
		UserRole:    NewUserRoleHandler(deps.Service.UserRole, deps.Logger),
	}
}
