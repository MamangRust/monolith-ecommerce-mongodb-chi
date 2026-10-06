package handler

import (
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
)

type ProductQueryHandler interface {
	pbproducts.ProductQueryServiceServer
}

type ProductCommandHandler interface {
	pbproducts.ProductCommandServiceServer
}
