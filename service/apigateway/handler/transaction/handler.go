package transactionhandler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/transaction"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsTransaction struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	ApiHandler apierror.ApiHandler
}

func RegisterTransactionHandler(deps *DepsTransaction) {
	mapper := apimapper.NewTransactionResponseMapper()
	statsMapper := apimapper.NewTransactionStatsResponseMapper()
	cache := transaction_cache.NewTransactionMencache(deps.CacheStore)

	queryClient := pb.NewTransactionQueryServiceClient(deps.Client)
	commandClient := pb.NewTransactionCommandServiceClient(deps.Client)
	statsClient := pb.NewTransactionStatsServiceClient(deps.Client)
	statsByMerchantClient := pb.NewTransactionStatsByMerchantServiceClient(deps.Client)

	NewTransactionQueryHandleApi(&transactionQueryHandleDeps{
		queryClient: queryClient,
		router:      deps.Router,
		logger:      deps.Logger,
		mapper:      mapper.QueryMapper(),
		cache:       cache,
		apiHandler:  deps.ApiHandler,
	})

	NewTransactionCommandHandleApi(&transactionCommandHandleDeps{
		client:     commandClient,
		router:     deps.Router,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewTransactionStatsHandleApi(&transactionStatsHandleDeps{
		statsClient:           statsClient,
		statsByMerchantClient: statsByMerchantClient,
		router:                deps.Router,
		logger:                deps.Logger,
		statsMapper:           statsMapper,
		statsCache:            cache,
		statsByMerchantCache:  cache,
		apiHandler:            deps.ApiHandler,
	})
}
