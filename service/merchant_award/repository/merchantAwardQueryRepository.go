package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantAwardQueryRepository struct {
	db *mongo.Database
}

func NewMerchantAwardQueryRepository(db *mongo.Database) MerchantAwardQueryRepository {
	return &merchantAwardQueryRepository{db: db}
}

func (r *merchantAwardQueryRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_certifications_and_awards")
}

func derefAwardTime2(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func awardToListRow(a *salemongo.MerchantCertificationsAndAward, total int64) *salemongo.MerchantCertificationListRow {
	return &salemongo.MerchantCertificationListRow{
		MerchantCertificationRow: salemongo.MerchantCertificationRow{
			MerchantCertificationID: a.MerchantCertificationID,
			MerchantID:              a.MerchantID,
			Title:                   a.Title,
			Description:             a.Description,
			IssuedBy:                a.IssuedBy,
			IssueDate:               a.IssueDate,
			ExpiryDate:              a.ExpiryDate,
			CertificateUrl:          a.CertificateUrl,
			CreatedAt:               derefAwardTime2(a.CreatedAt),
			UpdatedAt:               derefAwardTime2(a.UpdatedAt),
		},
		TotalCount: total,
	}
}

func (r *merchantAwardQueryRepository) findAwards(ctx context.Context, filter bson.M, page, pageSize int) ([]*salemongo.MerchantCertificationListRow, error) {
	page = max(page, 1)
	pageSize = max(pageSize, 1)
	offset := (page - 1) * pageSize

	total, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "merchant_certification_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards
	}
	defer cursor.Close(ctx)

	var awards []*salemongo.MerchantCertificationsAndAward
	if err := cursor.All(ctx, &awards); err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards
	}

	var results []*salemongo.MerchantCertificationListRow
	for _, a := range awards {
		results = append(results, awardToListRow(a, total))
	}
	return results, nil
}

func (r *merchantAwardQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationListRow, error) {
	return r.findAwards(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
}

func (r *merchantAwardQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationActiveRow, error) {
	listRows, err := r.findAwards(ctx, bson.M{"deleted_at": nil}, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	var active []*salemongo.MerchantCertificationActiveRow
	for _, row := range listRows {
		active = append(active, &salemongo.MerchantCertificationActiveRow{
			MerchantCertificationRow: row.MerchantCertificationRow,
			DeletedAt:                nil,
			TotalCount:               row.TotalCount,
		})
	}
	return active, nil
}

func (r *merchantAwardQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*salemongo.MerchantCertificationTrashedRow, error) {
	listRows, err := r.findAwards(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	var trashed []*salemongo.MerchantCertificationTrashedRow
	for _, row := range listRows {
		trashed = append(trashed, &salemongo.MerchantCertificationTrashedRow{
			MerchantCertificationRow: row.MerchantCertificationRow,
			TotalCount:               row.TotalCount,
		})
	}
	return trashed, nil
}

func (r *merchantAwardQueryRepository) FindByID(ctx context.Context, merchantID int) (*salemongo.MerchantCertificationRow, error) {
	var a salemongo.MerchantCertificationsAndAward
	err := r.coll().FindOne(ctx, bson.M{"merchant_certification_id": int32(merchantID), "deleted_at": nil}).Decode(&a)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchantaward_errors.ErrMerchantAwardNotFound
		}
		return nil, merchantaward_errors.ErrFindByIdMerchantAward
	}
	return awardToRow(&a), nil
}
