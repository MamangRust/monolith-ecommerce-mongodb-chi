package orderhandler

import (
	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/order"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/order"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsOrder struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterOrderHandler(deps *DepsOrder) {
	mapper := apimapper.NewOrderResponseMapper()
	cache := order_cache.OrderNewMencache(deps.CacheStore)

	queryClient := pb.NewOrderQueryServiceClient(deps.Client)
	statsClient := pb.NewOrderStatsServiceClient(deps.Client)

	NewOrderQueryHandleApi(&orderQueryHandleDeps{
		client: queryClient,
		router: deps.Router,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewOrderCommandHandleApi(&orderCommandHandleDeps{
		client: pb.NewOrderCommandServiceClient(deps.Client),
		router: deps.Router,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})

	NewOrderStatsHandleApi(&orderStatsHandleDeps{
		client:             statsClient,
		router:             deps.Router,
		logger:             deps.Logger,
		mapper:             mapper.StatsMapper(),
		cache:              cache,
		merchantStatsCache: cache,
	})
}
