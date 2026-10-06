package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressResult struct {
	ShippingAddressID int32
	OrderID           int32
	Alamat            string
	Provinsi          string
	Kota              string
	Negara            string
	Courier           string
	ShippingMethod    string
	ShippingCost      float64
	CreatedAt         *string
	UpdatedAt         *string
	DeletedAt         *string
	TotalCount        int64
}

type ShippingAddressQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error)
	FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error)
	FindByOrder(ctx context.Context, shipping_id int) (*ShippingAddressResult, error)
	FindByID(ctx context.Context, shipping_id int) (*ShippingAddressResult, error)
}

type ShippingAddressCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*salemongo.ShippingAddress, error)
	Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*salemongo.ShippingAddress, error)
	Trash(ctx context.Context, shipping_id int) (*salemongo.ShippingAddress, error)
	Restore(ctx context.Context, shipping_id int) (*salemongo.ShippingAddress, error)
	DeletePermanent(ctx context.Context, shipping_id int) (bool, error)
	DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
