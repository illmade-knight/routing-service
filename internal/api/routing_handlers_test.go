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
	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/illmade-knight/routing-service/internal/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks using testify/mock ---

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

// REFACTOR: Implement the new routing.MessageStore interface.
func (m *mockMessageStore) StoreMessages(ctx context.Context, recipient urn.URN, envelopes []*transport.SecureEnvelope) error {
	args := m.Called(ctx, recipient, envelopes)
	return args.Error(0)
}

func (m *mockMessageStore) RetrieveMessages(ctx context.Context, recipient urn.URN) ([]*transport.SecureEnvelope, error) {
	args := m.Called(ctx, recipient)
	var result []*transport.SecureEnvelope
	if val, ok := args.Get(0).([]*transport.SecureEnvelope); ok {
		result = val
	}
	return result, args.Error(1)
}

func (m *mockMessageStore) DeleteMessages(ctx context.Context, recipient urn.URN, messageIDs []string) error {
	args := m.Called(ctx, recipient, messageIDs)
	return args.Error(0)
}

// --- Test Suites ---

func TestSendHandler(t *testing.T) {
	// REFACTOR: Use a valid URN for the test envelope.
	testURN, err := urn.Parse("urn:sm:user:user-bob")
	require.NoError(t, err)

	validEnvelope := transport.SecureEnvelope{
		SenderID:    testURN, // Sender can also be a URN
		RecipientID: testURN,
	}
	validBody, err := json.Marshal(validEnvelope)
	require.NoError(t, err, "Setup: failed to marshal valid envelope")

	// REFACTOR: Test backward compatibility with a legacy userID.
	legacyEnvelope := map[string]string{
		"senderId":    "user-alice",
		"recipientId": "user-bob",
	}
	legacyBody, err := json.Marshal(legacyEnvelope)
	require.NoError(t, err, "Setup: failed to marshal legacy envelope")

	testCases := []struct {
		name               string
		requestBody        []byte
		setupMock          func(producer *mockIngestionProducer)
		expectedStatusCode int
	}{
		{
			name:        "Happy Path - Valid URN Request",
			requestBody: validBody,
			setupMock: func(producer *mockIngestionProducer) {
				producer.On("Publish", mock.Anything, mock.AnythingOfType("*transport.SecureEnvelope")).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:        "Happy Path - Legacy UserID Request",
			requestBody: legacyBody,
			setupMock: func(producer *mockIngestionProducer) {
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
	testURN, err := urn.Parse("urn:sm:user:user-bob")
	require.NoError(t, err)

	testMessages := []*transport.SecureEnvelope{
		{MessageID: "msg-1", RecipientID: testURN},
		{MessageID: "msg-2", RecipientID: testURN},
	}

	testCases := []struct {
		name               string
		userIDHeader       string
		setupMock          func(store *mockMessageStore)
		expectedStatusCode int
		expectedBody       string // Optional, for checking response body
	}{
		{
			name:         "Happy Path - Messages Found (URN Header)",
			userIDHeader: "urn:sm:user:user-bob",
			setupMock: func(store *mockMessageStore) {
				store.On("RetrieveMessages", mock.Anything, testURN).Return(testMessages, nil).Once()
				store.On("DeleteMessages", mock.Anything, testURN, []string{"msg-1", "msg-2"}).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:         "Happy Path - Messages Found (Legacy Header)",
			userIDHeader: "user-bob",
			setupMock: func(store *mockMessageStore) {
				store.On("RetrieveMessages", mock.Anything, testURN).Return(testMessages, nil).Once()
				store.On("DeleteMessages", mock.Anything, testURN, []string{"msg-1", "msg-2"}).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:         "Happy Path - No Messages Found",
			userIDHeader: "user-bob",
			setupMock: func(store *mockMessageStore) {
				store.On("RetrieveMessages", mock.Anything, testURN).Return([]*transport.SecureEnvelope{}, nil).Once()
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "Failure - Missing Auth Header",
			userIDHeader:       "",
			setupMock:          func(store *mockMessageStore) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:         "Failure - Store Fails on Fetch",
			userIDHeader: "user-bob",
			setupMock: func(store *mockMessageStore) {
				store.On("RetrieveMessages", mock.Anything, testURN).Return([]*transport.SecureEnvelope(nil), errors.New("db is down")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
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
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, responseRecorder.Body.String())
			}
			store.AssertExpectations(t)
		})
	}
}
