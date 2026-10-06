package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := map[string][]mongo.IndexModel{
		"users": {
			{Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"roles": {
			{Keys: bson.D{{Key: "role_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "role_name", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"user_roles": {
			{Keys: bson.D{{Key: "user_role_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "role_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"refresh_tokens": {
			{Keys: bson.D{{Key: "refresh_token_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "token", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"reset_tokens": {
			{Keys: bson.D{{Key: "reset_token_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "token", Value: 1}}},
		},
		"merchants": {
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"merchant_documents": {
			{Keys: bson.D{{Key: "document_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"merchant_business_information": {
			{Keys: bson.D{{Key: "merchant_business_info_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"merchant_certifications_and_awards": {
			{Keys: bson.D{{Key: "merchant_certification_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"merchant_details": {
			{Keys: bson.D{{Key: "merchant_detail_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"merchant_social_media_links": {
			{Keys: bson.D{{Key: "merchant_social_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_detail_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"merchant_policies": {
			{Keys: bson.D{{Key: "merchant_policy_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"banners": {
			{Keys: bson.D{{Key: "banner_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"categories": {
			{Keys: bson.D{{Key: "category_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"carts": {
			{Keys: bson.D{{Key: "cart_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "product_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"orders": {
			{Keys: bson.D{{Key: "order_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"order_items": {
			{Keys: bson.D{{Key: "order_item_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "order_id", Value: 1}}},
			{Keys: bson.D{{Key: "product_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"order_stock_reservations": {
			{Keys: bson.D{{Key: "reservation_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "order_id", Value: 1}}},
			{Keys: bson.D{{Key: "product_id", Value: 1}}},
		},
		"products": {
			{Keys: bson.D{{Key: "product_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "category_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"product_stock_adjustments": {
			{Keys: bson.D{{Key: "operation_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "product_id", Value: 1}}},
		},
		"reviews": {
			{Keys: bson.D{{Key: "review_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "product_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"review_details": {
			{Keys: bson.D{{Key: "review_detail_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "review_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"shipping_addresses": {
			{Keys: bson.D{{Key: "shipping_address_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "order_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"sliders": {
			{Keys: bson.D{{Key: "slider_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"transactions": {
			{Keys: bson.D{{Key: "transaction_id", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "order_id", Value: 1}}},
			{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
			{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
		},
		"consumer_inbox": {
			{Keys: bson.D{{Key: "consumer_name", Value: 1}, {Key: "event_key", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "status", Value: 1}}},
			{Keys: bson.D{{Key: "lease_until", Value: 1}}},
		},
		"outbox_events": {
			{Keys: bson.D{{Key: "topic", Value: 1}, {Key: "status", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: 1}}},
		},
	}

	for collectionName, models := range indexes {
		coll := db.Collection(collectionName)
		_, err := coll.Indexes().CreateMany(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}
