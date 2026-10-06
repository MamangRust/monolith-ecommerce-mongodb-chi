package product

import (
	"context"
	"time"

	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/products"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	product_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)

// QueryRepository is the contract consumers depend on for product reads.
type QueryRepository interface {
	FindByID(ctx context.Context, productID int) (*mongo.Product, error)
}

// CommandRepository is the contract consumers depend on for product writes
// (stock adjustments). CreateProduct/UpdateProduct live entirely inside the
// owning service and are not exposed here.
type CommandRepository interface {
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*mongo.Product, error)
	AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*mongo.Product, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated product query/command clients.
type Repository struct {
	query   pbproduct.ProductQueryServiceClient
	command pbproduct.ProductCommandServiceClient
	guard   *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *Repository) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// NewAdapter wraps the generated product clients. Resilience (timeout/circuit-breaker/
// bulkhead) is applied per call through the optional dependency guard; a nil
// guard is a passthrough.
func NewAdapter(query pbproduct.ProductQueryServiceClient, command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	a := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// New wraps the generated product query/command clients. It is the canonical
// constructor; modules build it inside their repositories with the guard
// options they were given.
func New(query pbproduct.ProductQueryServiceClient, command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return NewAdapter(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface
// for consumers that only read (cart, review).
func NewQueryAdapter(query pbproduct.ProductQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewAdapter(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository surface
// for consumers that only write.
func NewCommandAdapter(command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return NewAdapter(nil, command, opts...)
}

func (a *Repository) FindByID(ctx context.Context, productID int) (*mongo.Product, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return a.query.FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(productID)})
	})
	if err != nil {
		return nil, product_errors.ErrProductNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*mongo.Product, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return a.command.UpdateProductCountStock(ctx, &pbproduct.UpdateProductCountStockRequest{
			ProductId: int32(productID),
			Stock:     int32(stock),
		})
	})
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*mongo.Product, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return a.command.AdjustProductStock(ctx, &pbproduct.AdjustProductStockRequest{
			ProductId:   int32(productID),
			Delta:       int32(delta),
			OperationId: operationID,
		})
	})
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func toModel(data *pbproduct.ProductResponse) *mongo.Product {
	if data == nil {
		return nil
	}
	description := data.Description
	weight := data.Weight
	image := data.ImageProduct
	return &mongo.Product{
		ProductID:    data.Id,
		MerchantID:   data.MerchantId,
		CategoryID:   data.CategoryId,
		Name:         data.Name,
		Description:  &description,
		Price:        data.Price,
		CountInStock: data.CountInStock,
		Weight:       &weight,
		ImageProduct: &image,
		CreatedAt:    parseTime(data.CreatedAt),
		UpdatedAt:    parseTime(data.UpdatedAt),
	}
}

// parseTime accepts the loose RFC3339-ish timestamp strings the product service
// emits and returns the zero time for empty or unparseable values.
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
