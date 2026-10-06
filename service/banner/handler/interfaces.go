package handler

import (
	pbbanners "github.com/MamangRust/monolith-ecommerce-pb/banners"
)

type BannerQueryHandler interface {
	pbbanners.BannerQueryServiceServer
}

type BannerCommandHandler interface {
	pbbanners.BannerCommandServiceServer
}
