// REFACTOR: The main package is now a thin wrapper that uses the /pkg/routing
// component to configure and run the service. It is responsible for constructing
// the concrete dependencies and injecting them into the service wrapper.

// The server command is responsible for initializing and running the routing service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub/v2"
	"github.com/illmade-knight/go-dataflow/pkg/cache"
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/illmade-knight/routing-service/routingservice"
	"github.com/rs/zerolog"
)

func main() {
	// In a production application, we would use a structured logger with
	// configurable levels and output.
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Initializing routing service...")

	// --- 1. Load Configuration ---
	// In a real application, these flags would be populated from a config file,
	// environment variables, or command-line flags to allow for different
	// configurations in different environments (dev, staging, prod).
	useRedis := true
	useFirestore := true

	cfg := &routing.Config{
		ProjectID:             "your-gcp-project-id",
		IngressSubscriptionID: "envelope-ingress-sub",
		IngressTopicID:        "envelope-ingress",
		HTTPListenAddr:        ":8080",
		NumPipelineWorkers:    50,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- 2. Initialize External Clients ---
	// These clients are created here so their lifecycle (including shutdown) can be
	// managed by the main package.
	pubsubClient, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	defer pubsubClient.Close()

	// --- 3. Construct Concrete Storage Dependencies based on Config ---
	var presenceCache cache.Fetcher[string, routing.ConnectionInfo]
	var deviceTokenFetcher cache.Fetcher[string, []routing.DeviceToken]

	// This block demonstrates the flexibility of the new pattern. We can
	// conditionally construct the cache based on configuration.
	if useRedis {
		redisConfig := &cache.RedisConfig{Addr: "localhost:6379", CacheTTL: 5 * time.Minute}
		presenceCache, err = cache.NewRedisCache[string, routing.ConnectionInfo](ctx, redisConfig, logger, nil)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create presence cache")
		}
		defer presenceCache.Close()
	} else {
		// For local development or testing, we can use a simple in-memory cache
		// without needing a Redis instance.
		presenceCache = cache.NewInMemoryCache[string, routing.ConnectionInfo](nil)
	}

	if useFirestore {
		firestoreClient, firestoreErr := firestore.NewClient(ctx, cfg.ProjectID)
		if firestoreErr != nil {
			logger.Fatal().Err(firestoreErr).Msg("Failed to create firestore client")
		}
		defer firestoreClient.Close()

		firestoreConfig := &cache.FirestoreConfig{ProjectID: cfg.ProjectID, CollectionName: "device-tokens"}
		deviceTokenFetcher, err = cache.NewFirestore[string, []routing.DeviceToken](ctx, firestoreConfig, firestoreClient, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create device token fetcher")
		}
	} else {
		// For testing, we can use an in-memory fetcher that returns empty results.
		deviceTokenFetcher = cache.NewInMemoryCache[string, []routing.DeviceToken](nil)
	}

	// --- 4. Create and Run the Service Wrapper ---
	// Assemble the dependencies struct and inject it into the service wrapper.
	deps := &routing.Dependencies{
		PresenceCache:      presenceCache,
		DeviceTokenFetcher: deviceTokenFetcher,
	}
	service, err := routingservice.NewRoutingServiceWrapper(cfg, deps)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create routing service wrapper")
	}

	err = service.Start(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start service")
	}

	// --- 5. Handle Graceful Shutdown ---
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan // Block until a shutdown signal is received.

	log.Println("Shutdown signal received. Starting graceful shutdown...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := service.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Service shutdown failed.")
	}
}
