package pipeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/internal/pipeline" // Updated import
	"github.com/illmade-knight/routing-service/pkg/routing"       // Updated import
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mocks for Processor Dependencies ---

// mockFetcher is a test double for the generic cache.Fetcher interface.
type mockFetcher[K comparable, V any] struct {
	FetchFunc func(ctx context.Context, key K) (V, error)
}

func (m *mockFetcher[K, V]) Fetch(ctx context.Context, key K) (V, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, key)
	}
	var zeroValue V
	return zeroValue, errors.New("fetch function not implemented")
}

func (m *mockFetcher[K, V]) Close() error { return nil }

// mockDeliveryProducer is a test double for the routing.DeliveryProducer interface.
type mockDeliveryProducer struct {
	PublishFunc func(ctx context.Context, topicID string, data *transport.SecureEnvelope) error
}

func (m *mockDeliveryProducer) Publish(ctx context.Context, topicID string, data *transport.SecureEnvelope) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, topicID, data)
	}
	return errors.New("publish function not implemented")
}

// mockPushNotifier is a test double for the routing.PushNotifier interface.
type mockPushNotifier struct {
	NotifyFunc func(ctx context.Context, tokens []routing.DeviceToken, envelope *transport.SecureEnvelope) error
}

func (m *mockPushNotifier) Notify(ctx context.Context, tokens []routing.DeviceToken, envelope *transport.SecureEnvelope) error {
	if m.NotifyFunc != nil {
		return m.NotifyFunc(ctx, tokens, envelope)
	}
	return errors.New("notify function not implemented")
}

// --- Test Suite ---

func TestRoutingProcessor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	testEnvelope := &transport.SecureEnvelope{
		RecipientID:           "user-bob",
		EncryptedData:         []byte("test-data"),
		EncryptedSymmetricKey: []byte("test-key"),
	}
	testMessage := messagepipeline.Message{}

	t.Run("Happy Path - User is Online", func(t *testing.T) {
		// Arrange
		presenceCache := &mockFetcher[string, routing.ConnectionInfo]{
			FetchFunc: func(ctx context.Context, key string) (routing.ConnectionInfo, error) {
				assert.Equal(t, "user-bob", key)
				return routing.ConnectionInfo{ServerInstanceID: "pod-123"}, nil
			},
		}
		producerCalled := false
		deliveryProducer := &mockDeliveryProducer{
			PublishFunc: func(ctx context.Context, topicID string, data *transport.SecureEnvelope) error {
				producerCalled = true
				assert.Equal(t, "delivery-pod-123", topicID)
				assert.Equal(t, testEnvelope, data)
				return nil
			},
		}
		processor := pipeline.NewRoutingProcessor(presenceCache, nil, deliveryProducer, nil)

		// Act
		err := processor(ctx, testMessage, testEnvelope)

		// Assert
		require.NoError(t, err)
		assert.True(t, producerCalled)
	})

	t.Run("Happy Path - User is Offline with Device Tokens", func(t *testing.T) {
		// Arrange
		presenceCache := &mockFetcher[string, routing.ConnectionInfo]{
			FetchFunc: func(ctx context.Context, key string) (routing.ConnectionInfo, error) {
				return routing.ConnectionInfo{}, errors.New("not found")
			},
		}
		deviceTokenFetcher := &mockFetcher[string, []routing.DeviceToken]{
			FetchFunc: func(ctx context.Context, key string) ([]routing.DeviceToken, error) {
				assert.Equal(t, "user-bob", key)
				return []routing.DeviceToken{{Token: "device-abc"}}, nil
			},
		}
		notifierCalled := false
		pushNotifier := &mockPushNotifier{
			NotifyFunc: func(ctx context.Context, tokens []routing.DeviceToken, envelope *transport.SecureEnvelope) error {
				notifierCalled = true
				require.Len(t, tokens, 1)
				assert.Equal(t, "device-abc", tokens[0].Token)
				return nil
			},
		}
		processor := pipeline.NewRoutingProcessor(presenceCache, deviceTokenFetcher, nil, pushNotifier)

		// Act
		err := processor(ctx, testMessage, testEnvelope)

		// Assert
		require.NoError(t, err)
		assert.True(t, notifierCalled)
	})
}
