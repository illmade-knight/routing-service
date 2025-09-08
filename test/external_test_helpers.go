// Package test provides public helpers for running end-to-end tests against the service.
package test

import (
	"context"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub/v2"
	"github.com/illmade-knight/go-dataflow/pkg/cache"
	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	psadapter "github.com/illmade-knight/routing-service/internal/platform/pubsub"
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/rs/zerolog"
)

// deviceTokenDoc matches the structure of the document stored in Firestore.
type deviceTokenDoc struct {
	Tokens []routing.DeviceToken
}

// FirestoreTokenAdapter wraps the generic Firestore fetcher to satisfy the specific
// interface required by the routing service's dependencies.
type FirestoreTokenAdapter struct {
	docFetcher cache.Fetcher[string, deviceTokenDoc]
}

// Fetch implements the `cache.Fetcher[string, []routing.DeviceToken]` interface.
// It calls the underlying fetcher and extracts the `.Tokens` field.
func (a *FirestoreTokenAdapter) Fetch(ctx context.Context, key string) ([]routing.DeviceToken, error) {
	doc, err := a.docFetcher.Fetch(ctx, key)
	if err != nil {
		return nil, err
	}
	return doc.Tokens, nil
}

// Close satisfies the io.Closer part of the interface.
func (a *FirestoreTokenAdapter) Close() error {
	return a.docFetcher.Close()
}

// --- Test Helper Constructors ---

// NewTestFirestoreTokenFetcher creates a Firestore fetcher and wraps it in the required adapter.
func NewTestFirestoreTokenFetcher(
	ctx context.Context,
	fsClient *firestore.Client,
	projectID string,
	logger zerolog.Logger,
) (cache.Fetcher[string, []routing.DeviceToken], error) {

	// 1. Create the generic fetcher that returns the full document.
	firestoreDocFetcher, err := cache.NewFirestore[string, deviceTokenDoc](
		ctx,
		&cache.FirestoreConfig{ProjectID: projectID, CollectionName: "device-tokens"},
		fsClient,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// 2. Wrap it in the adapter to match the required interface.
	return &FirestoreTokenAdapter{docFetcher: firestoreDocFetcher}, nil
}

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
