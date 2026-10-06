package merchant

import (
	"context"
	"time"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
)

// QueryRepository is the contract consumers depend on for merchant reads.
// It is implemented by grpcQueryAdapter and can be satisfied by a test double.
type QueryRepository interface {
	FindByID(ctx context.Context, merchantID int) (*mongo.Merchant, error)
}

type grpcQueryAdapter struct {
	client pbmerchant.MerchantQueryServiceClient
	guard  *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *grpcQueryAdapter) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// New wraps a generated merchant query client into a QueryRepository. It is the
// canonical constructor; modules build it inside their repositories with the
// guard options they were given.
func New(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewQueryAdapter(client, opts...)
}

// NewQueryAdapter wraps a generated merchant query client into a QueryRepository.
// Resilience (timeout/circuit-breaker/bulkhead) is applied per call through the
// optional dependency guard; a nil guard is a passthrough.
func NewQueryAdapter(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	a := &grpcQueryAdapter{client: client}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *grpcQueryAdapter) FindByID(ctx context.Context, merchantID int) (*mongo.Merchant, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbmerchant.ApiResponseMerchant, error) {
		return a.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{Id: int32(merchantID)})
	})
	if err != nil {
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func toModel(m *pbmerchant.MerchantResponse) *mongo.Merchant {
	if m == nil {
		return nil
	}
	return &mongo.Merchant{
		MerchantID:   m.Id,
		UserID:       m.UserId,
		Name:         m.Name,
		Description:  &m.Description,
		Address:      &m.Address,
		ContactEmail: &m.ContactEmail,
		ContactPhone: &m.ContactPhone,
		Status:       m.Status,
		CreatedAt:    parseTimeString(m.CreatedAt),
		UpdatedAt:    parseTimeString(m.UpdatedAt),
	}
}

func parseTimeString(s string) *time.Time {
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
