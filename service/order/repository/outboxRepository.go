package repository

import (
	"context"
	"sort"
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

func (r *outboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*outboxpkg.OutboxEvent, error) {
	outboxID, err := salemongo.NextID(ctx, r.db, "outbox_events")
	if err != nil {
		return nil, err
	}

	now := time.Now()
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
	return event, nil
}

func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*outboxpkg.OutboxEvent, error) {
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
	return events, nil
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*outboxpkg.OutboxEvent, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "outbox_id", Value: 1}}).
		SetLimit(int64(limit))
	cursor, err := r.db.Collection("outbox_events").Find(ctx, bson.M{
		"status":      "processing",
		"lease_until": bson.M{"$lt": time.Now()},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*outboxpkg.OutboxEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	// Also claim pending events whose retry window elapsed: atomically flip
	// them to processing with a lease so concurrent relays stay exclusive.
	_ = leaseUntil
	pending, err := r.GetPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	events = append(events, pending...)
	sort.Slice(events, func(i, j int) bool {
		return events[i].OutboxID < events[j].OutboxID
	})
	if len(events) > limit {
		events = events[:limit]
	}

	for _, event := range events {
		_, _ = r.db.Collection("outbox_events").UpdateOne(ctx,
			bson.M{"outbox_id": event.OutboxID},
			bson.M{"$set": bson.M{"status": "processing", "lease_until": leaseUntil}})
	}
	return events, nil
}

func (r *outboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*outboxpkg.OutboxEvent, error) {
	now := time.Now()
	filter := bson.M{"outbox_id": outboxID}
	update := bson.M{"$set": bson.M{"status": "delivered", "updated_at": now}}
	res := r.db.Collection("outbox_events").FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After))
	var event outboxpkg.OutboxEvent
	if err := res.Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *outboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*outboxpkg.OutboxEvent, error) {
	now := time.Now()
	filter := bson.M{"outbox_id": outboxID}
	update := bson.M{
		"$set": bson.M{
			"status":          "pending",
			"next_attempt_at": nextAttemptAt,
			"updated_at":      now,
		},
		"$inc": bson.M{"attempts": int32(1)},
	}
	res := r.db.Collection("outbox_events").FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After))
	var event outboxpkg.OutboxEvent
	if err := res.Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *outboxRepository) MarkDead(ctx context.Context, outboxID int64) (*outboxpkg.OutboxEvent, error) {
	now := time.Now()
	filter := bson.M{"outbox_id": outboxID}
	update := bson.M{"$set": bson.M{"status": "dead", "updated_at": now}}
	res := r.db.Collection("outbox_events").FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After))
	var event outboxpkg.OutboxEvent
	if err := res.Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
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
