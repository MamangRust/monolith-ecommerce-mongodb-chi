package transaction

import (
	"context"
	"fmt"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CommandRepository is the contract consumers depend on for transaction writes.
type CommandRepository interface {
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type commandAdapter struct {
	client pbtransaction.TransactionCommandServiceClient
	guard  *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *commandAdapter) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// New wraps the generated transaction command client. It is the canonical
// constructor; modules build it inside their repositories with the guard
// options they were given.
func New(client pbtransaction.TransactionCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return NewCommandAdapter(client, opts...)
}

// NewCommandAdapter wraps the generated transaction command client. Resilience
// (timeout/circuit-breaker/bulkhead) is applied per call through the optional
// dependency guard; a nil guard is a passthrough.
func NewCommandAdapter(client pbtransaction.TransactionCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	a := &commandAdapter{client: client}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *commandAdapter) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	if a.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	if _, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbtransaction.ApiResponseTransactionDelete, error) {
		return a.client.DeleteTransactionByOrderPermanent(ctx, &pbtransaction.FindByIdTransactionRequest{
			Id: int32(orderID),
		})
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (a *commandAdapter) DeleteAll(ctx context.Context) (bool, error) {
	if a.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbtransaction.ApiResponseTransactionAll, error) {
		return a.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	})
	if err != nil {
		return false, err
	}

	return res.Status == "success", nil
}
