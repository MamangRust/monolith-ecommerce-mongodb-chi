package handler

import (
	pbcarts "github.com/MamangRust/monolith-ecommerce-pb/carts"
)

type CartQueryHandler interface {
	pbcarts.CartQueryServiceServer
}

type CartCommandHandler interface {
	pbcarts.CartCommandServiceServer
}
