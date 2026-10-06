package outbox

import (
	"context"
	"errors"
	"testing"
)

func TestReserveValidatesKeys(t *testing.T) {
	for name, tc := range map[string]struct {
		consumerName string
		eventKey     string
	}{
		"empty consumer":  {consumerName: "", eventKey: "topic:evt-1"},
		"empty event key": {consumerName: "email-service-group", eventKey: ""},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, _, err := Reserve(context.Background(), nil, tc.consumerName, tc.eventKey, "topic", 0, 1); !errors.Is(err, ErrInvalidInboxKey) {
				t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
			}
		})
	}
}

func TestMarkProcessedValidatesKeys(t *testing.T) {
	if err := MarkProcessed(context.Background(), nil, "email-service-group", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := MarkProcessed(context.Background(), nil, "", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestReleaseValidatesKeys(t *testing.T) {
	if err := Release(context.Background(), nil, "email-service-group", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := Release(context.Background(), nil, "", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestNewMongoInbox_NilDB(t *testing.T) {
	if _, err := NewMongoInbox(nil); err == nil {
		t.Fatal("expected error when creating MongoInbox with nil db")
	}
}

func TestReleasePassesProcessingError(t *testing.T) {
	// Verify the function signature accepts a processing error (nil and non-nil).
	// The actual behavior requires a live MongoDB, so we only test nil-db validation here.
	if err := Release(context.Background(), nil, "group", "key", 1, errors.New("smtp down")); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}
