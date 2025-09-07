package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/illmade-knight/go-dataflow/pkg/cache"
	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	psadapter "github.com/illmade-knight/routing-service/internal/platform/pubsub" // Adapter
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/illmade-knight/routing-service/routingservice" // Service Library
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 1. Load configuration from environment variables, flags, or a file
	cfg := &routing.Config{
		ProjectID:             os.Getenv("GCP_PROJECT_ID"),
		HTTPListenAddr:        ":8080",
		IngressSubscriptionID: "ingress-sub",
		IngressTopicID:        "ingress-topic",
		NumPipelineWorkers:    10,
	}

	// 2. Initialize REAL external clients
	psClient, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
	defer psClient.Close()

	// 3. Instantiate CONCRETE adapters using the real clients
	//    Here we choose our Google Cloud implementations.
	consumerConfig := messagepipeline.NewGooglePubsubConsumerDefaults(cfg.IngressSubscriptionID)
	consumer, err := messagepipeline.NewGooglePubsubConsumer(consumerConfig, psClient, logger)
	if err != nil {
		log.Fatalf("Failed to create pubsub consumer: %v", err)
	}

	ingestionTopic := psClient.Publisher(cfg.IngressTopicID)
	ingestionProducer := psadapter.NewProducer(ingestionTopic)

	// 4. Create other real dependencies (using nil/mocks for now)
	deps := &routing.Dependencies{
		// In a real application, you would initialize real Redis, Firestore, etc. clients here.
		PresenceCache:      cache.NewInMemoryCache[string, routing.ConnectionInfo](nil),
		DeviceTokenFetcher: cache.NewInMemoryCache[string, []routing.DeviceToken](nil),
		DeliveryProducer:   &mockDeliveryProducer{}, // Placeholder
		PushNotifier:       &mockPushNotifier{},     // Placeholder
	}

	// 5. Create the service using the public wrapper, injecting all dependencies
	service, err := routingservice.New(cfg, deps, consumer, ingestionProducer, logger)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	// 6. Start the service and handle graceful shutdown
	go func() {
		if err := service.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("Service failed to start")
			stop()
		}
	}()

	logger.Info().Msg("Routing service started. Press Ctrl+C to shut down.")
	<-ctx.Done() // Wait for shutdown signal

	logger.Info().Msg("Shutdown signal received. Gracefully stopping service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := service.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().Err(err).Msg("Service shutdown failed.")
	}

	logger.Info().Msg("Service stopped gracefully.")
}

// mockDeliveryProducer and mockPushNotifier are placeholders for the main executable.
// In a real scenario, these would be concrete implementations (e.g., another Pub/Sub producer).
type mockDeliveryProducer struct{}

func (m *mockDeliveryProducer) Publish(ctx context.Context, topicID string, data *transport.SecureEnvelope) error {
	log.Printf("MOCK: Publishing to delivery topic %s for user %s", topicID, data.RecipientID)
	return nil
}

type mockPushNotifier struct{}

func (m *mockPushNotifier) Notify(ctx context.Context, tokens []routing.DeviceToken, envelope *transport.SecureEnvelope) error {
	log.Printf("MOCK: Sending push notification to %d devices for user %s", len(tokens), envelope.RecipientID)
	return nil
}
