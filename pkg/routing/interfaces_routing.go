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
