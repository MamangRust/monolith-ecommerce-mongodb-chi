package mongo

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsDuplicateKeyError checks if the error is a MongoDB duplicate key error.
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		for _, we := range writeErr.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	// findAndModify surfaces duplicates as a command error carrying code 11000.
	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) && cmdErr.Code == 11000 {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Core entities (bson-tagged, collection-backed)
// ---------------------------------------------------------------------------

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	UserID           int32              `bson:"user_id"`
	Firstname        string             `bson:"firstname"`
	Lastname         string             `bson:"lastname"`
	Email            string             `bson:"email"`
	Password         string             `bson:"password"`
	VerificationCode string             `bson:"verification_code"`
	IsVerified       *bool              `bson:"is_verified"`
	CreatedAt        *time.Time         `bson:"created_at"`
	UpdatedAt        *time.Time         `bson:"updated_at"`
	DeletedAt        *time.Time         `bson:"deleted_at,omitempty"`
}

func (User) CollectionName() string { return "users" }

type UserRole struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserRoleID int32              `bson:"user_role_id"`
	UserID     int32              `bson:"user_id"`
	RoleID     int32              `bson:"role_id"`
	CreatedAt  *time.Time         `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

func (UserRole) CollectionName() string { return "user_roles" }

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoleID    int32              `bson:"role_id"`
	RoleName  string             `bson:"role_name"`
	CreatedAt *time.Time         `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (Role) CollectionName() string { return "roles" }

type RefreshToken struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	RefreshTokenID int32              `bson:"refresh_token_id"`
	UserID         int32              `bson:"user_id"`
	Token          string             `bson:"token"`
	Expiration     time.Time          `bson:"expiration"`
	CreatedAt      *time.Time         `bson:"created_at"`
	UpdatedAt      *time.Time         `bson:"updated_at"`
	DeletedAt      *time.Time         `bson:"deleted_at,omitempty"`
}

func (RefreshToken) CollectionName() string { return "refresh_tokens" }

type ResetToken struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ResetTokenID int32              `bson:"reset_token_id"`
	UserID       int64              `bson:"user_id"`
	Token        string             `bson:"token"`
	ExpiryDate   time.Time          `bson:"expiry_date"`
}

func (ResetToken) CollectionName() string { return "reset_tokens" }

// ---------------------------------------------------------------------------
// Banner
// ---------------------------------------------------------------------------

type Banner struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BannerID  int32              `bson:"banner_id"`
	Name      string             `bson:"name"`
	StartDate time.Time          `bson:"start_date"`
	EndDate   time.Time          `bson:"end_date"`
	StartTime string             `bson:"start_time"`
	EndTime   string             `bson:"end_time"`
	IsActive  *bool              `bson:"is_active"`
	CreatedAt *time.Time         `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (Banner) CollectionName() string { return "banners" }

// ---------------------------------------------------------------------------
// Cart
// ---------------------------------------------------------------------------

type Cart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CartID    int32              `bson:"cart_id"`
	UserID    int32              `bson:"user_id"`
	ProductID int32              `bson:"product_id"`
	Name      string             `bson:"name"`
	Price     int32              `bson:"price"`
	Image     string             `bson:"image"`
	Quantity  int32              `bson:"quantity"`
	Weight    int32              `bson:"weight"`
	CreatedAt *time.Time         `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (Cart) CollectionName() string { return "carts" }

// ---------------------------------------------------------------------------
// Category
// ---------------------------------------------------------------------------

type Category struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID    int32              `bson:"category_id"`
	Name          string             `bson:"name"`
	Description   *string            `bson:"description"`
	SlugCategory  *string            `bson:"slug_category"`
	ImageCategory *string            `bson:"image_category"`
	CreatedAt     *time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at"`
	DeletedAt     *time.Time         `bson:"deleted_at,omitempty"`
}

func (Category) CollectionName() string { return "categories" }

// ---------------------------------------------------------------------------
// Merchant
// ---------------------------------------------------------------------------

type Merchant struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	MerchantID   int32              `bson:"merchant_id"`
	UserID       int32              `bson:"user_id"`
	Name         string             `bson:"name"`
	Description  *string            `bson:"description"`
	Address      *string            `bson:"address"`
	ContactEmail *string            `bson:"contact_email"`
	ContactPhone *string            `bson:"contact_phone"`
	Status       string             `bson:"status"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (Merchant) CollectionName() string { return "merchants" }

type MerchantDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DocumentID   int32              `bson:"document_id"`
	MerchantID   int32              `bson:"merchant_id"`
	DocumentType string             `bson:"document_type"`
	DocumentUrl  string             `bson:"document_url"`
	Status       string             `bson:"status"`
	Note         *string            `bson:"note"`
	UploadedAt   *time.Time         `bson:"uploaded_at"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantDocument) CollectionName() string { return "merchant_documents" }

type MerchantBusinessInformation struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	MerchantBusinessInfoID int32              `bson:"merchant_business_info_id"`
	MerchantID             int32              `bson:"merchant_id"`
	BusinessType           *string            `bson:"business_type"`
	TaxID                  *string            `bson:"tax_id"`
	EstablishedYear        *int32             `bson:"established_year"`
	NumberOfEmployees      *int32             `bson:"number_of_employees"`
	WebsiteUrl             *string            `bson:"website_url"`
	CreatedAt              *time.Time         `bson:"created_at"`
	UpdatedAt              *time.Time         `bson:"updated_at"`
	DeletedAt              *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantBusinessInformation) CollectionName() string { return "merchant_business_information" }

type MerchantCertificationsAndAward struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty"`
	MerchantCertificationID int32              `bson:"merchant_certification_id"`
	MerchantID              int32              `bson:"merchant_id"`
	Title                   string             `bson:"title"`
	Description             *string            `bson:"description"`
	IssuedBy                *string            `bson:"issued_by"`
	IssueDate               time.Time          `bson:"issue_date"`
	ExpiryDate              time.Time          `bson:"expiry_date"`
	CertificateUrl          *string            `bson:"certificate_url"`
	CreatedAt               *time.Time         `bson:"created_at"`
	UpdatedAt               *time.Time         `bson:"updated_at"`
	DeletedAt               *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantCertificationsAndAward) CollectionName() string {
	return "merchant_certifications_and_awards"
}

type MerchantDetail struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	MerchantDetailID int32              `bson:"merchant_detail_id"`
	MerchantID       int32              `bson:"merchant_id"`
	DisplayName      *string            `bson:"display_name"`
	CoverImageUrl    *string            `bson:"cover_image_url"`
	LogoUrl          *string            `bson:"logo_url"`
	ShortDescription *string            `bson:"short_description"`
	WebsiteUrl       *string            `bson:"website_url"`
	CreatedAt        *time.Time         `bson:"created_at"`
	UpdatedAt        *time.Time         `bson:"updated_at"`
	DeletedAt        *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantDetail) CollectionName() string { return "merchant_details" }

type MerchantSocialMediaLink struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	MerchantSocialID int32              `bson:"merchant_social_id"`
	MerchantDetailID int32              `bson:"merchant_detail_id"`
	Platform         string             `bson:"platform"`
	Url              string             `bson:"url"`
	CreatedAt        *time.Time         `bson:"created_at"`
	UpdatedAt        *time.Time         `bson:"updated_at"`
	DeletedAt        *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantSocialMediaLink) CollectionName() string { return "merchant_social_media_links" }

type MerchantPolicy struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	MerchantPolicyID int32              `bson:"merchant_policy_id"`
	MerchantID       int32              `bson:"merchant_id"`
	PolicyType       string             `bson:"policy_type"`
	Title            string             `bson:"title"`
	Description      string             `bson:"description"`
	CreatedAt        *time.Time         `bson:"created_at"`
	UpdatedAt        *time.Time         `bson:"updated_at"`
	DeletedAt        *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantPolicy) CollectionName() string { return "merchant_policies" }

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OrderID    int32              `bson:"order_id"`
	UserID     int32              `bson:"user_id"`
	MerchantID int32              `bson:"merchant_id"`
	TotalPrice int32              `bson:"total_price"`
	CreatedAt  *time.Time         `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

func (Order) CollectionName() string { return "orders" }

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OrderItemID int32              `bson:"order_item_id"`
	OrderID     int32              `bson:"order_id"`
	ProductID   int32              `bson:"product_id"`
	Quantity    int32              `bson:"quantity"`
	Price       int32              `bson:"price"`
	CreatedAt   *time.Time         `bson:"created_at"`
	UpdatedAt   *time.Time         `bson:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty"`
}

func (OrderItem) CollectionName() string { return "order_items" }

type OrderStockReservation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ReservationID int32              `bson:"reservation_id"`
	OrderID       int32              `bson:"order_id"`
	ProductID     int32              `bson:"product_id"`
	Quantity      int32              `bson:"quantity"`
	Status        string             `bson:"status"`
	CreatedAt     *time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at"`
}

func (OrderStockReservation) CollectionName() string { return "order_stock_reservations" }

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProductID    int32              `bson:"product_id"`
	MerchantID   int32              `bson:"merchant_id"`
	CategoryID   int32              `bson:"category_id"`
	Name         string             `bson:"name"`
	Description  *string            `bson:"description"`
	Price        int32              `bson:"price"`
	CountInStock int32              `bson:"count_in_stock"`
	Brand        *string            `bson:"brand"`
	Weight       *int32             `bson:"weight"`
	Rating       *float64           `bson:"rating"`
	SlugProduct  *string            `bson:"slug_product"`
	ImageProduct *string            `bson:"image_product"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (Product) CollectionName() string { return "products" }

type ProductStockAdjustment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OperationID string             `bson:"operation_id"`
	ProductID   int32              `bson:"product_id"`
	Delta       int32              `bson:"delta"`
	CreatedAt   *time.Time         `bson:"created_at"`
}

func (ProductStockAdjustment) CollectionName() string { return "product_stock_adjustments" }

// ---------------------------------------------------------------------------
// Review
// ---------------------------------------------------------------------------

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ReviewID  int32              `bson:"review_id"`
	UserID    int32              `bson:"user_id"`
	ProductID int32              `bson:"product_id"`
	Name      string             `bson:"name"`
	Comment   string             `bson:"comment"`
	Rating    int32              `bson:"rating"`
	CreatedAt *time.Time         `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (Review) CollectionName() string { return "reviews" }

type ReviewDetail struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ReviewDetailID int32              `bson:"review_detail_id"`
	ReviewID       int32              `bson:"review_id"`
	Type           string             `bson:"type"`
	Url            string             `bson:"url"`
	Caption        *string            `bson:"caption"`
	CreatedAt      *time.Time         `bson:"created_at"`
	UpdatedAt      *time.Time         `bson:"updated_at"`
	DeletedAt      *time.Time         `bson:"deleted_at,omitempty"`
}

func (ReviewDetail) CollectionName() string { return "review_details" }

// ---------------------------------------------------------------------------
// ShippingAddress
// ---------------------------------------------------------------------------

type ShippingAddress struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	ShippingAddressID int32              `bson:"shipping_address_id"`
	OrderID           int32              `bson:"order_id"`
	Alamat            string             `bson:"alamat"`
	Provinsi          string             `bson:"provinsi"`
	Negara            string             `bson:"negara"`
	Kota              string             `bson:"kota"`
	Courier           string             `bson:"courier"`
	ShippingMethod    string             `bson:"shipping_method"`
	ShippingCost      float64            `bson:"shipping_cost"`
	CreatedAt         *time.Time         `bson:"created_at"`
	UpdatedAt         *time.Time         `bson:"updated_at"`
	DeletedAt         *time.Time         `bson:"deleted_at,omitempty"`
}

func (ShippingAddress) CollectionName() string { return "shipping_addresses" }

// ---------------------------------------------------------------------------
// Slider
// ---------------------------------------------------------------------------

type Slider struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SliderID  int32              `bson:"slider_id"`
	Name      string             `bson:"name"`
	Image     string             `bson:"image"`
	CreatedAt *time.Time         `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (Slider) CollectionName() string { return "sliders" }

// ---------------------------------------------------------------------------
// Transaction
// ---------------------------------------------------------------------------

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	TransactionID int32              `bson:"transaction_id"`
	OrderID       int32              `bson:"order_id"`
	MerchantID    int32              `bson:"merchant_id"`
	PaymentMethod string             `bson:"payment_method"`
	Amount        int32              `bson:"amount"`
	PaymentStatus string             `bson:"payment_status"`
	CreatedAt     *time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at"`
	DeletedAt     *time.Time         `bson:"deleted_at,omitempty"`
}

func (Transaction) CollectionName() string { return "transactions" }

// ---------------------------------------------------------------------------
// ConsumerInbox
// ---------------------------------------------------------------------------

type ConsumerInbox struct {
	ConsumerName       string     `bson:"consumer_name"`
	EventKey           string     `bson:"event_key"`
	Topic              string     `bson:"topic"`
	PartitionID        int32      `bson:"partition_id"`
	MessageOffset      int64      `bson:"message_offset"`
	Status             string     `bson:"status"`
	Attempts           int32      `bson:"attempts"`
	LeaseUntil         time.Time  `bson:"lease_until"`
	LastError          string     `bson:"last_error"`
	ProcessedAt        *time.Time `bson:"processed_at,omitempty"`
	ReservationVersion int64      `bson:"reservation_version"`
}

func (ConsumerInbox) CollectionName() string { return "consumer_inbox" }

// ---------------------------------------------------------------------------
// Row types (sqlc-compatible result rows)
// ---------------------------------------------------------------------------

type RoleRow struct {
	RoleID     int32      `bson:"role_id"`
	RoleName   string     `bson:"role_name"`
	CreatedAt  time.Time  `bson:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
	TotalCount int64      `bson:"total_count"`
}

type RoleActiveRow = RoleRow
type RoleTrashedRow = RoleRow

type MerchantRow struct {
	MerchantID   int32     `bson:"merchant_id"`
	UserID       int32     `bson:"user_id"`
	Name         string    `bson:"name"`
	Description  *string   `bson:"description"`
	Address      *string   `bson:"address"`
	ContactEmail *string   `bson:"contact_email"`
	ContactPhone *string   `bson:"contact_phone"`
	Status       string    `bson:"status"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

type MerchantActiveRow struct {
	MerchantRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type MerchantTrashedRow = MerchantActiveRow

type MerchantWithDeletedAt struct {
	Merchant
}

type MerchantDocumentRow struct {
	DocumentID   int32     `bson:"document_id"`
	MerchantID   int32     `bson:"merchant_id"`
	DocumentType string    `bson:"document_type"`
	DocumentUrl  string    `bson:"document_url"`
	Status       string    `bson:"status"`
	Note         *string   `bson:"note"`
	UploadedAt   time.Time `bson:"uploaded_at"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

type MerchantDocumentsRow struct {
	MerchantDocumentRow
	TotalCount int64 `bson:"total_count"`
}

type MerchantDocumentsActiveRow struct {
	MerchantDocumentRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type MerchantDocumentsTrashedRow = MerchantDocumentsActiveRow

type BannerRow struct {
	BannerID  int32     `bson:"banner_id"`
	Name      string    `bson:"name"`
	StartDate time.Time `bson:"start_date"`
	EndDate   time.Time `bson:"end_date"`
	StartTime string    `bson:"start_time"`
	EndTime   string    `bson:"end_time"`
	IsActive  *bool     `bson:"is_active"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type BannerListRow struct {
	BannerRow
	TotalCount int64 `bson:"total_count"`
}

type BannerActiveRow struct {
	BannerRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type BannerTrashedRow = BannerActiveRow

type SliderRow struct {
	SliderID  int32     `bson:"slider_id"`
	Name      string    `bson:"name"`
	Image     string    `bson:"image"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type SliderListRow struct {
	SliderRow
	TotalCount int64 `bson:"total_count"`
}

type SliderActiveRow struct {
	SliderRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type SliderTrashedRow = SliderActiveRow

type CategoryRow struct {
	CategoryID    int32     `bson:"category_id"`
	Name          string    `bson:"name"`
	Description   *string   `bson:"description"`
	SlugCategory  *string   `bson:"slug_category"`
	ImageCategory *string   `bson:"image_category"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

type CategoryListRow struct {
	CategoryRow
	TotalCount int64 `bson:"total_count"`
}

type CategoryActiveRow struct {
	CategoryRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type CategoryTrashedRow = CategoryActiveRow

type CategoryStatsRow struct {
	Month              string  `bson:"month"`
	Year               string  `bson:"year"`
	CategoryID         int32   `bson:"category_id"`
	CategoryName       string  `bson:"category_name"`
	OrderCount         int64   `bson:"order_count"`
	ItemsSold          int64   `bson:"items_sold"`
	TotalRevenue       float64 `bson:"total_revenue"`
	UniqueProductsSold int64   `bson:"unique_products_sold"`
}

type CategoryRevenueRow struct {
	Year         string  `bson:"year"`
	Month        string  `bson:"month"`
	TotalRevenue float64 `bson:"total_revenue"`
}

type CategoryYearlyRevenueRow struct {
	Year         string  `bson:"year"`
	TotalRevenue float64 `bson:"total_revenue"`
}

type CategoryYearlyStatsRow struct {
	Year               string  `bson:"year"`
	CategoryID         int32   `bson:"category_id"`
	CategoryName       string  `bson:"category_name"`
	OrderCount         int64   `bson:"order_count"`
	ItemsSold          int64   `bson:"items_sold"`
	TotalRevenue       float64 `bson:"total_revenue"`
	UniqueProductsSold int64   `bson:"unique_products_sold"`
}

type ProductRow struct {
	ProductID    int32     `bson:"product_id"`
	MerchantID   int32     `bson:"merchant_id"`
	CategoryID   int32     `bson:"category_id"`
	Name         string    `bson:"name"`
	Description  *string   `bson:"description"`
	Price        int32     `bson:"price"`
	CountInStock int32     `bson:"count_in_stock"`
	Brand        *string   `bson:"brand"`
	Weight       *int32    `bson:"weight"`
	Rating       *float64  `bson:"rating"`
	SlugProduct  *string   `bson:"slug_product"`
	ImageProduct *string   `bson:"image_product"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

type ProductListRow struct {
	ProductRow
	TotalCount int64 `bson:"total_count"`
}

type ProductActiveRow struct {
	ProductRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type ProductTrashedRow = ProductActiveRow

type ProductByCategoryRow struct {
	ProductRow
	TotalCount   int64  `bson:"total_count"`
	CategoryName string `bson:"category_name"`
}

type ProductByMerchantRow = ProductByCategoryRow

type ProductStockRow struct {
	ProductID    int32 `bson:"product_id"`
	CountInStock int32 `bson:"count_in_stock"`
}

type CartRow struct {
	CartID    int32     `bson:"cart_id"`
	UserID    int32     `bson:"user_id"`
	ProductID int32     `bson:"product_id"`
	Name      string    `bson:"name"`
	Price     int32     `bson:"price"`
	Image     string    `bson:"image"`
	Quantity  int32     `bson:"quantity"`
	Weight    int32     `bson:"weight"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CartListRow struct {
	CartRow
	TotalCount int64 `bson:"total_count"`
}

type ReviewRow struct {
	ReviewID  int32     `bson:"review_id"`
	UserID    int32     `bson:"user_id"`
	ProductID int32     `bson:"product_id"`
	Name      string    `bson:"name"`
	Comment   string    `bson:"comment"`
	Rating    int32     `bson:"rating"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type ReviewListRow struct {
	ReviewRow
	TotalCount int64 `bson:"total_count"`
}

type ReviewActiveRow struct {
	ReviewRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type ReviewTrashedRow = ReviewActiveRow

type ReviewByProductOrMerchantRow struct {
	ReviewRow
	TotalCount    int64      `bson:"total_count"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty"`
	ReviewDetails string     `bson:"review_details"`
}

type ReviewDetailRow struct {
	ReviewDetailID int32     `bson:"review_detail_id"`
	ReviewID       int32     `bson:"review_id"`
	Type           string    `bson:"type"`
	Url            string    `bson:"url"`
	Caption        *string   `bson:"caption"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}

type ReviewDetailListRow struct {
	ReviewDetailRow
	TotalCount int64 `bson:"total_count"`
}

type ReviewDetailActiveRow struct {
	ReviewDetailRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type ReviewDetailTrashedRow = ReviewDetailActiveRow

type OrderRow struct {
	OrderID    int32     `bson:"order_id"`
	UserID     int32     `bson:"user_id"`
	MerchantID int32     `bson:"merchant_id"`
	TotalPrice int32     `bson:"total_price"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

type OrderListRow struct {
	OrderRow
	TotalCount int64 `bson:"total_count"`
}

type OrderActiveRow struct {
	OrderRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type OrderTrashedRow = OrderActiveRow

type OrderByMerchantRow = OrderActiveRow

type OrderStatsRow struct {
	Year               string  `bson:"year"`
	Month              string  `bson:"month"`
	OrderCount         int64   `bson:"order_count"`
	TotalRevenue       float64 `bson:"total_revenue"`
	TotalItemsSold     int64   `bson:"total_items_sold"`
	UniqueProductsSold int64   `bson:"unique_products_sold"`
}

type OrderRevenueRow struct {
	Year         string  `bson:"year"`
	Month        string  `bson:"month"`
	TotalRevenue float64 `bson:"total_revenue"`
}

type OrderMonthlyOrderRow struct {
	Month          string  `bson:"month"`
	OrderCount     int64   `bson:"order_count"`
	TotalRevenue   float64 `bson:"total_revenue"`
	TotalItemsSold int64   `bson:"total_items_sold"`
}

type OrderYearlyOrderRow struct {
	Year               string  `bson:"year"`
	OrderCount         int64   `bson:"order_count"`
	TotalRevenue       float64 `bson:"total_revenue"`
	TotalItemsSold     int64   `bson:"total_items_sold"`
	UniqueProductsSold int64   `bson:"unique_products_sold"`
}

type OrderMonthlyRevenueRow struct {
	Year         string `bson:"year"`
	Month        string `bson:"month"`
	TotalRevenue int32  `bson:"total_revenue"`
}

type OrderYearlyRevenueRow struct {
	Year         string `bson:"year"`
	TotalRevenue int32  `bson:"total_revenue"`
}

type OrderItemRow struct {
	OrderItemID int32     `bson:"order_item_id"`
	OrderID     int32     `bson:"order_id"`
	ProductID   int32     `bson:"product_id"`
	Quantity    int32     `bson:"quantity"`
	Price       int32     `bson:"price"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

type OrderItemListRow struct {
	OrderItemRow
	TotalCount int64 `bson:"total_count"`
}

type OrderItemActiveRow struct {
	OrderItemRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type OrderItemTrashedRow = OrderItemActiveRow

type ShippingAddressRow struct {
	ShippingAddressID int32     `bson:"shipping_address_id"`
	OrderID           int32     `bson:"order_id"`
	Alamat            string    `bson:"alamat"`
	Provinsi          string    `bson:"provinsi"`
	Negara            string    `bson:"negara"`
	Kota              string    `bson:"kota"`
	Courier           string    `bson:"courier"`
	ShippingMethod    string    `bson:"shipping_method"`
	ShippingCost      float64   `bson:"shipping_cost"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
}

type ShippingAddressListRow struct {
	ShippingAddressRow
	TotalCount int64 `bson:"total_count"`
}

type ShippingAddressActiveRow struct {
	ShippingAddressRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type ShippingAddressTrashedRow = ShippingAddressActiveRow

type MerchantCertificationRow struct {
	MerchantCertificationID int32     `bson:"merchant_certification_id"`
	MerchantID              int32     `bson:"merchant_id"`
	Title                   string    `bson:"title"`
	Description             *string   `bson:"description"`
	IssuedBy                *string   `bson:"issued_by"`
	IssueDate               time.Time `bson:"issue_date"`
	ExpiryDate              time.Time `bson:"expiry_date"`
	CertificateUrl          *string   `bson:"certificate_url"`
	CreatedAt               time.Time `bson:"created_at"`
	UpdatedAt               time.Time `bson:"updated_at"`
}

type MerchantCertificationListRow struct {
	MerchantCertificationRow
	MerchantName string `bson:"merchant_name"`
	TotalCount   int64  `bson:"total_count"`
}

type MerchantCertificationActiveRow struct {
	MerchantCertificationRow
	TotalCount   int64      `bson:"total_count"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty"`
	MerchantName string     `bson:"merchant_name"`
}

type MerchantCertificationTrashedRow = MerchantCertificationActiveRow

type MerchantBusinessInfoRow struct {
	MerchantBusinessInfoID int32     `bson:"merchant_business_info_id"`
	MerchantID             int32     `bson:"merchant_id"`
	BusinessType           *string   `bson:"business_type"`
	TaxID                  *string   `bson:"tax_id"`
	EstablishedYear        *int32    `bson:"established_year"`
	NumberOfEmployees      *int32    `bson:"number_of_employees"`
	WebsiteUrl             *string   `bson:"website_url"`
	CreatedAt              time.Time `bson:"created_at"`
	UpdatedAt              time.Time `bson:"updated_at"`
}

type MerchantBusinessInfoListRow struct {
	MerchantBusinessInfoRow
	MerchantName string `bson:"merchant_name"`
	TotalCount   int64  `bson:"total_count"`
}

type MerchantBusinessInfoActiveRow struct {
	MerchantBusinessInfoRow
	TotalCount   int64      `bson:"total_count"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty"`
	MerchantName string     `bson:"merchant_name"`
}

type MerchantBusinessInfoTrashedRow = MerchantBusinessInfoActiveRow

type MerchantDetailRow struct {
	MerchantDetailID int32     `bson:"merchant_detail_id"`
	MerchantID       int32     `bson:"merchant_id"`
	DisplayName      *string   `bson:"display_name"`
	CoverImageUrl    *string   `bson:"cover_image_url"`
	LogoUrl          *string   `bson:"logo_url"`
	ShortDescription *string   `bson:"short_description"`
	WebsiteUrl       *string   `bson:"website_url"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

type MerchantDetailWithJoinsRow struct {
	MerchantDetailRow
	MerchantName     string `bson:"merchant_name"`
	SocialMediaLinks string `bson:"social_media_links"`
}

type MerchantDetailListRow struct {
	MerchantDetailRow
	MerchantName     string `bson:"merchant_name"`
	TotalCount       int64  `bson:"total_count"`
	SocialMediaLinks string `bson:"social_media_links"`
}

type MerchantDetailActiveRow struct {
	MerchantDetailRow
	TotalCount       int64      `bson:"total_count"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty"`
	MerchantName     string     `bson:"merchant_name"`
	SocialMediaLinks string     `bson:"social_media_links"`
}

type MerchantDetailTrashedRow = MerchantDetailActiveRow

type MerchantPolicyRow struct {
	MerchantPolicyID int32     `bson:"merchant_policy_id"`
	MerchantID       int32     `bson:"merchant_id"`
	PolicyType       string    `bson:"policy_type"`
	Title            string    `bson:"title"`
	Description      string    `bson:"description"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

type MerchantPolicyListRow struct {
	MerchantPolicyRow
	MerchantName string `bson:"merchant_name"`
	TotalCount   int64  `bson:"total_count"`
}

type MerchantPolicyActiveRow struct {
	MerchantPolicyRow
	TotalCount   int64      `bson:"total_count"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty"`
	MerchantName string     `bson:"merchant_name"`
}

type MerchantPolicyTrashedRow = MerchantPolicyActiveRow

type MerchantSocialMediaLinkRow struct {
	MerchantSocialID int32     `bson:"merchant_social_id"`
	MerchantDetailID int32     `bson:"merchant_detail_id"`
	Platform         string    `bson:"platform"`
	Url              string    `bson:"url"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

type TransactionRow struct {
	TransactionID int32     `bson:"transaction_id"`
	OrderID       int32     `bson:"order_id"`
	MerchantID    int32     `bson:"merchant_id"`
	PaymentMethod string    `bson:"payment_method"`
	Amount        int32     `bson:"amount"`
	PaymentStatus string    `bson:"payment_status"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

type TransactionListRow struct {
	TransactionRow
	TotalCount int64      `bson:"total_count"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
}

type TransactionActiveRow = TransactionListRow
type TransactionTrashedRow = TransactionListRow

type TransactionByMerchantRow = TransactionListRow

type TransactionStatsAmountRow struct {
	Year         string  `bson:"year"`
	Month        string  `bson:"month"`
	TotalSuccess int64   `bson:"total_success"`
	TotalFailed  int64   `bson:"total_failed"`
	TotalAmount  float64 `bson:"total_amount"`
}

type TransactionStatsMethodRow struct {
	Month             string  `bson:"month"`
	Year              string  `bson:"year"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionMonthlyAmountSuccessRow struct {
	Year         string `bson:"year"`
	Month        string `bson:"month"`
	TotalSuccess int64  `bson:"total_success"`
	TotalAmount  int32  `bson:"total_amount"`
}

type TransactionMonthlyAmountFailedRow struct {
	Year        string `bson:"year"`
	Month       string `bson:"month"`
	TotalFailed int64  `bson:"total_failed"`
	TotalAmount int32  `bson:"total_amount"`
}

type TransactionYearlyAmountSuccessRow struct {
	Year         string `bson:"year"`
	TotalSuccess int64  `bson:"total_success"`
	TotalAmount  int32  `bson:"total_amount"`
}

type TransactionYearlyAmountFailedRow struct {
	Year        string `bson:"year"`
	TotalFailed int64  `bson:"total_failed"`
	TotalAmount int32  `bson:"total_amount"`
}

type TransactionMonthlyMethodSuccessRow struct {
	Month             string  `bson:"month"`
	Year              string  `bson:"year"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionMonthlyMethodFailedRow = TransactionMonthlyMethodSuccessRow
type TransactionYearlyMethodSuccessRow struct {
	Year              string  `bson:"year"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionYearlyMethodFailedRow = TransactionYearlyMethodSuccessRow

type TransactionMonthlyAmountSuccessByMerchantRow = TransactionMonthlyAmountSuccessRow
type TransactionMonthlyAmountFailedByMerchantRow = TransactionMonthlyAmountFailedRow
type TransactionYearlyAmountSuccessByMerchantRow = TransactionYearlyAmountSuccessRow
type TransactionYearlyAmountFailedByMerchantRow = TransactionYearlyAmountFailedRow
type TransactionMonthlyMethodSuccessByMerchantRow = TransactionMonthlyMethodSuccessRow
type TransactionMonthlyMethodFailedByMerchantRow = TransactionMonthlyMethodFailedRow
type TransactionYearlyMethodSuccessByMerchantRow = TransactionYearlyMethodSuccessRow
type TransactionYearlyMethodFailedByMerchantRow = TransactionYearlyMethodFailedRow

type UserRow struct {
	UserID     int32     `bson:"user_id"`
	Firstname  string    `bson:"firstname"`
	Lastname   string    `bson:"lastname"`
	Email      string    `bson:"email"`
	Password   string    `bson:"password"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	TotalCount int64     `bson:"total_count"`
}

type UserActiveRow struct {
	UserRow
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
}

type UserTrashedRow = UserActiveRow

type UserByIDRow struct {
	UserID    int32     `bson:"user_id"`
	Firstname string    `bson:"firstname"`
	Lastname  string    `bson:"lastname"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UserByEmailWithPasswordRow struct {
	UserID   int32  `bson:"user_id"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
