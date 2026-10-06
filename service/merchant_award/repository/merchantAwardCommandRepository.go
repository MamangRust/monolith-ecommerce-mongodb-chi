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

type merchantAwardCommandRepository struct {
	db *mongo.Database
}

func NewMerchantAwardCommandRepository(db *mongo.Database) MerchantAwardCommandRepository {
	return &merchantAwardCommandRepository{db: db}
}

func (r *merchantAwardCommandRepository) coll() *mongo.Collection {
	return r.db.Collection("merchant_certifications_and_awards")
}

func (r *merchantAwardCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*salemongo.MerchantCertificationRow, error) {
	now := time.Now()
	id, err := salemongo.NextID(ctx, r.db, "merchant_certifications_and_awards")
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward
	}

	issueDate, _ := time.Parse("2006-01-02", req.IssueDate)
	expiryDate, _ := time.Parse("2006-01-02", req.ExpiryDate)

	doc := &salemongo.MerchantCertificationsAndAward{
		MerchantCertificationID: id,
		MerchantID:              int32(req.MerchantID),
		Title:                   req.Title,
		Description:             &req.Description,
		IssuedBy:                &req.IssuedBy,
		IssueDate:               issueDate,
		ExpiryDate:              expiryDate,
		CertificateUrl:          &req.CertificateUrl,
		CreatedAt:               &now,
		UpdatedAt:               &now,
	}
	if _, err := r.coll().InsertOne(ctx, doc); err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward
	}
	return awardToRow(doc), nil
}

func awardToRow(a *salemongo.MerchantCertificationsAndAward) *salemongo.MerchantCertificationRow {
	return &salemongo.MerchantCertificationRow{
		MerchantCertificationID: a.MerchantCertificationID,
		MerchantID:              a.MerchantID,
		Title:                   a.Title,
		Description:             a.Description,
		IssuedBy:                a.IssuedBy,
		IssueDate:               a.IssueDate,
		ExpiryDate:              a.ExpiryDate,
		CertificateUrl:          a.CertificateUrl,
		CreatedAt:               derefAwardTime(a.CreatedAt),
		UpdatedAt:               derefAwardTime(a.UpdatedAt),
	}
}

func derefAwardTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func (r *merchantAwardCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*salemongo.MerchantCertificationRow, error) {
	now := time.Now()
	issueDate, _ := time.Parse("2006-01-02", req.IssueDate)
	expiryDate, _ := time.Parse("2006-01-02", req.ExpiryDate)

	filter := bson.M{"merchant_certification_id": int32(*req.MerchantCertificationID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{
		"title":           req.Title,
		"description":     req.Description,
		"issued_by":       req.IssuedBy,
		"issue_date":      issueDate,
		"expiry_date":     expiryDate,
		"certificate_url": req.CertificateUrl,
		"updated_at":      now,
	}}

	var doc salemongo.MerchantCertificationsAndAward
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward
	}
	return awardToRow(&doc), nil
}

func (r *merchantAwardCommandRepository) Trash(ctx context.Context, merchantID int) (*salemongo.MerchantCertificationsAndAward, error) {
	now := time.Now()
	filter := bson.M{"merchant_id": int32(merchantID), "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": now}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var doc salemongo.MerchantCertificationsAndAward
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward
	}
	return &doc, nil
}

func (r *merchantAwardCommandRepository) Restore(ctx context.Context, merchantID int) (*salemongo.MerchantCertificationsAndAward, error) {
	filter := bson.M{"merchant_id": int32(merchantID), "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var doc salemongo.MerchantCertificationsAndAward
	if err := r.coll().FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward
	}
	return &doc, nil
}

func (r *merchantAwardCommandRepository) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	if _, err := r.coll().DeleteMany(ctx, bson.M{"merchant_id": int32(merchantID)}); err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.coll().UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	return err == nil, err
}

func (r *merchantAwardCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.coll().DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	return err == nil, err
}
