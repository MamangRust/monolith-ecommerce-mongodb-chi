package merchantdetailhandler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	merchant_detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_detail"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_detail"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsMerchantDetail struct {
	Client      *grpc.ClientConn
	Router      chi.Router
	Logger      logger.LoggerInterface
	CacheStore  *cache.CacheStore
	UploadImage upload_image.ImageUploads
	ApiHandler  apierror.ApiHandler
}

func RegisterMerchantDetailHandler(deps *DepsMerchantDetail) {
	mapper := apimapper.NewMerchantDetailResponseMapper()
	cache := merchant_detail_cache.NewMerchantDetailMencache(deps.CacheStore)

	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	handlers := []func(){
		setupMerchantDetailQueryHandler(deps, mapper.QueryMapper(), cache),
		setupMerchantDetailCommandHandler(deps, mapper.CommandMapper(), merchantMapper.CommandMapper(), cache),
	}

	for _, h := range handlers {
		h()
	}
}

func setupMerchantDetailQueryHandler(deps *DepsMerchantDetail, mapper apimapper.MerchantDetailQueryResponseMapper, cache merchant_detail_cache.MerchantDetailQueryCache) func() {
	return func() {
		NewMerchantDetailQueryHandleApi(&merchantDetailQueryHandleDeps{
			client:     pb.NewMerchantDetailQueryServiceClient(deps.Client),
			router:     deps.Router,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}

func setupMerchantDetailCommandHandler(deps *DepsMerchantDetail, mapper apimapper.MerchantDetailCommandResponseMapper, merchantMapper merchantapimapper.MerchantCommandResponseMapper, cache merchant_detail_cache.MerchantDetailCommandCache) func() {
	return func() {
		NewMerchantDetailCommandHandleApi(&merchantDetailCommandHandleDeps{
			client:         pb.NewMerchantDetailCommandServiceClient(deps.Client),
			router:         deps.Router,
			logger:         deps.Logger,
			mapper:         mapper,
			merchantMapper: merchantMapper,
			cache:          cache,
			upload_image:   deps.UploadImage,
			apiHandler:     deps.ApiHandler,
		})
	}
}
