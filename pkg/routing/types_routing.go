// Package routing contains the public domain models, interfaces, and configuration
// for the routing service. It defines the contract for interacting with the service.
package routing

// ConnectionInfo holds real-time presence data for a connected user.
// This struct will be stored in Redis for fast lookups.
type ConnectionInfo struct {
	// ServerInstanceID is the unique identifier of the server instance
	// handling the user's persistent connection (e.g., a pod name in Kubernetes).
	ServerInstanceID string `json:"server_instance_id"`

	// Protocol indicates the type of connection the user has established.
	// Valid values are "websocket" or "mqtt".
	Protocol string `json:"protocol"`
}

// DeviceToken represents a single device token for push notifications.
// This struct will be stored in Firestore for persistent, long-term storage.
type DeviceToken struct {
	// Token is the actual push notification token provided by the client's OS.
	Token string `firestore:"token"`

	// Platform indicates the mobile operating system, e.g., "ios" or "android".
	Platform string `firestore:"platform"`
}
