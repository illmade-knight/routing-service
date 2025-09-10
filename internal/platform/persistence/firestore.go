// Package persistence contains concrete storage implementations.
package persistence

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	usersCollection    = "user-messages"
	messagesCollection = "messages"
)

// FirestoreStore implements the routing.MessageStore interface using Google Cloud Firestore.
type FirestoreStore struct {
	client *firestore.Client
	logger zerolog.Logger
}

// NewFirestoreStore is the constructor for the FirestoreStore.
func NewFirestoreStore(client *firestore.Client, logger zerolog.Logger) (*FirestoreStore, error) {
	if client == nil {
		return nil, fmt.Errorf("firestore client cannot be nil")
	}
	store := &FirestoreStore{
		client: client,
		logger: logger,
	}
	return store, nil
}

// Store saves a message envelope for a specific user in Firestore.
func (s *FirestoreStore) Store(ctx context.Context, userID string, envelope *transport.SecureEnvelope) error {
	// REFACTOR: If the incoming envelope does not have a MessageID, generate one.
	// This makes the service more robust and removes the burden from the client,
	// as the ID is only used for internal storage management.
	if envelope.MessageID == "" {
		envelope.MessageID = uuid.NewString()
		s.logger.Warn().
			Str("newId", envelope.MessageID).
			Str("recipientID", userID).
			Msg("Generated new MessageID for incoming envelope.")
	}

	docRef := s.client.Collection(usersCollection).Doc(userID).Collection(messagesCollection).Doc(envelope.MessageID)
	_, err := docRef.Set(ctx, envelope)
	if err != nil {
		return fmt.Errorf("failed to store message %s for user %s: %w", envelope.MessageID, userID, err)
	}
	return nil
}

// FetchUndelivered retrieves all messages for a user that have not yet been delivered.
func (s *FirestoreStore) FetchUndelivered(ctx context.Context, userID string) ([]*transport.SecureEnvelope, error) {
	var envelopes []*transport.SecureEnvelope

	colRef := s.client.Collection(usersCollection).Doc(userID).Collection(messagesCollection)
	docs, err := colRef.Documents(ctx).GetAll()
	if err != nil {
		// It's not an error if the user's document or collection doesn't exist yet.
		if status.Code(err) == codes.NotFound {
			return envelopes, nil
		}
		return nil, fmt.Errorf("failed to fetch documents for user %s: %w", userID, err)
	}

	for _, doc := range docs {
		var envelope transport.SecureEnvelope
		err = doc.DataTo(&envelope)
		if err != nil {
			// Log the problematic document but continue processing others.
			s.logger.Error().
				Err(err).
				Str("documentID", doc.Ref.ID).
				Str("userID", userID).
				Msg("Failed to decode message from firestore")
			continue
		}
		envelopes = append(envelopes, &envelope)
	}

	return envelopes, nil
}

// MarkDelivered marks a set of messages as delivered by deleting them from Firestore.
func (s *FirestoreStore) MarkDelivered(ctx context.Context, userID string, envelopes []*transport.SecureEnvelope) error {
	if len(envelopes) == 0 {
		return nil
	}

	batch := s.client.Batch()
	colRef := s.client.Collection(usersCollection).Doc(userID).Collection(messagesCollection)

	for _, envelope := range envelopes {
		if envelope.MessageID == "" {
			s.logger.Warn().Str("userID", userID).Msg("Encountered envelope with no ID during MarkDelivered")
			continue
		}
		docRef := colRef.Doc(envelope.MessageID)
		batch.Delete(docRef)
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit batch delete for user %s: %w", userID, err)
	}

	s.logger.Info().
		Str("userID", userID).
		Int("count", len(envelopes)).
		Msg("Successfully marked messages as delivered")

	return nil
}
