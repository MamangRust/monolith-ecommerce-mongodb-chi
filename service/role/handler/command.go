package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleCommandHandler struct {
	pbroles.UnimplementedRoleCommandServiceServer
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewRoleCommandHandler(roleCommand service.RoleCommandService, logger logger.LoggerInterface) pbroles.RoleCommandServiceServer {
	return &roleCommandHandler{
		roleCommand: roleCommand,
		logger:      logger,
	}
}

func (s *roleCommandHandler) CreateRole(ctx context.Context, request *pbroles.CreateRoleRequest) (*pbroles.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{
		Name: request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully created role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) UpdateRole(ctx context.Context, request *pbroles.UpdateRoleRequest) (*pbroles.ApiResponseRole, error) {
	id := int(request.GetId())
	req := &requests.UpdateRoleRequest{
		ID:   &id,
		Name: request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully updated role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) TrashedRole(ctx context.Context, request *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRole, error) {
	id := int(request.GetRoleId())
	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully trashed role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) RestoreRole(ctx context.Context, request *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRole, error) {
	id := int(request.GetRoleId())
	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully restored role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) DeleteRolePermanent(ctx context.Context, request *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRoleDelete, error) {
	id := int(request.GetRoleId())
	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRoleDelete{
		Status:  "success",
		Message: "Successfully deleted role permanently",
	}, nil
}

func (s *roleCommandHandler) RestoreAllRole(ctx context.Context, _ *emptypb.Empty) (*pbroles.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully restored all roles",
	}, nil
}

func (s *roleCommandHandler) DeleteAllRolePermanent(ctx context.Context, _ *emptypb.Empty) (*pbroles.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbroles.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully deleted all roles permanently",
	}, nil
}
