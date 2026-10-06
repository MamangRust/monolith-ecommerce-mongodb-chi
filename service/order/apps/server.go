package apps

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	pborders "github.com/MamangRust/monolith-ecommerce-pb/orders"
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbshipping_addresses "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// gRPC Client Connections
	userAddr := viper.GetString("GRPC_USER_ADDR")

	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	userQueryClient := pbusers.NewUserQueryServiceClient(userConn)

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")

	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pbproducts.NewProductQueryServiceClient(productConn)
	productCommandClient := pbproducts.NewProductCommandServiceClient(productConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")
	if merchantAddr == "" {
		merchantAddr = "50055"
	}
	merchantConn, err := grpc.NewClient("localhost:"+merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchants.NewMerchantQueryServiceClient(merchantConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "50056"
	}
	orderItemConn, err := grpc.NewClient("localhost:"+orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pborder_items.NewOrderItemQueryServiceClient(orderItemConn)
	orderItemCommandClient := pborder_items.NewOrderItemCommandServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "50063"
	}
	shippingConn, err := grpc.NewClient("localhost:"+shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingCommandClient := pbshipping_addresses.NewShippingCommandServiceClient(shippingConn)
	shippingQueryClient := pbshipping_addresses.NewShippingQueryServiceClient(shippingConn)

	transactionAddr := viper.GetString("GRPC_TRANSACTION_ADDR")
	if transactionAddr == "" {
		transactionAddr = "50061"
	}
	transactionConn, err := grpc.NewClient("localhost:"+transactionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to transaction service: %w", err)
	}
	transactionCommandClient := pbtransactions.NewTransactionCommandServiceClient(transactionConn)

	repos := repository.NewRepositories(&repository.Deps{
		DB:                 srv.MongoDB,
		UserQuery:          userQueryClient,
		ProductQuery:       productQueryClient,
		ProductCommand:     productCommandClient,
		MerchantQuery:      merchantQueryClient,
		OrderItemQuery:     orderItemQueryClient,
		OrderItemCommand:   orderItemCommandClient,
		ShippingCommand:    shippingCommandClient,
		ShippingQuery:      shippingQueryClient,
		TransactionCommand: transactionCommandClient,
		Guards: repository.GuardOptions{
			User: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Product: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Merchant: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			OrderItem: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("order-item", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Shipping: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("shipping-address", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Transaction: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("transaction", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
		},
	})

	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	obs, _ := observability.NewObservability("order-server", srv.Logger)
	mencache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         mencache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
		Kafka:         myKafka,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pborders.RegisterOrderQueryServiceServer(gs, h.OrderQuery)
		pborders.RegisterOrderStatsServiceServer(gs, h.OrderStats)
		pborders.RegisterOrderStatsByMerchantServiceServer(gs, h.OrderStatsByMerchant)
		pborders.RegisterOrderCommandServiceServer(gs, h.OrderCommand)
	}

	return srv, nil
}
