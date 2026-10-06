package order_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	orderhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/order"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type OrderStatsApiTestSuite struct {
	tests.BaseTestSuite
	echo       chi.Router
	merchantID int
	userID     int
}

func (s *OrderStatsApiTestSuite) SetupSuite() {
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

	orderhandler.RegisterOrderHandler(&orderhandler.DepsOrder{
		Client:     s.Conns["order"],
		Router:     s.echo,
		Logger:     s.Log,
		CacheStore: s.GetCacheStore(),
	})

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	db, err := s.MongoDB()
	s.Require().NoError(err)
	_, err = db.Collection("orders").UpdateOne(ctx,
		bson.M{"order_id": int32(orderID)},
		bson.M{"$set": bson.M{"created_at": time.Now()}})
	s.Require().NoError(err)
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenue() {
	now := time.Now()
	url := fmt.Sprintf("/api/order/monthly-total-revenue?year=%d&month=%d", now.Year(), int(now.Month()))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *OrderStatsApiTestSuite) TestFindYearlyTotalRevenue() {
	year := time.Now().Year()
	url := fmt.Sprintf("/api/order/yearly-total-revenue?year=%d", year)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyOrder() {
	year := time.Now().Year()
	url := fmt.Sprintf("/api/order/monthly-revenue?year=%d", year)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	now := time.Now()
	url := fmt.Sprintf("/api/order/merchant/monthly-total-revenue?year=%d&month=%d&merchant_id=%d",
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

func TestOrderStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsApiTestSuite))
}
