package transaction_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	transactionhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/transaction"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionStatsApiTestSuite struct {
	tests.BaseTestSuite
	echo       chi.Router
	merchantID int
	userID     int
}

func (s *TransactionStatsApiTestSuite) SetupSuite() {
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

	s.echo = chi.NewRouter()
	apiHandler := apierror.NewApiHandler(s.Obs, s.Log)

	transactionhandler.RegisterTransactionHandler(&transactionhandler.DepsTransaction{
		Client:     s.Conns["transaction"],
		Router:     s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
		ApiHandler: apiHandler,
	})

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Seed a successful transaction
	db, err := s.MongoDB()
	s.Require().NoError(err)
	_, err = db.Collection("transactions").InsertOne(ctx, bson.M{
		"merchant_id":    int32(s.merchantID),
		"transaction_id": int32(1),

		"order_id":       int32(orderID),
		"amount":         int32(100000),
		"payment_method": "credit_card",
		"payment_status": "success",
		"created_at":     time.Now(),
	})
	s.Require().NoError(err)
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccess() {
	now := time.Now()
	url := fmt.Sprintf("/api/transaction-stats/monthly-success?year=%d&month=%d", now.Year(), int(now.Month()))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindYearStatusSuccess() {
	year := time.Now().Year()
	url := fmt.Sprintf("/api/transaction-stats/yearly-success?year=%d", year)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindMonthMethodSuccess() {
	now := time.Now()
	url := fmt.Sprintf("/api/transaction-stats/monthly-method-success?year=%d&month=%d", now.Year(), int(now.Month()))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccessByMerchant() {
	now := time.Now()
	url := fmt.Sprintf("/api/transaction-stats/merchant/monthly-success?year=%d&month=%d&merchant_id=%d",
		now.Year(), int(now.Month()), s.merchantID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func TestTransactionStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsApiTestSuite))
}
