package merchantawardhandler

import (
	merchantaward_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_awards"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_award"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsMerchantAward struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterMerchantAwardHandler(deps *DepsMerchantAward) {
	mapper := apimapper.NewMerchantAwardResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantaward_cache.NewMerchantAward(deps.CacheStore)

	NewMerchantAwardQueryHandleApi(&merchantAwardQueryHandleDeps{
		client: pb.NewMerchantAwardQueryServiceClient(deps.Client),
		router: deps.Router,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantAwardCommandHandleApi(&merchantAwardCommandHandleDeps{
		client:         pb.NewMerchantAwardCommandServiceClient(deps.Client),
		router:         deps.Router,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
