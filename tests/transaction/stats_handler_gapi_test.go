package transaction_test

import (
	"context"
	"testing"
	"time"

	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionStatsGapiTestSuite struct {
	tests.BaseTestSuite
	client           pbtransactions.TransactionStatsServiceClient
	clientByMerchant pbtransactions.TransactionStatsByMerchantServiceClient
	merchantID       int
	userID           int
}

func (s *TransactionStatsGapiTestSuite) SetupSuite() {
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

	s.client = pbtransactions.NewTransactionStatsServiceClient(s.Conns["transaction"])
	s.clientByMerchant = pbtransactions.NewTransactionStatsByMerchantServiceClient(s.Conns["transaction"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction for this order
	db, err := s.MongoDB()
	s.Require().NoError(err)
	_, err = db.Collection("transactions").InsertOne(ctx, bson.M{
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
	_, err = db.Collection("transactions").InsertOne(ctx, bson.M{
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

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyAmountSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.MonthAmountTransactionRequest{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.GetMonthlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyAmountSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.YearAmountTransactionRequest{
		Year: int32(year),
	}

	res, err := s.client.GetYearlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyAmountFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.MonthAmountTransactionRequest{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.GetMonthlyAmountFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyAmountFailed() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.YearAmountTransactionRequest{
		Year: int32(year),
	}

	res, err := s.client.GetYearlyAmountFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyTransactionMethodSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.MonthMethodTransactionRequest{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.GetMonthlyTransactionMethodSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyTransactionMethodSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.YearMethodTransactionRequest{
		Year: int32(year),
	}

	res, err := s.client.GetYearlyTransactionMethodSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.MonthAmountTransactionMerchantRequest{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.GetMonthlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.YearAmountTransactionMerchantRequest{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.GetYearlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestTransactionStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsGapiTestSuite))
}
