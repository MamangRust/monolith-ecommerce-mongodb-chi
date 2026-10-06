// Package role adapts the Role service gRPC API into the shared domain model.
// It is the only place that talks to the role service for role lookups; the
// (un)assignment of a user's roles lives in the user_role adapter.
package role

import (
	"context"
	"time"

	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	role_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
)

// QueryRepository is the contract consumers depend on for role reads.
type QueryRepository interface {
	FindByID(ctx context.Context, roleID int) (*mongo.Role, error)
	FindByName(ctx context.Context, name string) (*mongo.Role, error)
}

// Repository implements QueryRepository on top of the generated role query
// client.
type Repository struct {
	query pbroles.RoleQueryServiceClient
	guard *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) { r.guard = g }

// New wraps the generated role query client. Resilience (timeout/circuit-
// breaker/bulkhead) is applied per call through the optional dependency guard;
// a nil guard is a passthrough. Modules build it inside their repositories with
// the guard options they were given.
func New(query pbroles.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface.
func NewQueryAdapter(query pbroles.RoleQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, opts...)
}

func (r *Repository) FindByID(ctx context.Context, roleID int) (*mongo.Role, error) {
	res, err := adapter.Call(r.guard, ctx, func(ctx context.Context) (*pbroles.ApiResponseRole, error) {
		return r.query.FindByIdRole(ctx, &pbroles.FindByIdRoleRequest{RoleId: int32(roleID)})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return roleToModel(res.Data), nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (*mongo.Role, error) {
	res, err := adapter.Call(r.guard, ctx, func(ctx context.Context) (*pbroles.ApiResponseRole, error) {
		return r.query.FindByNameRole(ctx, &pbroles.FindByNameRoleRequest{Name: name})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return roleToModel(res.Data), nil
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

// parseTime accepts the loose RFC3339-ish timestamp strings the role service
// emits and returns nil for empty or unparseable values (the mongo models use
// nullable timestamps).
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
