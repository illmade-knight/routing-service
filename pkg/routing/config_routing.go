package routing

import (
	"cloud.google.com/go/pubsub/v2"
	"github.com/illmade-knight/go-dataflow/pkg/cache"
	"google.golang.org/api/option"
)

// Config holds all the necessary configuration for the routing service.
type Config struct {
	ProjectID             string
	HTTPListenAddr        string
	IngressSubscriptionID string
	IngressTopicID        string
	NumPipelineWorkers    int
	PubsubClientOptions   []option.ClientOption
}

// Dependencies holds all the external components the service wrapper needs to run.
// These are provided as interfaces to decouple the wrapper from concrete implementations.
type Dependencies struct {
	PubsubClient       *pubsub.Client
	PresenceCache      cache.Fetcher[string, ConnectionInfo]
	DeviceTokenFetcher cache.Fetcher[string, []DeviceToken]
	DeliveryProducer   DeliveryProducer
	PushNotifier       PushNotifier
	MessageStore       MessageStore
}
