// Package pipeline contains the core data processing logic of the routing service.
package pipeline

import (
	"context"
	"fmt"

	"github.com/illmade-knight/go-dataflow/pkg/cache"
	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/pkg/routing"
	"github.com/rs/zerolog"
)

// NewRoutingProcessor is a constructor that creates our main message handler.
// It accepts dependencies for caching, delivery, and offline storage, returning
// a configured messagepipeline.StreamProcessor.
func NewRoutingProcessor(
	presenceCache cache.Fetcher[string, routing.ConnectionInfo],
	deviceTokenFetcher cache.Fetcher[string, []routing.DeviceToken],
	deliveryProducer routing.DeliveryProducer,
	pushNotifier routing.PushNotifier,
	messageStore routing.MessageStore,
	logger zerolog.Logger,
) messagepipeline.StreamProcessor[transport.SecureEnvelope] {

	// The returned function is the StreamProcessor that will be executed by the pipeline.
	return func(ctx context.Context, original messagepipeline.Message, envelope *transport.SecureEnvelope) error {
		var err error

		// Strategy 1: Check if the user is online via the presence cache.
		connInfo, err := presenceCache.Fetch(ctx, envelope.RecipientID)
		if err == nil {
			// The user is online. Forward the envelope to the specific server
			// instance handling their connection.
			deliveryTopic := "delivery-" + connInfo.ServerInstanceID
			logger.Info().
				Str("recipientID", envelope.RecipientID).
				Str("deliveryTopic", deliveryTopic).
				Msg("User is online. Forwarding message.")

			err = deliveryProducer.Publish(ctx, deliveryTopic, envelope)
			if err != nil {
				return fmt.Errorf("failed to publish to delivery topic %s for user %s: %w", deliveryTopic, envelope.RecipientID, err)
			}
			return nil // Success
		}

		// Strategy 2: User is offline. Persist the message and then attempt a push notification.
		logger.Info().
			Str("recipientID", envelope.RecipientID).
			Msg("User is offline. Storing message and attempting push notification.")

		err = messageStore.Store(ctx, envelope.RecipientID, envelope)
		if err != nil {
			return fmt.Errorf("failed to store message for offline user %s: %w", envelope.RecipientID, err)
		}

		deviceTokens, err := deviceTokenFetcher.Fetch(ctx, envelope.RecipientID)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("recipientID", envelope.RecipientID).
				Msg("Could not retrieve device tokens for user. Message is stored.")
			return nil
		}

		err = pushNotifier.Notify(ctx, deviceTokens, envelope)
		if err != nil {
			return fmt.Errorf("failed to send push notification to user %s: %w", envelope.RecipientID, err)
		}

		return nil // Success
	}
}
