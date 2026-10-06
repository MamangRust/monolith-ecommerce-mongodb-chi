package bannerhandler

import (
	banner_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/banner"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/banner"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsBanner struct {
	Client        *grpc.ClientConn
	Router        chi.Router
	Logger        logger.LoggerInterface
	CacheStore    *cache.CacheStore
	Observability observability.TraceLoggerObservability
}

func RegisterBannerHandler(deps *DepsBanner) {
	mapper := apimapper.NewBannerResponseMapper()
	cache := banner_cache.NewBannerMencache(deps.CacheStore)

	NewBannerQueryHandleApi(&bannerQueryHandleDeps{
		client:        pb.NewBannerQueryServiceClient(deps.Client),
		router:        deps.Router,
		logger:        deps.Logger,
		mapper:        mapper.QueryMapper(),
		cache:         cache,
		observability: deps.Observability,
	})

	NewBannerCommandHandleApi(&bannerCommandHandleDeps{
		client:        pb.NewBannerCommandServiceClient(deps.Client),
		router:        deps.Router,
		logger:        deps.Logger,
		mapper:        mapper.CommandMapper(),
		cache:         cache,
		observability: deps.Observability,
	})
}
