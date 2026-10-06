package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type UserRoleService interface {
	AssignRoleToUser(ctx context.Context, request *requests.CreateUserRoleRequest) (*salemongo.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, request *requests.RemoveUserRoleRequest) error
}

type userRoleService struct {
	observability      observability.TraceLoggerObservability
	userRoleRepository repository.UserRoleRepository
	logger             logger.LoggerInterface
}

type UserRoleServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	UserRoleRepository repository.UserRoleRepository
	Logger             logger.LoggerInterface
}

func NewUserRoleService(deps *UserRoleServiceDeps) UserRoleService {
	return &userRoleService{
		observability:      deps.Observability,
		userRoleRepository: deps.UserRoleRepository,
		logger:             deps.Logger,
	}
}

func (s *userRoleService) AssignRoleToUser(ctx context.Context, request *requests.CreateUserRoleRequest) (*salemongo.UserRole, error) {
	const method = "AssignRoleToUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", request.UserId),
		attribute.Int("role_id", request.RoleId))

	defer func() {
		end(status)
	}()

	userRole, err := s.userRoleRepository.AssignRoleToUser(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*salemongo.UserRole](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", request.UserId),
			zap.Int("role_id", request.RoleId),
		)
	}

	logSuccess("Successfully assigned role to user", zap.Int("user_id", request.UserId), zap.Int("role_id", request.RoleId))

	return userRole, nil
}

func (s *userRoleService) RemoveRoleFromUser(ctx context.Context, request *requests.RemoveUserRoleRequest) error {
	const method = "RemoveRoleFromUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", request.UserId),
		attribute.Int("role_id", request.RoleId))

	defer func() {
		end(status)
	}()

	err := s.userRoleRepository.RemoveRoleFromUser(ctx, request)
	if err != nil {
		status = "error"
		_, err := errorhandler.HandleError[any](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", request.UserId),
			zap.Int("role_id", request.RoleId),
		)
		return err
	}

	logSuccess("Successfully removed role from user", zap.Int("user_id", request.UserId), zap.Int("role_id", request.RoleId))

	return nil
}
