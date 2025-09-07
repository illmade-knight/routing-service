// Package test provides public helpers for running end-to-end tests against the service.
package test

import (
	"cloud.google.com/go/pubsub/v2"
	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	psadapter "github.com/illmade-knight/routing-service/internal/platform/pubsub"
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/rs/zerolog"
)

// NewTestConsumer creates a concrete GooglePubsubConsumer for testing purposes.
// It returns the generic interface type to hide the implementation.
func NewTestConsumer(
	subID string,
	client *pubsub.Client,
	logger zerolog.Logger,
) (messagepipeline.MessageConsumer, error) {
	cfg := messagepipeline.NewGooglePubsubConsumerDefaults(subID)
	return messagepipeline.NewGooglePubsubConsumer(cfg, client, logger)
}

// NewTestProducer creates a concrete PubsubProducer for testing purposes.
// It returns the generic interface type.
func NewTestProducer(topic *pubsub.Publisher) routing.IngestionProducer {
	return psadapter.NewProducer(topic)
}
