package category

import (
	"context"
	"time"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/categories"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

// QueryRepository is the contract consumers depend on for category reads.
// It is implemented by grpcQueryAdapter and can be satisfied by a test double.
type QueryRepository interface {
	FindByID(ctx context.Context, categoryID int) (*mongo.Category, error)
	FindAll(ctx context.Context, page, pageSize int, search string) ([]*mongo.Category, error)
}

type grpcQueryAdapter struct {
	client pbcategory.CategoryQueryServiceClient
	guard  *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *grpcQueryAdapter) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// New wraps a generated category query client into a QueryRepository. It is the
// canonical constructor; modules build it inside their repositories with the
// guard options they were given.
func New(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewQueryAdapter(client, opts...)
}

// NewQueryAdapter wraps a generated category query client into a QueryRepository.
// Resilience (timeout/circuit-breaker/bulkhead) is applied per call through the
// optional dependency guard; a nil guard is a passthrough.
func NewQueryAdapter(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	a := &grpcQueryAdapter{client: client}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *grpcQueryAdapter) FindByID(ctx context.Context, categoryID int) (*mongo.Category, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbcategory.ApiResponseCategory, error) {
		return a.client.FindById(ctx, &pbcategory.FindByIdCategoryRequest{Id: int32(categoryID)})
	})
	if err != nil {
		return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *grpcQueryAdapter) FindAll(ctx context.Context, page, pageSize int, search string) ([]*mongo.Category, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbcategory.ApiResponsePaginationCategory, error) {
		return a.client.FindAll(ctx, &pbcategory.FindAllCategoryRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		})
	})
	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModels(res.Data), nil
}

func toModel(c *pbcategory.CategoryResponse) *mongo.Category {
	if c == nil {
		return nil
	}
	return &mongo.Category{
		CategoryID:    c.Id,
		Name:          c.Name,
		Description:   &c.Description,
		SlugCategory:  &c.SlugCategory,
		ImageCategory: &c.ImageCategory,
		CreatedAt:     parseTimeString(c.CreatedAt),
		UpdatedAt:     parseTimeString(c.UpdatedAt),
	}
}

func toModels(cs []*pbcategory.CategoryResponse) []*mongo.Category {
	out := make([]*mongo.Category, 0, len(cs))
	for _, c := range cs {
		out = append(out, toModel(c))
	}
	return out
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
