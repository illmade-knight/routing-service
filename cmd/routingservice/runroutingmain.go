package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/routing-service/pkg/queue"
)

type api struct {
	queue queue.Queue
}

func (a *api) sendHandler(w http.ResponseWriter, r *http.Request) {
	var envelope transport.SecureEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err := a.queue.Enqueue(envelope); err != nil {
		http.Error(w, "Failed to enqueue message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	log.Printf("Accepted envelope from %s for %s", envelope.SenderID, envelope.RecipientID)
}

// ... you would also have a handler for fetching messages that uses a.queue.Dequeue ...

func main() {
	queue := queue.NewInMemoryQueue()
	app := &api{queue: queue}

	http.HandleFunc("/send", app.sendHandler)
	log.Println("Routing Service listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
