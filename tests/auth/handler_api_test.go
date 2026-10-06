package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	auth_cache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/handler"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-auth/service"
	auth_cache_api "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/auth"
	authhandler "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/handler/auth"
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-role/cache"
	role_handler "github.com/MamangRust/monolith-ecommerce-grpc-role/handler"
	role_repo "github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	role_service "github.com/MamangRust/monolith-ecommerce-grpc-role/service"
	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	user_handler "github.com/MamangRust/monolith-ecommerce-grpc-user/handler"
	user_repo "github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	user_service "github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	pb "github.com/MamangRust/monolith-ecommerce-pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/apierror"
	pbroles "github.com/MamangRust/monolith-ecommerce-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-ecommerce-pb/user_roles"
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthHandlerApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	redisClient *redis.Client
	server      chi.Router
	email       string
	password    string
	accessToken string
	userID      int
}

func (s *AuthHandlerApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)
	s.redisClient.FlushAll(context.Background())

	mongoDB, err := s.ts.MongoDB()
	s.Require().NoError(err)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// 1. Setup Role Service & gRPC Server
	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(mongoDB)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repository:    roleRepos,
		Logger:        log,
		Cache:         roleMencache,
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(roleServer, roleGapi.RoleQuery)
	pbroles.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserroles.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)
	roleLis, _ := net.Listen("tcp", "localhost:0")
	go roleServer.Serve(roleLis)
	roleConn, _ := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 2. Setup User Service & gRPC Server
	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{
		Db:              mongoDB,
		RoleQueryClient: pbroles.NewRoleQueryServiceClient(roleConn),
		UserRoleClient:  pbuserroles.NewUserRoleServiceClient(roleConn),
	})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:          hasher,
		Cache:         userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(userServer, userGapi.UserQuery)
	pbusers.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)
	userLis, _ := net.Listen("tcp", "localhost:0")
	go userServer.Serve(userLis)
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 3. Setup Auth Service with gRPC clients
	userAdapter := useradapter.NewAdapter(
		pbusers.NewUserQueryServiceClient(userConn),
		pbusers.NewUserCommandServiceClient(userConn),
	)
	repos := repository.NewRepositories(&repository.Deps{
		Db:              mongoDB,
		User:            userAdapter,
		RoleQueryClient: pbroles.NewRoleQueryServiceClient(roleConn),
		UserRoleClient:  pbuserroles.NewUserRoleServiceClient(roleConn),
	})

	tokenManager, _ := auth.NewManager("mysecret")
	mencache := auth_cache.NewMencache(cacheStore)
	apiAuthCache := auth_cache_api.NewMencache(cacheStore)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	h := handler.NewAuthHandleGrpc(svc, log)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, h)

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.server = chi.NewRouter()
	apiHandler := apierror.NewApiHandler(obs, log)

	// Auth bypass middleware for /api/auth/me — must mimic middlewares/auth.go:
	// it stores the authenticated user id as an int under "user_id".
	s.server.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.userID != 0 {
				r = r.WithContext(context.WithValue(r.Context(), "user_id", s.userID))
			}
			next.ServeHTTP(w, r)
		})
	})

	fmt.Println("Calling RegisterAuthHandler...")
	authhandler.RegisterAuthHandler(&authhandler.DepsAuth{
		Client:     conn,
		Router:     s.server,
		Logger:     log,
		Cache:      apiAuthCache,
		ApiHandler: apiHandler,
	})
	fmt.Println("RegisterAuthHandler called.")

	s.server.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	s.email = "auth.handler.api.test@example.com"
	s.password = "password123"

	// Seed ROLE_ADMIN via gRPC to ensure visibility
	roleCommandClient := pbroles.NewRoleCommandServiceClient(roleConn)
	createdRoleRes, err := roleCommandClient.CreateRole(context.Background(), &pbroles.CreateRoleRequest{
		Name: "ROLE_ADMIN",
	})
	if err != nil {
		fmt.Printf("DEBUG: Seed ROLE_ADMIN gRPC error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Seed ROLE_ADMIN gRPC success, ID: %d\n", createdRoleRes.Data.Id)
	}

	// Verify via direct queries
	roleCount, err := mongoDB.Collection("roles").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		fmt.Printf("DEBUG: direct queries error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: direct queries found %d roles\n", roleCount)
	}

	// Verify via GetRole by ID
	var roleByID bson.M
	err = mongoDB.Collection("roles").FindOne(context.Background(), bson.M{"role_id": createdRoleRes.Data.Id}).Decode(&roleByID)
	if err != nil {
		fmt.Printf("DEBUG: GetRole by ID error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: GetRole by ID found: %s\n", roleByID["role_name"])
	}

	// Verify via gRPC with empty search
	roleClient := pbroles.NewRoleQueryServiceClient(roleConn)
	roleRes, err := roleClient.FindAllRole(context.Background(), &pbroles.FindAllRoleRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("DEBUG: gRPC verify error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: gRPC verify (empty search) found %d roles\n", len(roleRes.Data))
		for _, r := range roleRes.Data {
			fmt.Printf("- gRPC role: %s\n", r.Name)
		}
	}
}

func (s *AuthHandlerApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	s.ts.Teardown()
}

func (s *AuthHandlerApiTestSuite) Test0_Ping() {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("pong", rec.Body.String())
}

func (s *AuthHandlerApiTestSuite) Test0_Hello() {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/hello", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("Hello", rec.Body.String())
}

func (s *AuthHandlerApiTestSuite) Test1_Register() {
	body := map[string]string{
		"firstname":        "Auth",
		"lastname":         "API",
		"email":            s.email,
		"password":         s.password,
		"confirm_password": s.password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusCreated, rec.Code, "Expected StatusCreated, got %d. Body: %s", rec.Code, rec.Body.String())

	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.userID = int(data["id"].(float64))
}

func (s *AuthHandlerApiTestSuite) Test2_Login() {
	body := map[string]string{
		"email":    s.email,
		"password": s.password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())

	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.accessToken = data["access_token"].(string)
}

func (s *AuthHandlerApiTestSuite) Test4_LoginLockout() {
	email := "locked.api@example.com"
	password := "wrongpassword"

	// Register user first
	regBody := map[string]string{
		"firstname":        "Locked",
		"lastname":         "API",
		"email":            email,
		"password":         "correctpassword",
		"confirm_password": "correctpassword",
	}
	jsonRegBody, _ := json.Marshal(regBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonRegBody))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	s.server.ServeHTTP(regRec, regReq)
	s.Equal(http.StatusCreated, regRec.Code)

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonLoginBody, _ := json.Marshal(loginBody)

	// Fail login 5 times (total 5)
	for i := 0; i < 5; i++ {
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonLoginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		s.server.ServeHTTP(loginRec, loginReq)
		s.Equal(http.StatusUnauthorized, loginRec.Code)
	}

	// 6th attempt should return 403 Forbidden (ErrAccountLocked)
	lockedReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonLoginBody))
	lockedReq.Header.Set("Content-Type", "application/json")
	lockedRec := httptest.NewRecorder()
	s.server.ServeHTTP(lockedRec, lockedReq)
	s.Equal(http.StatusForbidden, lockedRec.Code)
}

func (s *AuthHandlerApiTestSuite) Test3_GetMe() {
	s.Require().NotZero(s.userID)
	s.Require().NotEmpty(s.accessToken)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())

	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.Equal(s.email, data["email"])
}

func TestAuthHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerApiTestSuite))
}
