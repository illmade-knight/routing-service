package queue

import (
	"sync"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
)

// InMemoryQueue is a thread-safe in-memory message queue.
type InMemoryQueue struct {
	sync.RWMutex
	messages map[string][]transport.SecureEnvelope
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{messages: make(map[string][]transport.SecureEnvelope)}
}

func (q *InMemoryQueue) Enqueue(envelope transport.SecureEnvelope) error {
	q.Lock()
	defer q.Unlock()
	q.messages[envelope.RecipientID] = append(q.messages[envelope.RecipientID], envelope)
	return nil
}

func (q *InMemoryQueue) Dequeue(userID string) ([]transport.SecureEnvelope, error) {
	q.Lock()
	defer q.Unlock()
	envelopes := q.messages[userID]
	delete(q.messages, userID) // Clear the queue for the user after retrieval
	return envelopes, nil
}
