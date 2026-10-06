package service

import (
	"context"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
)

type transactionStatsByMerchantService struct {
	observability observability.TraceLoggerObservability
	cache         cache.TransactionStatsByMerchantCache
	repository    repository.TransactionStatsByMerchantRepository
	logger        logger.LoggerInterface
}

type TransactionStatsByMerchantServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.TransactionStatsByMerchantCache
	Repository    repository.TransactionStatsByMerchantRepository
	Logger        logger.LoggerInterface
}

func NewTransactionStatsByMerchantService(deps *TransactionStatsByMerchantServiceDeps) *transactionStatsByMerchantService {
	return &transactionStatsByMerchantService{
		observability: deps.Observability,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
	}
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error) {
	const method = "FindMonthlyAmountSuccessByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyAmountSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyAmountSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error) {
	const method = "FindYearlyAmountSuccessByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyAmountSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionYearlyAmountSuccessByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyAmountSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error) {
	const method = "FindMonthlyAmountFailedByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyAmountFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionMonthlyAmountFailedByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyAmountFailedByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error) {
	const method = "FindYearlyAmountFailedByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyAmountFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionYearlyAmountFailedByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyAmountFailedByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, error) {
	const method = "FindMonthlyMethodByMerchantSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyMethodSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyMethodSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodSuccessByMerchantRow, error) {
	const method = "FindYearlyMethodByMerchantSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyMethodSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionYearlyMethodSuccessByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyMethodSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodFailedByMerchantRow, error) {
	const method = "FindMonthlyMethodByMerchantFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyMethodFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionMonthlyMethodFailedByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyMethodFailedByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodFailedByMerchantRow, error) {
	const method = "FindYearlyMethodByMerchantFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyMethodFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*salemongo.TransactionYearlyMethodFailedByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyMethodFailedByMerchant(ctx, req, res)
	return res, nil
}
