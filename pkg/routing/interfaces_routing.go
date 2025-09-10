package routing

import (
	"context"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
)

// IngestionProducer defines the interface for a component that can publish a
// received envelope to the internal message bus for processing.
type IngestionProducer interface {
	Publish(ctx context.Context, envelope *transport.SecureEnvelope) error
}

// DeliveryProducer defines the interface for a component that can publish a message
// to a specific, targeted topic for real-time delivery to an online user.
type DeliveryProducer interface {
	Publish(ctx context.Context, topicID string, data *transport.SecureEnvelope) error
}

// PushNotifier defines the interface for a component that can send push
// notifications to offline users' devices.
type PushNotifier interface {
	Notify(ctx context.Context, tokens []DeviceToken, envelope *transport.SecureEnvelope) error
}

// REFACTOR: Add the new MessageStore interface for offline message persistence.

// MessageStore defines the contract for a persistence layer that stores messages
// for users who are offline.
type MessageStore interface {
	// Store saves a message envelope for a specific user.
	Store(ctx context.Context, userID string, envelope *transport.SecureEnvelope) error

	// FetchUndelivered retrieves all messages for a user that have not yet
	// been delivered.
	FetchUndelivered(ctx context.Context, userID string) ([]*transport.SecureEnvelope, error)

	// MarkDelivered marks a set of messages as delivered so they are not
	// retrieved again. This could be implemented as a delete or a status update.
	MarkDelivered(ctx context.Context, userID string, envelopes []*transport.SecureEnvelope) error
}
