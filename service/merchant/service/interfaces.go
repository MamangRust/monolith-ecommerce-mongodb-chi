package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, error)
	FindByID(ctx context.Context, merchantID int) (*repository.MerchantResult, error)
}

type MerchantCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantRequest) (*salemongo.Merchant, error)
	Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*salemongo.Merchant, error)
	Trash(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
	Restore(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
	DeletePermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*salemongo.Merchant, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindByID(ctx context.Context, documentID int) (*repository.MerchantDocumentResult, error)
}

type MerchantDocumentCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*salemongo.MerchantDocument, error)
	Trash(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error)
	Restore(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error)
	DeletePermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
