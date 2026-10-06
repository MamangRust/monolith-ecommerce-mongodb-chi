package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	outboxpkg "github.com/MamangRust/monolith-ecommerce-pkg/outbox"
)

type outboxRepository struct {
	db *mongo.Database
}

func NewOutboxRepository(db *mongo.Database) OutboxRepository {
	return &outboxRepository{db: db}
}

func toOutboxResult(e *outboxpkg.OutboxEvent) *OutboxEventResult {
	if e == nil {
		return nil
	}
	return &OutboxEventResult{
		OutboxID:      e.OutboxID,
		Topic:         e.Topic,
		EventKey:      e.EventKey,
		Payload:       e.Payload,
		Status:        e.Status,
		Attempts:      e.Attempts,
		NextAttemptAt: e.NextAttemptAt,
		CreatedAt:     &e.CreatedAt,
		UpdatedAt:     &e.UpdatedAt,
	}
}

func (r *outboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*OutboxEventResult, error) {
	now := time.Now()

	outboxID, err := salemongo.NextID(ctx, r.db, "outbox_events")
	if err != nil {
		return nil, err
	}

	event := &outboxpkg.OutboxEvent{
		OutboxID:      int64(outboxID),
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		Attempts:      0,
		NextAttemptAt: now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if _, err := r.db.Collection("outbox_events").InsertOne(ctx, event); err != nil {
		return nil, err
	}
	return toOutboxResult(event), nil
}

func (r *outboxRepository) CreateInTx(ctx context.Context, _ interface{}, topic, key string, payload []byte) (*OutboxEventResult, error) {
	// MongoDB doesn't have gorm transactions; the outbox insert is atomic.
	return r.Create(ctx, topic, key, payload)
}

func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*OutboxEventResult, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "outbox_id", Value: 1}}).
		SetLimit(int64(limit))
	cursor, err := r.db.Collection("outbox_events").Find(ctx, bson.M{
		"status":          "pending",
		"next_attempt_at": bson.M{"$lte": time.Now()},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*outboxpkg.OutboxEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	var results []*OutboxEventResult
	for _, e := range events {
		results = append(results, toOutboxResult(e))
	}
	return results, nil
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*OutboxEventResult, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "outbox_id", Value: 1}}).
		SetLimit(int64(limit))
	cursor, err := r.db.Collection("outbox_events").Find(ctx, bson.M{
		"status":          "pending",
		"next_attempt_at": bson.M{"$lte": time.Now()},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*outboxpkg.OutboxEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	for _, e := range events {
		now := time.Now()
		r.db.Collection("outbox_events").UpdateOne(ctx,
			bson.M{"outbox_id": e.OutboxID},
			bson.M{"$set": bson.M{"next_attempt_at": leaseUntil, "updated_at": now}})
		e.NextAttemptAt = leaseUntil
		e.UpdatedAt = now
	}

	var results []*OutboxEventResult
	for _, e := range events {
		results = append(results, toOutboxResult(e))
	}
	return results, nil
}

func (r *outboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*OutboxEventResult, error) {
	now := time.Now()
	var event outboxpkg.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"outbox_id": outboxID},
		bson.M{"$set": bson.M{"status": "delivered", "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return toOutboxResult(&event), nil
}

func (r *outboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*OutboxEventResult, error) {
	now := time.Now()
	var event outboxpkg.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"outbox_id": outboxID},
		bson.M{
			"$set": bson.M{"next_attempt_at": nextAttemptAt, "updated_at": now},
			"$inc": bson.M{"attempts": int32(1)},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return toOutboxResult(&event), nil
}

func (r *outboxRepository) MarkDead(ctx context.Context, outboxID int64) (*OutboxEventResult, error) {
	now := time.Now()
	var event outboxpkg.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"outbox_id": outboxID},
		bson.M{"$set": bson.M{"status": "dead", "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return toOutboxResult(&event), nil
}

func (r *outboxRepository) DeleteOld(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.Collection("outbox_events").DeleteMany(ctx, bson.M{
		"status":     bson.M{"$in": []string{"delivered", "dead"}},
		"updated_at": bson.M{"$lt": cutoff},
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
