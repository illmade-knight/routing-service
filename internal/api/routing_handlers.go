// Package api contains the HTTP handlers and other API-related components
// for the routing service.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/rs/zerolog"
)

// API holds the dependencies for the HTTP handlers, such as the logger,
// message producer, and message store.
type API struct {
	producer routing.IngestionProducer
	// REFACTOR: Add the MessageStore to handle fetching offline messages.
	store  routing.MessageStore
	logger zerolog.Logger
}

// NewAPI is the constructor for the API struct.
func NewAPI(producer routing.IngestionProducer, store routing.MessageStore, logger zerolog.Logger) *API {
	api := &API{
		producer: producer,
		store:    store,
		logger:   logger,
	}
	return api
}

// SendHandler is the HTTP handler for the POST /send endpoint.
// It decodes a SecureEnvelope from the request body and passes it to the
// IngestionProducer to be queued for processing.
func (a *API) SendHandler(writer http.ResponseWriter, request *http.Request) {
	var err error

	var envelope transport.SecureEnvelope
	err = json.NewDecoder(request.Body).Decode(&envelope)
	if err != nil {
		a.logger.Warn().Err(err).Msg("Failed to decode request body")
		http.Error(writer, "Bad Request: malformed JSON", http.StatusBadRequest)
		return
	}

	err = a.producer.Publish(request.Context(), &envelope)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to publish envelope to message bus")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

// REFACTOR: Add the handler for retrieving offline messages.

// GetMessagesHandler is the HTTP handler for the GET /messages endpoint.
// It fetches all stored messages for an authenticated user, marks them as
// delivered, and returns them to the client.
func (a *API) GetMessagesHandler(writer http.ResponseWriter, request *http.Request) {
	// NOTE: In a real application, the userID would be extracted from a JWT
	// or session token after an authentication middleware has run.
	// For now, we'll read it from a header as a placeholder.
	userID := request.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(writer, "Unauthorized: Missing X-User-ID header", http.StatusUnauthorized)
		return
	}

	// 1. Fetch all undelivered messages from the store.
	messages, err := a.store.FetchUndelivered(request.Context(), userID)
	if err != nil {
		a.logger.Error().Err(err).Str("userID", userID).Msg("Failed to fetch undelivered messages")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 2. Respond to the client with the messages.
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(messages)
	if err != nil {
		a.logger.Error().Err(err).Str("userID", userID).Msg("Failed to encode messages to response")
		// Can't send http.Error here as headers/status are already written.
		return
	}

	// 3. After successful delivery, mark the messages as delivered in the store.
	// This prevents them from being sent again.
	if len(messages) > 0 {
		err = a.store.MarkDelivered(request.Context(), userID, messages)
		if err != nil {
			// This is a problematic state: the client has the messages, but we
			// failed to mark them as delivered. They will receive duplicates next time.
			// A robust system might use a two-phase commit or transactional outbox
			// pattern here, but for now, we just log this critical failure.
			a.logger.Error().Err(err).Str("userID", userID).Int("count", len(messages)).
				Msg("CRITICAL: Failed to mark messages as delivered after sending to client")
		}
	}
}
