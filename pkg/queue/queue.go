package queue

import (
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
)

// Queue defines the interface for message persistence.
type Queue interface {
	Enqueue(envelope transport.SecureEnvelope) error
	Dequeue(userID string) ([]transport.SecureEnvelope, error)
}
