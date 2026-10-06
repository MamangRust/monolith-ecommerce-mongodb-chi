package category_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"time"

	pbcategories "github.com/MamangRust/monolith-ecommerce-pb/categories"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryStatsGapiTestSuite struct {
	tests.BaseTestSuite
	client           pbcategories.CategoryStatsServiceClient
	clientById       pbcategories.CategoryStatsByIdServiceClient
	clientByMerchant pbcategories.CategoryStatsByMerchantServiceClient
	categoryID       int
	merchantID       int
	userID           int
}

func (s *CategoryStatsGapiTestSuite) SetupSuite() {
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

	s.client = pbcategories.NewCategoryStatsServiceClient(s.Conns["category"])
	s.clientById = pbcategories.NewCategoryStatsByIdServiceClient(s.Conns["category"])
	s.clientByMerchant = pbcategories.NewCategoryStatsByMerchantServiceClient(s.Conns["category"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.categoryID = s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, s.categoryID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	db, err := s.MongoDB()
	s.Require().NoError(err)
	_, err = db.Collection("orders").UpdateOne(ctx, bson.M{"order_id": int32(orderID)}, bson.M{"$set": bson.M{"created_at": time.Now()}})
	s.Require().NoError(err)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPrices() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearMonthTotalPrices{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.FindMonthlyTotalPrices(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPrices() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategories.FindYearTotalPrices{
		Year: int32(year),
	}

	res, err := s.client.FindYearlyTotalPrices(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPrice() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategories.FindYearCategory{
		Year: int32(year),
	}

	res, err := s.client.FindMonthPrice(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPrice() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategories.FindYearCategory{
		Year: int32(year),
	}

	res, err := s.client.FindYearPrice(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPricesById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearMonthTotalPriceById{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindMonthlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPricesById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearTotalPriceById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindYearlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindMonthPriceById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindYearPriceById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPricesByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearMonthTotalPriceByMerchant{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindMonthlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPricesByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearTotalPriceByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindYearlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindMonthPriceByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindYearPriceByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestCategoryStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsGapiTestSuite))
}
