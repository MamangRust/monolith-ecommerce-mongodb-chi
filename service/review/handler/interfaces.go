package handler

import (
	pbreviews "github.com/MamangRust/monolith-ecommerce-pb/reviews"
)

type ReviewHandleGrpc interface {
	pbreviews.ReviewQueryServiceServer
	pbreviews.ReviewCommandServiceServer
}

type ReviewQueryHandler interface {
	pbreviews.ReviewQueryServiceServer
}

type ReviewCommandHandler interface {
	pbreviews.ReviewCommandServiceServer
}
