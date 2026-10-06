package seeder

import (
	"context"
	"log"

	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Deps struct {
	Db     *mongo.Database
	Ctx    context.Context
	Logger logger.LoggerInterface
	Hash   hash.HashPassword
}

type Seeder struct {
	User             *userSeeder
	Role             *roleSeeder
	UserRole         *userRoleSeeder
	Merchant         *merchantSeeder
	MerchantDetail   *merchantDetailSeeder
	MerchantAward    *merchantAwardSeeder
	MerchantPolicy   *merchantPolicySeeder
	MerchantBusiness *merchantBusinessSeeder
	Category         *categorySeeder
	Product          *productSeeder
	Order            *orderSeeder
	Shipping         *shippingAddressSeeder
	Review           *reviewSeeder
	Transaction      *transactionSeeder
	Slider           *sliderSeeder
	Banner           *bannerSeeder
}

func NewSeeder(deps Deps) *Seeder {
	return &Seeder{
		User:             NewUserSeeder(deps.Db, deps.Hash, deps.Ctx, deps.Logger),
		Role:             NewRoleSeeder(deps.Db, deps.Ctx, deps.Logger),
		UserRole:         NewUserRoleSeeder(deps.Db, deps.Ctx, deps.Logger),
		Merchant:         NewMerchantSeeder(deps.Db, deps.Ctx, deps.Logger),
		MerchantDetail:   NewMerchantDetailSeeder(deps.Db, deps.Ctx, deps.Logger),
		MerchantAward:    NewMerchantAwardSeeder(deps.Db, deps.Ctx, deps.Logger),
		MerchantPolicy:   NewMerchantPolicySeeder(deps.Db, deps.Ctx, deps.Logger),
		MerchantBusiness: NewMerchantBusinessSeeder(deps.Db, deps.Ctx, deps.Logger),
		Category:         NewCategorySeeder(deps.Db, deps.Ctx, deps.Logger),
		Product:          NewProductSeeder(deps.Db, deps.Ctx, deps.Logger),
		Order:            NewOrderSeeder(deps.Db, deps.Ctx, deps.Logger),
		Shipping:         NewShippingAddressSeeder(deps.Db, deps.Ctx, deps.Logger),
		Review:           NewReviewSeeder(deps.Db, deps.Ctx, deps.Logger),
		Transaction:      NewTransactionSeeder(deps.Db, deps.Ctx, deps.Logger),
		Slider:           NewSliderSeeder(deps.Db, deps.Ctx, deps.Logger),
		Banner:           NewBannerSeeder(deps.Db, deps.Ctx, deps.Logger),
	}
}

func (s *Seeder) Run() error {
	if err := s.RunSession1(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(35 * time.Second)

	if err := s.RunSession2(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(35 * time.Second)

	if err := s.RunSession3(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Seeder) RunSession1() error {
	if err := s.seedWithDelay("users", s.User.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("roles", s.Role.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("user_roles", s.UserRole.Seed); err != nil {
		return err
	}

	return nil
}

func (s *Seeder) RunSession2() error {
	if err := s.seedWithDelay("merchants", s.Merchant.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchant_details", s.MerchantDetail.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchant_awards", s.MerchantAward.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchant_policies", s.MerchantPolicy.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchant_businesses", s.MerchantBusiness.Seed); err != nil {
		return err
	}

	return nil
}

func (s *Seeder) RunSession3() error {
	if err := s.seedWithDelay("categories", s.Category.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("products", s.Product.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("orders", s.Order.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("shipping_addresses", s.Shipping.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("transactions", s.Transaction.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("banners", s.Banner.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("reviews", s.Review.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("sliders", s.Slider.Seed); err != nil {
		return err
	}

	return nil
}

func (s *Seeder) seedWithDelay(entityName string, seedFunc func() error) error {
	if err := seedFunc(); err != nil {
		return fmt.Errorf("failed to seed %s: %w", entityName, err)
	}

	time.Sleep(20 * time.Second)
	return nil
}

func toBoolPtr(b bool) *bool {
	return &b
}

func toStringPtr(s string) *string {
	return &s
}

func toInt32Ptr(i int32) *int32 {
	return &i
}

func toFloat64Ptr(f float64) *float64 {
	return &f
}
