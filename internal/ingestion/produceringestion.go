// REFACTOR: This producer is updated to publish directly to a pubsub.Topic
// to prevent a "double-wrapping" serialization error that was occurring when
// it called another generic producer.

// Package ingestion contains components responsible for receiving data from external
// sources (like HTTP or MQTT) and producing it to the internal pipeline.
package ingestion

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
)

// pubsubTopicClient defines the interface for the underlying pubsub.Topic.
// This allows us to use a mock for testing.
type pubsubTopicClient interface {
	Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult
}

// PubsubProducer implements the api.IngestionProducer interface.
// It acts as an adapter, serializing the SecureEnvelope from the API layer and
// publishing it directly to the configured Pub/Sub topic.
type PubsubProducer struct {
	topic pubsubTopicClient
}

// NewPubsubProducer is the constructor for our ingestion producer.
func NewPubsubProducer(topic pubsubTopicClient) *PubsubProducer {
	producer := PubsubProducer{
		topic: topic,
	}
	return &producer
}

// Publish serializes the SecureEnvelope and sends it to the message bus.
func (p PubsubProducer) Publish(ctx context.Context, envelope *transport.SecureEnvelope) error {
	var err error

	// Serialize the envelope to a JSON byte slice for the message payload.
	payloadBytes, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope for publishing: %w", err)
	}

	// Create the pubsub.Message directly.
	message := &pubsub.Message{
		Data: payloadBytes,
	}

	// Publish the message and wait for the result.
	result := p.topic.Publish(ctx, message)
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("ingestion producer failed to publish message: %w", err)
	}

	return nil
}
