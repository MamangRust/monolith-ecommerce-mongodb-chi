package merchantdocumenthandler

import (
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_documents"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsMerchantDocument struct {
	Client      *grpc.ClientConn
	Router      chi.Router
	Logger      logger.LoggerInterface
	UploadImage upload_image.ImageUploads
}

func RegisterMerchantDocumentHandler(deps *DepsMerchantDocument) {
	mapper := apimapper.NewMerchantDocumentResponseMapper()

	NewMerchantDocumentQueryHandleApi(&merchantDocumentQueryHandleDeps{
		client: pb.NewMerchantDocumentQueryServiceClient(deps.Client),
		router: deps.Router,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
	})

	NewMerchantDocumentCommandHandleApi(&merchantDocumentCommandHandleDeps{
		client:       pb.NewMerchantDocumentCommandServiceClient(deps.Client),
		router:       deps.Router,
		logger:       deps.Logger,
		mapper:       mapper.CommandMapper(),
		upload_image: deps.UploadImage,
	})
}
