package transaction_test

import (
	"context"
	"testing"
	"time"

	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionStatsServiceTestSuite struct {
	tests.BaseTestSuite
	svc           service.TransactionStatsService
	svcByMerchant service.TransactionStatsByMerchantService
	merchantID    int
	userID        int
}

func (s *TransactionStatsServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()
	s.SetupOrderService()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)
	repos := repository.NewRepositories(&repository.Deps{
		DB: mongoDB,
	})
	cacheStore := s.GetCacheStore()
	mencache := transaction_cache.NewMencache(cacheStore)

	deps := &service.TransactionStatsServiceDeps{
		Cache:         mencache.TransactionStatsCache,
		Repository:    repos.TransactionStats,
		Logger:        s.Log,
		Observability: s.Obs,
	}
	s.svc = service.NewTransactionStatsService(deps)

	depsMerchant := &service.TransactionStatsByMerchantServiceDeps{
		Cache:         mencache.TransactionStatsByMerchantCache,
		Repository:    repos.TransactionStatsByMerchant,
		Logger:        s.Log,
		Observability: s.Obs,
	}
	s.svcByMerchant = service.NewTransactionStatsByMerchantService(depsMerchant)

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction for this order
	_, err = mongoDB.Collection("transactions").InsertOne(ctx, bson.M{
		"transaction_id": int32(1),

		"order_id":       int32(orderID),
		"merchant_id":    int32(s.merchantID),
		"amount":         int32(10000),
		"payment_method": "credit_card",
		"payment_status": "success",
		"created_at":     time.Now(),
	})
	s.Require().NoError(err)

	// Create a failed transaction for this order
	_, err = mongoDB.Collection("transactions").InsertOne(ctx, bson.M{
		"transaction_id": int32(2),

		"order_id":       int32(orderID),
		"merchant_id":    int32(s.merchantID),
		"amount":         int32(10000),
		"payment_method": "credit_card",
		"payment_status": "failed",
		"created_at":     time.Now(),
	})
	s.Require().NoError(err)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyAmountSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthAmountTransaction{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyAmountSuccess() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyAmountSuccess(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyAmountFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthAmountTransaction{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyAmountFailed(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyAmountFailed() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyAmountFailed(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyMethodSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthMethodTransaction{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyMethodSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyMethodSuccess() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyMethodSuccess(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthAmountTransactionMerchant{
		Year:       now.Year(),
		Month:      int(now.Month()),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearAmountTransactionMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestTransactionStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsServiceTestSuite))
}
