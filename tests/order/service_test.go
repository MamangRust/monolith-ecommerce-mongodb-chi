package order_test

import (
	"context"
	"testing"

	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-ecommerce-pb/order_items"
	pbproducts "github.com/MamangRust/monolith-ecommerce-pb/products"
	pbshipping_addresses "github.com/MamangRust/monolith-ecommerce-pb/shipping_addresses"
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *OrderServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	// Order dependencies
	mencache := order_cache.NewMencache(cacheStore)
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)
	repos := repository.NewRepositories(&repository.Deps{
		DB:                 mongoDB,
		MerchantQuery:      pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQuery:       pbproducts.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommand:     pbproducts.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQuery:     pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommand:   pborder_items.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
		UserQuery:          pbusers.NewUserQueryServiceClient(s.Conns["user"]),
		ShippingCommand:    pbshipping_addresses.NewShippingCommandServiceClient(s.Conns["shipping-address"]),
		TransactionCommand: pbtransactions.NewTransactionCommandServiceClient(s.Conns["transaction"]),
	})

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *OrderServiceTestSuite) TestOrderLifecycle() {
	ctx := context.Background()

	// 1. Setup Dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// 2. Create Order
	req := &requests.CreateOrderRequest{
		UserID:     userID,
		MerchantID: merchID,
		TotalPrice: 20000,
		Items: []requests.CreateOrderItemRequest{
			{
				ProductID: prodID,
				Quantity:  1,
				Price:     10000,
			},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat:         "Test Address",
			Provinsi:       "West Java",
			Kota:           "Bandung",
			Courier:        "JNE",
			ShippingMethod: "REG",
			ShippingCost:   10000,
			Negara:         "Indonesia",
		},
	}

	created, err := s.svc.OrderCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	orderID := int(created.OrderID)

	// 3. FindByID
	found, err := s.svc.OrderQuery.FindByID(ctx, orderID)
	s.Require().NoError(err)
	s.Equal(int32(req.TotalPrice), found.TotalPrice)

	// 4. FindAll
	_, total, err := s.svc.OrderQuery.FindAll(ctx, &requests.FindAllOrder{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.OrderCommand.Trash(ctx, orderID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.OrderQuery.FindTrashed(ctx, &requests.FindAllOrder{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.OrderQuery.FindActive(ctx, &requests.FindAllOrder{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, o := range active {
		s.NotEqual(orderID, int(o.OrderID))
	}

	// 8. Restore
	_, err = s.svc.OrderCommand.Restore(ctx, orderID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.OrderCommand.Trash(ctx, orderID)
	s.Require().NoError(err)
	success, err := s.svc.OrderCommand.DeletePermanent(ctx, orderID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	o1, _ := s.svc.OrderCommand.Create(ctx, req)
	o2, _ := s.svc.OrderCommand.Create(ctx, req)

	s.svc.OrderCommand.Trash(ctx, int(o1.OrderID))
	s.svc.OrderCommand.Trash(ctx, int(o2.OrderID))

	resRestoreAll, err := s.svc.OrderCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.OrderCommand.Trash(ctx, int(o1.OrderID))
	s.svc.OrderCommand.Trash(ctx, int(o2.OrderID))

	// Note: We use the correct method name from the transaction dependency
	// if applicable, but this is the Order service test.
	// Ensuring Cleanup.
	resDeleteAll, err := s.svc.OrderCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestOrderServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderServiceTestSuite))
}
