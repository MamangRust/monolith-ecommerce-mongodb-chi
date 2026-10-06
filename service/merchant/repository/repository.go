package repository

import (
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/users"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency. Callers build them with adapter.WithDependencyGuard (typically
// around resilience.NewDependencyGuard) so the repository owns the wiring.
type GuardOptions struct {
	User []adapter.GuardOption
}

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

func NewRepositories(db *mongo.Database, userQuery pbuser.UserQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions

	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(db.Collection("merchants")),
		MerchantCommand:         NewMerchantCommandRepository(db, db.Collection("merchants")),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(db.Collection("merchant_documents")),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(db.Collection("merchant_documents")),
		UserQuery:               useradapter.NewQueryAdapter(userQuery, g.User...),
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
