package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userRoleHandler struct {
	pbuserroles.UnimplementedUserRoleServiceServer
	userRole service.UserRoleService
	logger   logger.LoggerInterface
}

func NewUserRoleHandler(userRole service.UserRoleService, logger logger.LoggerInterface) pbuserroles.UserRoleServiceServer {
	return &userRoleHandler{
		userRole: userRole,
		logger:   logger,
	}
}

func (s *userRoleHandler) AssignRoleToUser(ctx context.Context, request *pbuserroles.AssignRoleToUserRequest) (*pbuserroles.ApiResponseUserRole, error) {
	req := &requests.CreateUserRoleRequest{
		UserId: int(request.GetUserId()),
		RoleId: int(request.GetRoleId()),
	}

	userRole, err := s.userRole.AssignRoleToUser(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbuserroles.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapToProtoUserRoleResponse(userRole),
	}, nil
}

func (s *userRoleHandler) RemoveRoleFromUser(ctx context.Context, request *pbuserroles.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	req := &requests.RemoveUserRoleRequest{
		UserId: int(request.GetUserId()),
		RoleId: int(request.GetRoleId()),
	}

	err := s.userRole.RemoveRoleFromUser(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func mapToProtoUserRoleResponse(v *salemongo.UserRole) *pbuserroles.UserRoleResponse {
	if v == nil {
		return nil
	}
	return &pbuserroles.UserRoleResponse{
		UserRoleId: v.UserRoleID,
		UserId:     v.UserID,
		RoleId:     v.RoleID,
		CreatedAt:  formatTimestamp(v.CreatedAt),
		UpdatedAt:  formatTimestamp(v.UpdatedAt),
	}
}
