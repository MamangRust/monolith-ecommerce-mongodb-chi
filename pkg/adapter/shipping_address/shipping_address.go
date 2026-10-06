package shipping_address

import (
	"context"

	pbshipping "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the contract consumers depend on for shipping-address reads.
type QueryRepository interface {
	FindByID(ctx context.Context, shippingID int) (*mongo.ShippingAddress, error)
	FindByOrder(ctx context.Context, orderID int) (*mongo.ShippingAddress, error)
}

// CommandRepository is the contract consumers depend on for shipping-address writes.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*mongo.ShippingAddress, error)
	Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*mongo.ShippingAddress, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated shipping query/command clients.
type Repository struct {
	query   pbshipping.ShippingQueryServiceClient
	command pbshipping.ShippingCommandServiceClient
	guard   *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *Repository) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// NewAdapter wraps the generated shipping clients. Resilience
// (timeout/circuit-breaker/bulkhead) is applied per call through the optional
// dependency guard; a nil guard is a passthrough.
func NewAdapter(query pbshipping.ShippingQueryServiceClient, command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	a := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// New wraps the generated shipping query/command clients. It is the canonical
// constructor; modules build it inside their repositories with the guard
// options they were given.
func New(query pbshipping.ShippingQueryServiceClient, command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return NewAdapter(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface.
func NewQueryAdapter(query pbshipping.ShippingQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewAdapter(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository surface.
func NewCommandAdapter(command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return NewAdapter(nil, command, opts...)
}

// Reads propagate the dependency's gRPC status unchanged (NotFound -> 404, and
// so on). Command failures keep the domain-error wrapping the order service
// already relied on; AppError.Unwrap still exposes the underlying gRPC status.
func (a *Repository) FindByID(ctx context.Context, shippingID int) (*mongo.ShippingAddress, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return a.query.FindById(ctx, &pbshipping.FindByIdShippingRequest{Id: int32(shippingID)})
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) FindByOrder(ctx context.Context, orderID int) (*mongo.ShippingAddress, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return a.query.FindByOrder(ctx, &pbshipping.FindByIdShippingRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*mongo.ShippingAddress, error) {
	var orderID int32
	if req.OrderID != nil {
		orderID = int32(*req.OrderID)
	}

	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return a.command.CreateShipping(ctx, &pbshipping.CreateShippingAddressRequest{
			OrderId:        orderID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*mongo.ShippingAddress, error) {
	var shippingID int32
	if req.ShippingID != nil {
		shippingID = int32(*req.ShippingID)
	}

	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return a.command.UpdateShipping(ctx, &pbshipping.UpdateShippingAddressRequest{
			ShippingId:     shippingID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbshipping.ApiResponseShippingDelete, error) {
		return a.command.DeleteShippingByOrderPermanent(ctx, &pbshipping.FindByIdShippingRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent.WithInternal(err)
	}
	return true, nil
}

func (a *Repository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbshipping.ApiResponseShippingAll, error) {
		return a.command.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress.WithInternal(err)
	}
	return res.Status == "success", nil
}

func toModel(data *pbshipping.ShippingResponse) *mongo.ShippingAddress {
	if data == nil {
		return nil
	}
	return &mongo.ShippingAddress{
		ShippingAddressID: data.Id,
		OrderID:           data.OrderId,
		Alamat:            data.Alamat,
		Provinsi:          data.Provinsi,
		Negara:            data.Negara,
		Kota:              data.Kota,
		Courier:           data.Courier,
		ShippingMethod:    data.ShippingMethod,
		ShippingCost:      float64(data.ShippingCost),
	}
}
