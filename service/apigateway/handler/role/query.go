package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache"
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/role"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/httpx"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	pb "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/pbcompat"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/role"
	"github.com/go-chi/chi/v5"
)

type roleQueryHandlerApi struct {
	kafka      *kafka.Kafka
	role       pb.RoleQueryServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.RoleQueryResponseMapper
	cache      role_cache.RoleQueryCache
	apiHandler apierror.ApiHandler
}

type roleQueryHandleDeps struct {
	client     pb.RoleQueryServiceClient
	router     chi.Router
	logger     logger.LoggerInterface
	mapper     apimapper.RoleQueryResponseMapper
	kafka      *kafka.Kafka
	cache_role mencache.RoleCache
	cache      role_cache.RoleQueryCache
	apiHandler apierror.ApiHandler
}

func NewRoleQueryHandleApi(params *roleQueryHandleDeps) *roleQueryHandlerApi {
	handler := &roleQueryHandlerApi{
		role:       params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		kafka:      params.kafka,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	roleMiddleware := middlewares.RoleValidatorGRPC(params.client, params.logger, params.cache_role)
	params.router.Route("/api/role-query", func(routerRole chi.Router) {
		requireAdmin := middlewares.RequireRoles("Admin", "ROLE_ADMIN", "Admin_Role_10")

		routerRole.With(roleMiddleware, requireAdmin).Get("/", httpx.Handler(handler.FindAll))
		routerRole.With(roleMiddleware, requireAdmin).Get("/{id}", httpx.Handler(handler.FindById))
		routerRole.With(roleMiddleware, requireAdmin).Get("/active", httpx.Handler(handler.FindByActive))
		routerRole.With(roleMiddleware, requireAdmin).Get("/trashed", httpx.Handler(handler.FindByTrashed))
		routerRole.With(roleMiddleware, requireAdmin).Get("/user/{user_id}", httpx.Handler(handler.FindByUserId))

	})
	return handler
}

// @Security Bearer
// @Summary Find all roles
// @Tags Role Query
// @Description Retrieve a list of all roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRole "List of roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query [get]
func (h *roleQueryHandlerApi) FindAll(w http.ResponseWriter, r *http.Request) error {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := r.URL.Query().Get("search")

	ctx := r.Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoles(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.role.FindAllRole(ctx, &pb.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationRole(res)
	h.cache.SetCachedRoles(ctx, req, apiResponse)

	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find role by ID
// @Tags Role Query
// @Description Retrieve a role by ID
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} errors.ErrorResponse "Invalid role ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/{id} [get]
func (h *roleQueryHandlerApi) FindById(w http.ResponseWriter, r *http.Request) error {
	roleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("id is required")
	}

	ctx := r.Context()
	if cachedData, found := h.cache.GetCachedRoleById(ctx, roleID); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.role.FindByIdRole(ctx, &pb.FindByIdRoleRequest{RoleId: int32(roleID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseRole(res)
	h.cache.SetCachedRoleById(ctx, roleID, apiResponse)

	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active roles
// @Tags Role Query
// @Description Retrieve a list of active roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/active [get]
func (h *roleQueryHandlerApi) FindByActive(w http.ResponseWriter, r *http.Request) error {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := r.URL.Query().Get("search")

	ctx := r.Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoleActive(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.role.FindByActive(ctx, &pb.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)
	h.cache.SetCachedRoleActive(ctx, req, apiResponse)

	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed roles
// @Tags Role Query
// @Description Retrieve a list of trashed role records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed role data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/trashed [get]
func (h *roleQueryHandlerApi) FindByTrashed(w http.ResponseWriter, r *http.Request) error {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := r.URL.Query().Get("search")

	ctx := r.Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoleTrashed(ctx, req); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.role.FindByTrashed(ctx, &pb.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)
	h.cache.SetCachedRoleTrashed(ctx, req, apiResponse)

	return httpx.JSON(w, http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find roles by user ID
// @Tags Role Query
// @Description Retrieve a list of roles assigned to a specific user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.ApiResponsesRole "List of user roles"
// @Failure 400 {object} errors.ErrorResponse "Invalid user ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/user/{user_id} [get]
func (h *roleQueryHandlerApi) FindByUserId(w http.ResponseWriter, r *http.Request) error {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil || userID <= 0 {
		return errors.NewBadRequestError("user_id is required")
	}

	ctx := r.Context()
	if cachedData, found := h.cache.GetCachedRoleByUserId(ctx, userID); found {
		return httpx.JSON(w, http.StatusOK, cachedData)
	}

	res, err := h.role.FindByUserId(ctx, &pb.FindByIdUserRoleRequest{UserId: int32(userID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsesRole(res)
	h.cache.SetCachedRoleByUserId(ctx, userID, apiResponse)

	return httpx.JSON(w, http.StatusOK, apiResponse)
}
