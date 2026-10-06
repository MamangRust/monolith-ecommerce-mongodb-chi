package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*repository.CartResult, *int, error)
}

type CartCommandService interface {
	Create(ctx context.Context, req *requests.CreateCartRequest) (*repository.CartCreateResult, error)
	DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error)
	DeleteAll(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error)
}
