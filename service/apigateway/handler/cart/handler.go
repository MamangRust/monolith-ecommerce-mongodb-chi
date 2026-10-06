package carthandler

import (
	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/cart"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/cart"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsCart struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterCartHandler(deps *DepsCart) {
	mapper := apimapper.NewCartResponseMapper()
	cache := cart_cache.NewCartMencache(deps.CacheStore)

	NewCartQueryHandleApi(&cartQueryHandleDeps{
		client: pb.NewCartQueryServiceClient(deps.Client),
		router: deps.Router,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewCartCommandHandleApi(&cartCommandHandleDeps{
		client: pb.NewCartCommandServiceClient(deps.Client),
		router: deps.Router,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})
}
