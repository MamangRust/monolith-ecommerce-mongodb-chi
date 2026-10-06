package tests

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type TestSuite struct {
	MongoContainer *mongodb.MongoDBContainer
	RedisContainer *redis.RedisContainer
	MongoURI       string
	RedisURL       string
	Ctx            context.Context
	client         *mongo.Client
	db             *mongo.Database
}

func SetupTestSuite() (*TestSuite, error) {
	ctx := context.Background()

	// Setup MongoDB
	mongoContainer, err := mongodb.Run(ctx,
		"mongo:7.0",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("27017").WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	mongoHost, err := mongoContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get mongo host: %w", err)
	}

	mongoPort, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		return nil, fmt.Errorf("failed to get mongo port: %w", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort.Port())

	// Setup Redis
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	redisURL, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis connection string: %w", err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	db := client.Database("testdb")
	if err := salemongo.CreateIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create mongo indexes: %w", err)
	}

	ts := &TestSuite{
		MongoContainer: mongoContainer,
		RedisContainer: redisContainer,
		MongoURI:       mongoURI,
		RedisURL:       redisURL,
		Ctx:            ctx,
		client:         client,
		db:             db,
	}

	return ts, nil
}

func (ts *TestSuite) MongoDB() (*mongo.Database, error) {
	return ts.db, nil
}

func (ts *TestSuite) RedisClient() *goredis.Client {
	opts, _ := goredis.ParseURL(ts.RedisURL)
	return goredis.NewClient(opts)
}

func (ts *TestSuite) Teardown() {
	if ts.client != nil {
		if err := ts.client.Disconnect(ts.Ctx); err != nil {
			log.Printf("failed to disconnect mongo client: %v", err)
		}
	}
	if ts.MongoContainer != nil {
		if err := ts.MongoContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate mongo container: %v", err)
		}
	}
	if ts.RedisContainer != nil {
		if err := ts.RedisContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate redis container: %v", err)
		}
	}
}

type GRPCServerRunner interface {
	Serve(lis net.Listener) error
	Stop()
	GracefulStop()
}

func RunGRPCServer(server *grpc.Server) (string, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()
	return lis.Addr().String(), nil
}
