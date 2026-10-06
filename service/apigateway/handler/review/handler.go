package reviewhandler

import (
	review_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/review"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/review"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsReview struct {
	Client        *grpc.ClientConn
	Router        chi.Router
	Logger        logger.LoggerInterface
	Cache         *cache.CacheStore
	Observability observability.TraceLoggerObservability
}

func RegisterReviewHandler(deps *DepsReview) {
	mapper := apimapper.NewReviewResponseMapper()
	cache := review_cache.NewReviewMencache(deps.Cache)

	NewReviewQueryHandleApi(&reviewQueryHandleDeps{
		client:        pb.NewReviewQueryServiceClient(deps.Client),
		router:        deps.Router,
		logger:        deps.Logger,
		mapper:        mapper.QueryMapper(),
		cache:         cache.QueryCache(),
		observability: deps.Observability,
	})

	NewReviewCommandHandleApi(&reviewCommandHandleDeps{
		client:        pb.NewReviewCommandServiceClient(deps.Client),
		router:        deps.Router,
		logger:        deps.Logger,
		mapper:        mapper.CommandMapper(),
		queryMapper:   mapper.QueryMapper(),
		cache:         cache.CommandCache(),
		observability: deps.Observability,
	})
}
