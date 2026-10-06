package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"
	pbmerchants "github.com/MamangRust/monolith-ecommerce-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-ecommerce-pb/merchants_document"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantQuery           pbmerchants.MerchantQueryServiceServer
	MerchantCommandHandler  pbmerchants.MerchantCommandServiceServer
	MerchantDocumentQuery   pbmerchants_document.MerchantDocumentQueryServiceServer
	MerchantDocumentCommand pbmerchants_document.MerchantDocumentCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantQuery:           NewMerchantQueryHandler(deps.Service.MerchantQuery, deps.Logger),
		MerchantCommandHandler:  NewMerchantCommandHandler(deps.Service.MerchantCommand, deps.Logger),
		MerchantDocumentQuery:   NewMerchantDocumentQueryHandler(deps.Service.MerchantDocumentQuery, deps.Logger),
		MerchantDocumentCommand: NewMerchantDocumentCommandHandler(deps.Service.MerchantDocumentCommand, deps.Logger),
	}
}
