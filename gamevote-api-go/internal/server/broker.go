package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"sync"
)

// Broker fans server-sent events out to the clients watching each party.
// Single process, in memory: all we need for one small instance.
type Broker struct {
	mu     sync.Mutex
	rooms  map[string]map[*client]struct{}
	closed chan struct{}
	once   sync.Once
}

type client struct {
	user string
	ch   chan []byte
}

func NewBroker() *Broker {
	return &Broker{rooms: map[string]map[*client]struct{}{}, closed: make(chan struct{})}
}

// Join registers a client and tells the room who is online.
func (b *Broker) Join(code, user string) *client {
	c := &client{user: user, ch: make(chan []byte, 8)}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.rooms[code] == nil {
		b.rooms[code] = map[*client]struct{}{}
	}
	b.rooms[code][c] = struct{}{}
	b.sendLocked(code, "online_users", b.onlineLocked(code))
	return c
}

func (b *Broker) Leave(code string, c *client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	room := b.rooms[code]
	delete(room, c)
	if len(room) == 0 {
		delete(b.rooms, code)
		return
	}
	b.sendLocked(code, "online_users", b.onlineLocked(code))
}

func (b *Broker) HasListeners(code string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.rooms[code]) > 0
}

func (b *Broker) Publish(code, event string, payload any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sendLocked(code, event, payload)
}

// Close ends all streams (used on shutdown).
func (b *Broker) Close() { b.once.Do(func() { close(b.closed) }) }

func (b *Broker) onlineLocked(code string) []string {
	users := []string{}
	for c := range b.rooms[code] {
		if !slices.Contains(users, c.user) {
			users = append(users, c.user)
		}
	}
	slices.Sort(users)
	return users
}

func (b *Broker) sendLocked(code, event string, payload any) {
	room := b.rooms[code]
	if len(room) == 0 {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("sse: marshal failed", "event", event, "error", err)
		return
	}
	msg := []byte("event: " + event + "\ndata: " + string(data) + "\n\n")
	for c := range room {
		select {
		case c.ch <- msg:
		default:
			// Slow client: drop its oldest message. Every event carries the
			// full state, so the newest one is the one that matters.
			select {
			case <-c.ch:
			default:
			}
			select {
			case c.ch <- msg:
			default:
			}
		}
	}
}

func formatEvent(event string, payload any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return []byte("event: " + event + "\ndata: " + string(data) + "\n\n"), nil
}

// writeEvent sends one event straight to a single client.
func writeEvent(w http.ResponseWriter, rc *http.ResponseController, event string, payload any) error {
	msg, err := formatEvent(event, payload)
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return rc.Flush()
}
