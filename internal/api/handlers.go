// Package api contains the HTTP handlers and other API-related components
// for the routing service.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/pkg/routing" // Updated import path
	"github.com/rs/zerolog"
)

// API holds the dependencies for the HTTP handlers, such as the logger
// and the message producer.
type API struct {
	// This now correctly refers to the public interface.
	producer routing.IngestionProducer
	logger   zerolog.Logger
}

// NewAPI is the constructor for the API struct.
func NewAPI(producer routing.IngestionProducer, logger zerolog.Logger) *API {
	api := API{
		producer: producer,
		logger:   logger,
	}
	return &api
}

// SendHandler is the HTTP handler for the POST /send endpoint.
// It decodes a SecureEnvelope from the request body and passes it to the
// IngestionProducer to be queued for processing.
func (a *API) SendHandler(writer http.ResponseWriter, request *http.Request) {
	var err error

	// Decode the JSON body into a SecureEnvelope struct.
	var envelope transport.SecureEnvelope
	err = json.NewDecoder(request.Body).Decode(&envelope)
	if err != nil {
		a.logger.Warn().Err(err).Msg("Failed to decode request body")
		http.Error(writer, "Bad Request: malformed JSON", http.StatusBadRequest)
		return
	}

	// Pass the envelope to the producer to be sent to the pipeline.
	err = a.producer.Publish(request.Context(), &envelope)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to publish envelope to message bus")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// On success, respond with 202 Accepted. This indicates that the request
	// has been accepted for processing, but is not yet complete.
	writer.WriteHeader(http.StatusAccepted)
}
