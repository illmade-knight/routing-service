package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/internal/api" // Updated import path
	"github.com/illmade-knight/routing-service/pkg/routing"  // Updated import path for interface
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockIngestionProducer is a test double for the routing.IngestionProducer interface.
type mockIngestionProducer struct {
	PublishFunc func(ctx context.Context, envelope *transport.SecureEnvelope) error
}

func (m *mockIngestionProducer) Publish(ctx context.Context, envelope *transport.SecureEnvelope) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, envelope)
	}
	return errors.New("publish function not implemented")
}

func TestSendHandler(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	validEnvelope := transport.SecureEnvelope{
		SenderID:    "user-alice",
		RecipientID: "user-bob",
	}
	validBody, err := json.Marshal(validEnvelope)
	require.NoError(t, err, "Setup: failed to marshal valid envelope")

	testCases := []struct {
		name                 string
		requestBody          []byte
		mockProducer         routing.IngestionProducer // Use the public interface type
		expectedStatusCode   int
		expectProducerToCall bool
	}{
		{
			name:        "Happy Path - Valid Request",
			requestBody: validBody,
			mockProducer: &mockIngestionProducer{
				PublishFunc: func(ctx context.Context, envelope *transport.SecureEnvelope) error {
					assert.Equal(t, validEnvelope.SenderID, envelope.SenderID)
					return nil
				},
			},
			expectedStatusCode:   http.StatusAccepted,
			expectProducerToCall: true,
		},
		{
			name:                 "Failure - Malformed JSON Body",
			requestBody:          []byte("{ not-json }"),
			mockProducer:         &mockIngestionProducer{},
			expectedStatusCode:   http.StatusBadRequest,
			expectProducerToCall: false,
		},
		{
			name:        "Failure - Producer Fails",
			requestBody: validBody,
			mockProducer: &mockIngestionProducer{
				PublishFunc: func(ctx context.Context, envelope *transport.SecureEnvelope) error {
					return errors.New("message bus is down")
				},
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectProducerToCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			request := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(tc.requestBody))
			responseRecorder := httptest.NewRecorder()
			producerWasCalled := false

			// We need to wrap the mock producer to track if it was called.
			// The mock must satisfy the public interface.
			wrappedProducer := &mockIngestionProducer{
				PublishFunc: func(ctx context.Context, envelope *transport.SecureEnvelope) error {
					producerWasCalled = true
					// We need to cast the test case's producer to call the underlying function
					if mock, ok := tc.mockProducer.(*mockIngestionProducer); ok {
						return mock.PublishFunc(ctx, envelope)
					}
					return errors.New("test setup error: mock producer not of correct type")
				},
			}

			apiHandler := api.NewAPI(wrappedProducer, zerolog.Nop())

			// Act
			apiHandler.SendHandler(responseRecorder, request)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, tc.expectProducerToCall, producerWasCalled)
		})
	}
}
