package handler

import (
	pbtransactions "github.com/MamangRust/monolith-ecommerce-pb/transactions"
)

type TransactionQueryHandler interface {
	pbtransactions.TransactionQueryServiceServer
}

type TransactionCommandHandler interface {
	pbtransactions.TransactionCommandServiceServer
}

type TransactionStatsHandler interface {
	pbtransactions.TransactionStatsServiceServer
}

type TransactionStatsByMerchantHandler interface {
	pbtransactions.TransactionStatsByMerchantServiceServer
}
