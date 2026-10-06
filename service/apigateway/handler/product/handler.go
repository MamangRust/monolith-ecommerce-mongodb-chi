package producthandler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	product_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/product"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/product"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsProduct struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	Upload     upload_image.ImageUploads
	ApiHandler apierror.ApiHandler
}

func RegisterProductHandler(deps *DepsProduct) {
	mapper := apimapper.NewProductResponseMapper()
	cache := product_cache.NewProductMencache(deps.CacheStore)

	queryClient := pb.NewProductQueryServiceClient(deps.Client)

	NewProductQueryHandleApi(&productQueryHandleDeps{
		client:     queryClient,
		router:     deps.Router,
		logger:     deps.Logger,
		mapper:     mapper.QueryMapper(),
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewProductCommandHandleApi(&productCommandHandleDeps{
		client:       pb.NewProductCommandServiceClient(deps.Client),
		router:       deps.Router,
		logger:       deps.Logger,
		mapper:       mapper.CommandMapper(),
		cache:        cache,
		upload_image: deps.Upload,
		apiHandler:   deps.ApiHandler,
	})

}
