package cache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*repository.TransactionResult, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*repository.TransactionResult, bool)
	SetCachedTransactionCache(ctx context.Context, data *repository.TransactionResult)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*repository.TransactionResult, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *repository.TransactionResult)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	InvalidateTransactionCache(ctx context.Context)
}

type TransactionStatsCache interface {
	GetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, bool)
	SetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountSuccessRow)
	GetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, bool)
	SetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountFailedRow)
	GetCachedYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, bool)
	SetCachedYearlyAmountSuccess(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountSuccessRow)
	GetCachedYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, bool)
	SetCachedYearlyAmountFailed(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountFailedRow)
	GetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, bool)
	SetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodSuccessRow)
	GetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, bool)
	SetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodFailedRow)
	GetCachedYearlyMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, bool)
	SetCachedYearlyMethodSuccess(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodSuccessRow)
	GetCachedYearlyMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, bool)
	SetCachedYearlyMethodFailed(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodFailedRow)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, bool)
	SetCachedMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*salemongo.TransactionMonthlyAmountSuccessByMerchantRow)
	GetCachedMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, bool)
	SetCachedMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*salemongo.TransactionMonthlyAmountFailedByMerchantRow)
	GetCachedYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, bool)
	SetCachedYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*salemongo.TransactionYearlyAmountSuccessByMerchantRow)
	GetCachedYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, bool)
	SetCachedYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*salemongo.TransactionYearlyAmountFailedByMerchantRow)
	GetCachedMonthlyMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodSuccessByMerchantRow, bool)
	SetCachedMonthlyMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*salemongo.TransactionMonthlyMethodSuccessByMerchantRow)
	GetCachedMonthlyMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodFailedByMerchantRow, bool)
	SetCachedMonthlyMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*salemongo.TransactionMonthlyMethodFailedByMerchantRow)
	GetCachedYearlyMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodSuccessByMerchantRow, bool)
	SetCachedYearlyMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*salemongo.TransactionYearlyMethodSuccessByMerchantRow)
	GetCachedYearlyMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodFailedByMerchantRow, bool)
	SetCachedYearlyMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*salemongo.TransactionYearlyMethodFailedByMerchantRow)
}
