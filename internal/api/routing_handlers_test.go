package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/internal/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockIngestionProducer struct {
	mock.Mock
}

func (m *mockIngestionProducer) Publish(ctx context.Context, envelope *transport.SecureEnvelope) error {
	args := m.Called(ctx, envelope)
	return args.Error(0)
}

type mockMessageStore struct {
	mock.Mock
}

func (m *mockMessageStore) Store(ctx context.Context, userID string, envelope *transport.SecureEnvelope) error {
	args := m.Called(ctx, userID, envelope)
	return args.Error(0)
}
func (m *mockMessageStore) FetchUndelivered(ctx context.Context, userID string) ([]*transport.SecureEnvelope, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*transport.SecureEnvelope), args.Error(1)
}
func (m *mockMessageStore) MarkDelivered(ctx context.Context, userID string, envelopes []*transport.SecureEnvelope) error {
	args := m.Called(ctx, userID, envelopes)
	return args.Error(0)
}

func TestSendHandler(t *testing.T) {
	validEnvelope := transport.SecureEnvelope{
		SenderID:    "user-alice",
		RecipientID: "user-bob",
	}
	validBody, err := json.Marshal(validEnvelope)
	require.NoError(t, err, "Setup: failed to marshal valid envelope")

	testCases := []struct {
		name               string
		requestBody        []byte
		setupMock          func(producer *mockIngestionProducer)
		expectedStatusCode int
	}{
		{
			name:        "Happy Path - Valid Request",
			requestBody: validBody,
			setupMock: func(producer *mockIngestionProducer) {
				// Use mock.AnythingOfType because the pointer will be different after unmarshaling.
				producer.On("Publish", mock.Anything, mock.AnythingOfType("*transport.SecureEnvelope")).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:               "Failure - Malformed JSON Body",
			requestBody:        []byte("{ not-json }"),
			setupMock:          func(producer *mockIngestionProducer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Failure - Producer Fails",
			requestBody: validBody,
			setupMock: func(producer *mockIngestionProducer) {
				producer.On("Publish", mock.Anything, mock.AnythingOfType("*transport.SecureEnvelope")).Return(errors.New("message bus is down"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			producer := new(mockIngestionProducer)
			store := new(mockMessageStore)
			tc.setupMock(producer)

			apiHandler := api.NewAPI(producer, store, zerolog.Nop())
			request := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(tc.requestBody))
			responseRecorder := httptest.NewRecorder()

			// Act
			apiHandler.SendHandler(responseRecorder, request)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, responseRecorder.Code)
			producer.AssertExpectations(t)
		})
	}
}

func TestGetMessagesHandler(t *testing.T) {
	testUserID := "user-bob"
	testMessages := []*transport.SecureEnvelope{
		{MessageID: "msg-1", RecipientID: testUserID},
		{MessageID: "msg-2", RecipientID: testUserID},
	}

	testCases := []struct {
		name               string
		userIDHeader       string
		setupMock          func(store *mockMessageStore)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:         "Happy Path - Messages Found",
			userIDHeader: testUserID,
			setupMock: func(store *mockMessageStore) {
				store.On("FetchUndelivered", mock.Anything, testUserID).Return(testMessages, nil).Once()
				store.On("MarkDelivered", mock.Anything, testUserID, testMessages).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			// REFACTOR: Update expected JSON to use camelCase keys.
			expectedBody: `[{"senderId":"","recipientId":"user-bob","messageId":"msg-1","encryptedData":null,"encryptedSymmetricKey":null,"signature":null},{"senderId":"","recipientId":"user-bob","messageId":"msg-2","encryptedData":null,"encryptedSymmetricKey":null,"signature":null}]` + "\n",
		},
		{
			name:         "Happy Path - No Messages Found",
			userIDHeader: testUserID,
			setupMock: func(store *mockMessageStore) {
				store.On("FetchUndelivered", mock.Anything, testUserID).Return([]*transport.SecureEnvelope{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[]` + "\n",
		},
		{
			name:               "Failure - Missing Auth Header",
			userIDHeader:       "",
			setupMock:          func(store *mockMessageStore) {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `Unauthorized: Missing X-User-ID header` + "\n",
		},
		{
			name:         "Failure - Store Fails on Fetch",
			userIDHeader: testUserID,
			setupMock: func(store *mockMessageStore) {
				store.On("FetchUndelivered", mock.Anything, testUserID).Return([]*transport.SecureEnvelope(nil), errors.New("db is down")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `Internal Server Error` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			store := new(mockMessageStore)
			producer := new(mockIngestionProducer)
			tc.setupMock(store)

			apiHandler := api.NewAPI(producer, store, zerolog.Nop())
			request := httptest.NewRequest(http.MethodGet, "/messages", nil)
			if tc.userIDHeader != "" {
				request.Header.Set("X-User-ID", tc.userIDHeader)
			}
			responseRecorder := httptest.NewRecorder()

			// Act
			apiHandler.GetMessagesHandler(responseRecorder, request)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, tc.expectedBody, responseRecorder.Body.String())
			store.AssertExpectations(t)
		})
	}
}
