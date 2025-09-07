package routingservice

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/internal/api"
	"github.com/illmade-knight/routing-service/internal/pipeline"
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/rs/zerolog"
)

// Wrapper encapsulates all components of the running service.
type Wrapper struct {
	cfg               *routing.Config
	apiServer         *http.Server
	processingService *messagepipeline.StreamingService[transport.SecureEnvelope]
	logger            zerolog.Logger
}

// New creates and wires up the entire routing service.
func New(
	cfg *routing.Config,
	deps *routing.Dependencies,
	consumer messagepipeline.MessageConsumer, // Accepts the generic consumer
	producer routing.IngestionProducer, // Accepts the generic producer
	logger zerolog.Logger,
) (*Wrapper, error) {
	var err error

	pipelineConfig := pipeline.Config{
		NumWorkers: cfg.NumPipelineWorkers,
	}

	// The wrapper assembles the internal pipeline service...
	processingService, err := pipeline.NewService(pipelineConfig, deps, consumer, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create pipeline service: %w", err)
	}

	// ...and the internal API handlers.
	apiHandler := api.NewAPI(producer, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/send", apiHandler.SendHandler)
	apiServer := &http.Server{Addr: cfg.HTTPListenAddr, Handler: mux}

	wrapper := &Wrapper{
		cfg:               cfg,
		apiServer:         apiServer,
		processingService: processingService,
		logger:            logger,
	}
	return wrapper, nil
}

// Start runs the service's background components.
func (w *Wrapper) Start(ctx context.Context) error {
	w.logger.Info().Msg("Core processing pipeline starting.")
	if err := w.processingService.Start(ctx); err != nil {
		return fmt.Errorf("failed to start processing service: %w", err)
	}

	listener, err := net.Listen("tcp", w.cfg.HTTPListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on http address: %w", err)
	}
	w.apiServer.Addr = listener.Addr().String()

	go func() {
		w.logger.Info().Str("address", w.apiServer.Addr).Msg("HTTP server starting.")
		if err := w.apiServer.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			w.logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()
	return nil
}

// GetHTTPPort returns the resolved address the HTTP server is listening on.
func (w *Wrapper) GetHTTPPort() string {
	if w.apiServer != nil && w.apiServer.Addr != "" {
		_, port, err := net.SplitHostPort(w.apiServer.Addr)
		if err == nil {
			return ":" + port
		}
	}
	return ""
}

// Shutdown gracefully stops all service components.
func (w *Wrapper) Shutdown(ctx context.Context) error {
	w.logger.Info().Msg("Shutting down service components...")
	var finalErr error

	if err := w.apiServer.Shutdown(ctx); err != nil {
		w.logger.Error().Err(err).Msg("HTTP server shutdown failed.")
		finalErr = err
	}
	if err := w.processingService.Stop(ctx); err != nil {
		w.logger.Error().Err(err).Msg("Processing service shutdown failed.")
		finalErr = err
	}
	w.logger.Info().Msg("Service shutdown complete.")
	return finalErr
}
