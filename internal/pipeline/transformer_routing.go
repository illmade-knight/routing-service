// Package pipeline contains the core data processing logic of the routing service.
package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/illmade-knight/go-dataflow/pkg/messagepipeline"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
)

// EnvelopeTransformer is a messagepipeline.MessageTransformer function. It decodes
// a raw message payload into a transport.SecureEnvelope struct.
//
// If the message payload is not valid JSON, it returns a wrapped error and a
// boolean flag indicating the message is malformed and should be skipped (not retried).
func EnvelopeTransformer(_ context.Context, msg *messagepipeline.Message) (*transport.SecureEnvelope, bool, error) {
	var envelope transport.SecureEnvelope
	var err error

	err = json.Unmarshal(msg.Payload, &envelope)
	if err != nil {
		// The incoming message payload is malformed. We wrap the error with context
		// and signal that the message should be skipped.
		wrappedErr := fmt.Errorf("failed to unmarshal secure envelope from message ID %s: %w", msg.ID, err)
		return nil, true, wrappedErr
	}

	// The message was successfully transformed.
	return &envelope, false, nil
}
