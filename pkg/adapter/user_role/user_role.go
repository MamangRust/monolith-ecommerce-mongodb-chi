// Package user_role adapts the role service's user-role RPCs into the shared
// domain model. It owns the only place that assigns or removes roles for a
// user and resolves the roles attached to a user.
package user_role

import (
	"context"
	"time"

	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_role_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the read path consumers use to resolve a user's roles.
type QueryRepository interface {
	FindByUserId(ctx context.Context, userID int) ([]*mongo.Role, error)
}

// CommandRepository is the write path consumers use to (un)assign roles.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*mongo.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements QueryRepository and CommandRepository on top of the
// generated role query client (reads) and user-role client (writes).
type Repository struct {
	roleQueryClient pbroles.RoleQueryServiceClient
	userRoleClient  pbuserroles.UserRoleServiceClient
	guard           *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) { r.guard = g }

// New wraps the generated role query and user-role clients. Resilience
// (timeout/circuit-breaker/bulkhead) is applied per call through the optional
// dependency guard; a nil guard is a passthrough. Modules build it inside their
// repositories with the guard options they were given.
func New(roleQueryClient pbroles.RoleQueryServiceClient, userRoleClient pbuserroles.UserRoleServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{roleQueryClient: roleQueryClient, userRoleClient: userRoleClient}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// FindByUserId implements QueryRepository. The role service returns the resolved
// roles, so we map the role responses straight into the domain model.
func (r *Repository) FindByUserId(ctx context.Context, userID int) ([]*mongo.Role, error) {
	res, err := adapter.Call(r.guard, ctx, func(ctx context.Context) (*pbroles.ApiResponsesRole, error) {
		return r.roleQueryClient.FindByUserId(ctx, &pbroles.FindByIdUserRoleRequest{UserId: int32(userID)})
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	roles := make([]*mongo.Role, 0, len(res.Data))
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		roles = append(roles, roleToModel(item))
	}
	return roles, nil
}

// AssignRoleToUser implements CommandRepository. The user-role service treats a
// repeated assignment as a no-op and returns the existing mapping, so this call
// is idempotent.
func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*mongo.UserRole, error) {
	res, err := adapter.Call(r.guard, ctx, func(ctx context.Context) (*pbuserroles.ApiResponseUserRole, error) {
		return r.userRoleClient.AssignRoleToUser(ctx, &pbuserroles.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
	})
	if err != nil || res == nil || res.Data == nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	return userRoleToModel(res.Data), nil
}

// RemoveRoleFromUser implements CommandRepository.
func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	_, err := adapter.Call(r.guard, ctx, func(ctx context.Context) (*emptypb.Empty, error) {
		return r.userRoleClient.RemoveRoleFromUser(ctx, &pbuserroles.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}

func roleToModel(role *pbroles.RoleResponse) *mongo.Role {
	if role == nil {
		return nil
	}
	return &mongo.Role{
		RoleID:    role.Id,
		RoleName:  role.Name,
		CreatedAt: parseTime(role.CreatedAt),
		UpdatedAt: parseTime(role.UpdatedAt),
	}
}

func userRoleToModel(ur *pbuserroles.UserRoleResponse) *mongo.UserRole {
	if ur == nil {
		return nil
	}
	return &mongo.UserRole{
		UserRoleID: ur.UserRoleId,
		UserID:     ur.UserId,
		RoleID:     ur.RoleId,
	}
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}
