package user

import (
	"context"
	"time"

	pbuser "github.com/MamangRust/monolith-ecommerce-pb/users"
	adapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	user_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
)

// QueryRepository is the contract consumers depend on for user reads.
type QueryRepository interface {
	FindByID(ctx context.Context, userID int) (*mongo.User, error)
	FindByEmail(ctx context.Context, email string) (*mongo.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*mongo.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*mongo.User, error)
}

// CommandRepository is the contract consumers depend on for user writes.
type CommandRepository interface {
	CreateUser(ctx context.Context, req *requests.RegisterRequest) (*mongo.User, error)
	UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*mongo.User, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) (*mongo.User, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated user query/command clients.
type Repository struct {
	query   pbuser.UserQueryServiceClient
	command pbuser.UserCommandServiceClient
	guard   *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (a *Repository) SetGuard(g *resilience.DependencyGuard) { a.guard = g }

// NewAdapter wraps the generated user clients. Resilience (timeout/circuit-breaker/
// bulkhead) is applied per call through the optional dependency guard; a nil
// guard is a passthrough.
func NewAdapter(query pbuser.UserQueryServiceClient, command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	a := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// New wraps the generated user query/command clients. It is the canonical
// constructor; modules build it inside their repositories with the guard
// options they were given.
func New(query pbuser.UserQueryServiceClient, command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return NewAdapter(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface
// for consumers that only read.
func NewQueryAdapter(query pbuser.UserQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return NewAdapter(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository
// surface for consumers that only write.
func NewCommandAdapter(command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return NewAdapter(nil, command, opts...)
}

func (a *Repository) FindByID(ctx context.Context, userID int) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUser, error) {
		return a.query.FindById(ctx, &pbuser.FindByIdUserRequest{Id: int32(userID)})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) FindByEmail(ctx context.Context, email string) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return a.query.FindByEmail(ctx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if res == nil || res.Data == nil {
		return nil, nil
	}
	return &mongo.User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
		CreatedAt: parseTime(res.Data.CreatedAt),
		UpdatedAt: parseTime(res.Data.UpdatedAt),
	}, nil
}

// FindByEmailAndVerify mirrors the legacy auth repo behaviour: lookup by email
// and surface a "user not found for verification" error if the gRPC RPC fails.
// Unlike FindByEmail, it does NOT short-circuit on codes.NotFound because the
// auth service expects a single failure mode for both pathways.
func (a *Repository) FindByEmailAndVerify(ctx context.Context, email string) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return a.query.FindByEmail(ctx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if res == nil || res.Data == nil {
		return nil, nil
	}
	return &mongo.User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
		CreatedAt: parseTime(res.Data.CreatedAt),
		UpdatedAt: parseTime(res.Data.UpdatedAt),
	}, nil
}

func (a *Repository) FindByVerificationCode(ctx context.Context, code string) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUser, error) {
		return a.query.FindByVerificationCode(ctx, &pbuser.FindByVerificationCodeRequest{VerificationCode: code})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) CreateUser(ctx context.Context, req *requests.RegisterRequest) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUser, error) {
		return a.command.Create(ctx, &pbuser.CreateUserRequest{
			Firstname:       req.FirstName,
			Lastname:        req.LastName,
			Email:           req.Email,
			Password:        req.Password,
			ConfirmPassword: req.ConfirmPassword,
		})
	})
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUser, error) {
		return a.command.UpdateIsVerified(ctx, &pbuser.UpdateUserIsVerifiedRequest{
			Id:         int32(userID),
			IsVerified: isVerified,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func (a *Repository) UpdateUserPassword(ctx context.Context, userID int, password string) (*mongo.User, error) {
	res, err := adapter.Call(a.guard, ctx, func(ctx context.Context) (*pbuser.ApiResponseUser, error) {
		return a.command.UpdatePassword(ctx, &pbuser.UpdateUserPasswordRequest{
			Id:       int32(userID),
			Password: password,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}
	if res == nil {
		return nil, nil
	}
	return toModel(res.Data), nil
}

func toModel(data *pbuser.UserResponse) *mongo.User {
	if data == nil {
		return nil
	}
	return &mongo.User{
		UserID:    data.Id,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		CreatedAt: parseTime(data.CreatedAt),
		UpdatedAt: parseTime(data.UpdatedAt),
	}
}

// parseTime accepts the loose RFC3339-ish timestamp strings the user service
// emits and returns the zero time for empty or unparseable values.
func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}
