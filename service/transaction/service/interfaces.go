package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type TransactionQueryService interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*repository.TransactionResult, *int, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*repository.TransactionResult, *int, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*repository.TransactionResult, *int, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*repository.TransactionResult, *int, error)

	FindByID(
		ctx context.Context,
		transaction_id int,
	) (*repository.TransactionResult, error)

	FindByOrderID(
		ctx context.Context,
		order_id int,
	) (*repository.TransactionResult, error)
}

type TransactionCommandService interface {
	Create(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*repository.TransactionResult, error)

	Update(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*repository.TransactionResult, error)

	Trash(
		ctx context.Context,
		transaction_id int,
	) (*repository.TransactionResult, error)

	Restore(
		ctx context.Context,
		transaction_id int,
	) (*repository.TransactionResult, error)

	DeletePermanent(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type TransactionStatsService interface {
	FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, error)
	FindYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, error)
	FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, error)
	FindYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, error)
	FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error)
	FindYearlyMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, error)
	FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, error)
	FindYearlyMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, error)
}

type TransactionStatsByMerchantService interface {
	FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error)
	FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error)
	FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error)
	FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error)
	FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, error)
	FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodSuccessByMerchantRow, error)
	FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodFailedByMerchantRow, error)
	FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodFailedByMerchantRow, error)
}
