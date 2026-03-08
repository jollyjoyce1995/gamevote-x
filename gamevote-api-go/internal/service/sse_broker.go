package service

import (
	"encoding/json"
	"log/slog"
	"sync"
)

// SSEClient represents a connected SSE client
type SSEClient struct {
	ID      string
	Channel chan string
}

// SSEBroker manages SSE connections keyed by party code
type SSEBroker struct {
	mu      sync.RWMutex
	clients map[string]map[string]*SSEClient // partyCode -> clientID -> client
}

var Broker = &SSEBroker{
	clients: make(map[string]map[string]*SSEClient),
}

// Register a new client for a party
func (b *SSEBroker) Register(partyCode string, clientID string) *SSEClient {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.clients[partyCode]; !ok {
		b.clients[partyCode] = make(map[string]*SSEClient)
	}

	client := &SSEClient{
		ID:      clientID,
		Channel: make(chan string, 10),
	}
	b.clients[partyCode][clientID] = client
	slog.Info("SSE: Client connected", "client_id", clientID, "party_code", partyCode)
	return client
}

// Unregister a client from a party
func (b *SSEBroker) Unregister(partyCode string, clientID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if clients, ok := b.clients[partyCode]; ok {
		if client, ok := clients[clientID]; ok {
			close(client.Channel)
			delete(clients, clientID)
			slog.Info("SSE: Client disconnected", "client_id", clientID, "party_code", partyCode)
		}
		if len(clients) == 0 {
			delete(b.clients, partyCode)
		}
	}
}

// Broadcast sends an event to all clients listening to a party
func (b *SSEBroker) Broadcast(partyCode string, eventType string, data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	clients, ok := b.clients[partyCode]
	if !ok {
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("SSE: Failed to marshal broadcast data", "error", err)
		return
	}

	msg := "event: " + eventType + "\ndata: " + string(jsonData) + "\n\n"
	for _, client := range clients {
		select {
		case client.Channel <- msg:
		default:
			slog.Warn("SSE: Client buffer full, skipping", "client_id", client.ID)
		}
	}
}

// OnlineUsers returns the number of clients connected to a particular party
func (b *SSEBroker) OnlineUsers(partyCode string) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	clients, ok := b.clients[partyCode]
	if !ok {
		return []string{}
	}

	names := make([]string, 0, len(clients))
	for id := range clients {
		names = append(names, id)
	}
	return names
}
