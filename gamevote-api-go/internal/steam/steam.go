// Package steam searches the public Steam store. It replaces syncing the whole
// app catalog into the database: we only ask Steam when somebody types a query.
package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gamevote-api-go/api"
	"gamevote-api-go/internal/store"
)

const (
	searchURL = "https://store.steampowered.com/api/storesearch/"
	cacheTTL  = 10 * time.Minute
	cacheMax  = 256
)

type Client struct {
	http *http.Client

	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	games   []api.Game
	expires time.Time
}

func New() *Client {
	return &Client{
		http:  &http.Client{Timeout: 5 * time.Second},
		cache: map[string]cacheEntry{},
	}
}

func (c *Client) Search(ctx context.Context, query string) ([]api.Game, error) {
	key := strings.ToLower(strings.TrimSpace(query))
	if games, ok := c.cached(key); ok {
		return games, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		searchURL+"?"+url.Values{"term": {key}, "l": {"english"}, "cc": {"us"}}.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, store.Unavailable("steam search failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, store.Unavailable("steam search returned %s", resp.Status)
	}

	var body struct {
		Items []struct {
			Type string `json:"type"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, store.Unavailable("steam search returned invalid data")
	}

	games := make([]api.Game, 0, len(body.Items))
	for _, it := range body.Items {
		if it.Type != "app" || it.Name == "" {
			continue
		}
		games = append(games, api.Game{
			AppId:    it.ID,
			Name:     it.Name,
			ImageUrl: fmt.Sprintf("https://cdn.akamai.steamstatic.com/steam/apps/%d/header.jpg", it.ID),
		})
	}
	c.store(key, games)
	return games, nil
}

func (c *Client) cached(key string) ([]api.Game, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.cache[key]
	if !ok || time.Now().After(e.expires) {
		delete(c.cache, key)
		return nil, false
	}
	return e.games, true
}

func (c *Client) store(key string, games []api.Game) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.cache) >= cacheMax {
		now := time.Now()
		for k, e := range c.cache {
			if now.After(e.expires) {
				delete(c.cache, k)
			}
		}
		if len(c.cache) >= cacheMax {
			clear(c.cache)
		}
	}
	c.cache[key] = cacheEntry{games: games, expires: time.Now().Add(cacheTTL)}
}
