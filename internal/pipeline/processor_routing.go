// Package pipeline contains the core data processing logic of the routing service.
package pipeline

import (
	"context"
	"fmt"
	"log" // NOTE: In production, this would be a structured logger like zerolog.

	"github.com/illmade-knight/go-dataflow/pkg/cache"
	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/pkg/routing" // Updated import
)

// NewRoutingProcessor is a constructor that creates our main message handler.
// It accepts dependencies for presence caching, device token fetching,
// and message delivery, returning a configured messagepipeline.StreamProcessor.
func NewRoutingProcessor(
	presenceCache cache.Fetcher[string, routing.ConnectionInfo],
	deviceTokenFetcher cache.Fetcher[string, []routing.DeviceToken],
	deliveryProducer routing.DeliveryProducer,
	pushNotifier routing.PushNotifier,
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
			log.Printf("User %s is online. Forwarding to topic %s", envelope.RecipientID, deliveryTopic)

			err = deliveryProducer.Publish(ctx, deliveryTopic, envelope)
			if err != nil {
				// If publishing fails, return the error.
				// This will cause the message to be Nacked and likely retried.
				return fmt.Errorf("failed to publish to delivery topic %s for user %s: %w", deliveryTopic, envelope.RecipientID, err)
			}
			return nil // Success
		}

		// Strategy 2: User is offline. Fetch their device tokens for a push notification.
		log.Printf("User %s is offline. Attempting to send push notification.", envelope.RecipientID)
		deviceTokens, err := deviceTokenFetcher.Fetch(ctx, envelope.RecipientID)
		if err != nil {
			// We could not find device tokens. There is nothing more to do.
			// Log the error but return nil, as retrying will not help.
			// This effectively Acks the message.
			log.Printf("Could not retrieve device tokens for user %s: %v. Dropping message.", envelope.RecipientID, err)
			return nil
		}

		err = pushNotifier.Notify(ctx, deviceTokens, envelope)
		if err != nil {
			// The push notification failed. This could be a transient error, so we
			// return the error to the pipeline for a potential retry.
			return fmt.Errorf("failed to send push notification to user %s: %w", envelope.RecipientID, err)
		}

		return nil // Success
	}
}
