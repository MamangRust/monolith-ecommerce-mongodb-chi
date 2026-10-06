package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	pbcarts "github.com/MamangRust/monolith-ecommerce-pb/carts"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

type cartQueryHandler struct {
	pbcarts.UnimplementedCartQueryServiceServer
	cartQuery service.CartQueryService
	logger    logger.LoggerInterface
}

func NewCartQueryHandler(cartQuery service.CartQueryService, logger logger.LoggerInterface) *cartQueryHandler {
	return &cartQueryHandler{
		cartQuery: cartQuery,
		logger:    logger,
	}
}

func (h *cartQueryHandler) FindAll(ctx context.Context, request *pbcarts.FindAllCartRequest) (*pbcarts.ApiResponsePaginationCart, error) {
	userID := int(request.GetUserId())
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllCarts{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cartItems, totalRecords, err := h.cartQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCartItems := make([]*pbcarts.CartResponse, len(cartItems))
	for i, cartItem := range cartItems {
		protoCartItems[i] = mapToProtoCartResponse(cartItem)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbcarts.ApiResponsePaginationCart{
		Status:     "success",
		Message:    "Successfully fetched cart items",
		Data:       protoCartItems,
		Pagination: paginationMeta,
	}, nil
}
