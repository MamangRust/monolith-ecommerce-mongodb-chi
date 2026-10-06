package merchantsociallinkhandler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_social_link"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsMerchantSocialLink struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	ApiHandler apierror.ApiHandler
}

func RegisterMerchantSocialLinkHandler(deps *DepsMerchantSocialLink) {
	mapper := apimapper.NewMerchantSocialLinkResponseMapper()

	NewMerchantSocialLinkCommandHandleApi(&merchantSocialLinkCommandHandleDeps{
		client:     pb.NewMerchantSocialCommandServiceClient(deps.Client),
		router:     deps.Router,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		apiHandler: deps.ApiHandler,
	})
}
