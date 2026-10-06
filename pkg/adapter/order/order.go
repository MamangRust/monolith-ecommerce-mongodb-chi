package order

import (
	"context"

	pborder "github.com/MamangRust/monolith-ecommerce-pb/orders"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
)

// QueryRepository is the contract consumers depend on for order reads.
type QueryRepository interface {
	FindByID(ctx context.Context, orderID int) (*mongo.Order, error)
}

type queryAdapter struct {
	client pborder.OrderQueryServiceClient
	guard  *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter so a dependency guard can be attached
// via adapter.WithDependencyGuard.
func (a *queryAdapter) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// New wraps the generated order query client. It is the canonical constructor;
// modules build it inside their repositories with the guard options they were
// given.
func New(client pborder.OrderQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewQueryAdapter(client, opts...)
}

// NewQueryAdapter wraps the generated order query client. Resilience
// (timeout/circuit-breaker/bulkhead) is applied per call through the optional
// dependency guard; a nil guard is a passthrough.
func NewQueryAdapter(client pborder.OrderQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	a := &queryAdapter{client: client}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// FindByID propagates the dependency's gRPC status unchanged (NotFound -> 404,
// and so on) instead of wrapping it in a domain error.
func (a *queryAdapter) FindByID(ctx context.Context, orderID int) (*mongo.Order, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pborder.ApiResponseOrder, error) {
		return a.client.FindById(ctx, &pborder.FindByIdOrderRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Data == nil {
		return nil, nil
	}

	return &mongo.Order{
		OrderID:    res.Data.Id,
		UserID:     res.Data.UserId,
		MerchantID: res.Data.MerchantId,
		TotalPrice: res.Data.TotalPrice,
	}, nil
}
