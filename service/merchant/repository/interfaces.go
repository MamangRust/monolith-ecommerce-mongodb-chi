package repository

import (
	"context"

	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantResult struct {
	MerchantID   int32
	UserID       int32
	Name         string
	Description  *string
	Address      *string
	ContactEmail *string
	ContactPhone *string
	Status       string
	CreatedAt    *string
	UpdatedAt    *string
	DeletedAt    *string
	TotalCount   int64
}

type MerchantDocumentResult struct {
	DocumentID   int32
	MerchantID   int32
	DocumentType string
	DocumentUrl  string
	Status       string
	Note         *string
	UploadedAt   *string
	CreatedAt    *string
	UpdatedAt    *string
	DeletedAt    *string
	TotalCount   int64
}

type MerchantDocumentQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
	FindByID(ctx context.Context, id int) (*MerchantDocumentResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
}

type MerchantDocumentCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*salemongo.MerchantDocument, error)
	Trash(ctx context.Context, merchant_document_id int) (*salemongo.MerchantDocument, error)
	Restore(ctx context.Context, merchant_document_id int) (*salemongo.MerchantDocument, error)
	DeletePermanent(ctx context.Context, merchant_document_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error)
	FindByID(ctx context.Context, user_id int) (*MerchantResult, error)
}

type MerchantCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantRequest) (*salemongo.Merchant, error)
	Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*salemongo.Merchant, error)
	Trash(ctx context.Context, merchant_id int) (*salemongo.Merchant, error)
	Restore(ctx context.Context, merchant_id int) (*salemongo.Merchant, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*salemongo.Merchant, error)
}

type UserQueryRepository = useradapter.QueryRepository
