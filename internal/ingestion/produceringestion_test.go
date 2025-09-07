// REFACTOR: This new file adds unit tests for the ingestion producer.

package ingestion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/v2/pstest"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/internal/ingestion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPubsubProducer_Publish(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	testEnvelope := &transport.SecureEnvelope{
		SenderID:    "user-alice",
		RecipientID: "user-bob",
	}

	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, err := grpc.NewClient(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// TODO: Handle error.
	}
	defer conn.Close()
	// Use the connection when creating a pubsub client.
	mockClient, err := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	require.NoError(t, err)

	topicName := "projects/project/topics/his"
	_, err = mockClient.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: topicName,
	})

	t.Run("Happy Path - Successful Publish", func(t *testing.T) {
		// Arrange
		publishCalled := false

		publisher := mockClient.Publisher("his")

		producer := ingestion.NewPubsubProducer(publisher)

		// Act
		err := producer.Publish(ctx, testEnvelope)

		// Assert
		require.NoError(t, err)
		assert.True(t, publishCalled, "Expected the underlying client's Publish method to be called")
	})

	t.Run("Failure - Underlying client returns an error", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("pubsub is unavailable")
		publisher := mockClient.Publisher("his")

		producer := ingestion.NewPubsubProducer(publisher)

		// Act
		err := producer.Publish(ctx, testEnvelope)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, expectedError, "The producer should wrap and return the client's error")
	})
}
