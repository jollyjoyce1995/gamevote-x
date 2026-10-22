package server_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"gamevote-api-go/api"
	"gamevote-api-go/internal/server"
	"gamevote-api-go/internal/steam"
	"gamevote-api-go/internal/store"
)

// Runs against a real Postgres: DATABASE_URL=postgres://... go test ./...
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	st, err := store.Open(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(st.Close)
	broker := server.NewBroker()
	ts := httptest.NewServer(server.New(st, broker, steam.New()).Handler(""))
	t.Cleanup(func() { broker.Close(); ts.Close() })
	return ts
}

func call[T any](t *testing.T, ts *httptest.Server, method, path string, body any, wantStatus int) T {
	t.Helper()
	var rd *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rd = bytes.NewReader(b)
	} else {
		rd = bytes.NewReader(nil)
	}
	req, _ := http.NewRequest(method, ts.URL+"/api"+path, rd)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out T
	if resp.StatusCode != wantStatus {
		var e api.Error
		_ = json.NewDecoder(resp.Body).Decode(&e)
		t.Fatalf("%s %s: status %d (%s), want %d", method, path, resp.StatusCode, e.Error, wantStatus)
	}
	if resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatalf("%s %s: decode: %v", method, path, err)
		}
	}
	return out
}

func TestPartyLifecycle(t *testing.T) {
	ts := newTestServer(t)
	suffix := time.Now().Format("150405.000")
	alice, bob := "alice-"+suffix, "bob-"+suffix

	call[api.User](t, ts, "POST", "/login", api.LoginRequest{Username: alice}, 200)
	call[api.User](t, ts, "POST", "/login", api.LoginRequest{Username: bob}, 200)
	call[api.Error](t, ts, "POST", "/login", api.LoginRequest{Username: "  "}, 400)

	// unknown attendee is rejected
	call[api.Error](t, ts, "POST", "/parties", api.CreatePartyRequest{Attendees: &[]string{"nobody-" + suffix}}, 400)

	p := call[api.Party](t, ts, "POST", "/parties", api.CreatePartyRequest{Attendees: &[]string{alice}}, 200)
	code := p.Code
	if len(code) != 6 || p.Status != api.NOMINATION || p.CurrentRound != 1 {
		t.Fatalf("unexpected new party: %+v", p)
	}

	// voting without options is refused, joining is idempotent
	call[api.Error](t, ts, "PUT", "/parties/"+code+"/status", api.SetStatusRequest{Status: api.VOTING}, 409)
	call[api.Party](t, ts, "PUT", "/parties/"+code+"/attendees/"+bob, nil, 200)
	p = call[api.Party](t, ts, "PUT", "/parties/"+code+"/attendees/"+bob, nil, 200)
	if len(p.Attendees) != 2 {
		t.Fatalf("attendees = %v", p.Attendees)
	}
	call[api.Error](t, ts, "PUT", "/parties/"+code+"/attendees/ghost-"+suffix, nil, 400)
	call[api.Error](t, ts, "PUT", "/parties/NOPE00/attendees/"+bob, nil, 404)

	// options: duplicates (case-insensitive) conflict, names with slashes work
	call[struct{}](t, ts, "POST", "/parties/"+code+"/options", api.PartyOption{Name: "Portal 2"}, 204)
	call[struct{}](t, ts, "POST", "/parties/"+code+"/options", api.PartyOption{Name: "portal 2"}, 409)
	call[struct{}](t, ts, "POST", "/parties/"+code+"/options", api.PartyOption{Name: "Tom/Jerry"}, 204)
	call[struct{}](t, ts, "POST", "/parties/"+code+"/options", api.PartyOption{Name: "Remove me"}, 204)
	call[struct{}](t, ts, "DELETE", "/parties/"+code+"/options/Remove%20me", nil, 204)
	call[api.Error](t, ts, "DELETE", "/parties/"+code+"/options/Remove%20me", nil, 404)
	call[struct{}](t, ts, "DELETE", "/parties/"+code+"/options/Tom%2FJerry", nil, 204)

	// beers
	call[struct{}](t, ts, "POST", "/parties/"+code+"/beers", api.AddBeerRequest{Attendee: alice}, 204)
	call[api.Error](t, ts, "POST", "/parties/"+code+"/beers", api.AddBeerRequest{Attendee: "ghost"}, 400)

	// voting
	p = call[api.Party](t, ts, "PUT", "/parties/"+code+"/status", api.SetStatusRequest{Status: api.VOTING}, 200)
	if len(p.PendingVoters) != 2 {
		t.Fatalf("pending = %v", p.PendingVoters)
	}
	call[api.Error](t, ts, "POST", "/parties/"+code+"/options", api.PartyOption{Name: "Late"}, 409)
	call[api.Error](t, ts, "PUT", "/parties/"+code+"/votes/"+alice, api.Votes{"Portal 2": 5}, 400)
	call[struct{}](t, ts, "PUT", "/parties/"+code+"/votes/"+alice, api.Votes{"Portal 2": -1}, 204)
	call[struct{}](t, ts, "PUT", "/parties/"+code+"/votes/"+alice, api.Votes{"Portal 2": 1}, 204) // re-vote replaces
	p = call[api.Party](t, ts, "GET", "/parties/"+code, nil, 200)
	if p.Status != api.VOTING || len(p.PendingVoters) != 1 || p.PendingVoters[0] != bob {
		t.Fatalf("after first vote: %+v", p)
	}
	call[struct{}](t, ts, "PUT", "/parties/"+code+"/votes/"+bob, api.Votes{"Portal 2": 1}, 204)
	p = call[api.Party](t, ts, "GET", "/parties/"+code, nil, 200)
	if p.Status != api.RESULTS || p.RoundResults["1"]["Portal 2"] != 2 {
		t.Fatalf("expected RESULTS with score 2, got %s %v", p.Status, p.RoundResults)
	}
	if p.BeerCount != 1 || p.BeerPerAttendee[alice] != 1 || p.BeerPerAttendee[bob] != 0 {
		t.Fatalf("beers: %+v", p)
	}

	// next round
	call[api.Error](t, ts, "PUT", "/parties/"+code+"/status", api.SetStatusRequest{Status: api.VOTING}, 409)
	p = call[api.Party](t, ts, "PUT", "/parties/"+code+"/status", api.SetStatusRequest{Status: api.NOMINATION}, 200)
	if p.CurrentRound != 2 || p.Status != api.NOMINATION || p.RoundResults["1"]["Portal 2"] != 2 {
		t.Fatalf("next round: %+v", p)
	}

	list := call[[]api.PartySummary](t, ts, "GET", "/parties", nil, 200)
	found := false
	for _, s := range list {
		if s.Code == code {
			found = true
			if s.AttendeeCount != 2 || s.OptionCount != 1 || s.BeerCount != 1 {
				t.Fatalf("summary: %+v", s)
			}
		}
	}
	if !found {
		t.Fatal("party missing from list")
	}
}

func TestDrinksAndUsers(t *testing.T) {
	ts := newTestServer(t)
	drinks := call[[]api.DrinkType](t, ts, "GET", "/drinks", nil, 200)
	if len(drinks) == 0 {
		t.Fatal("no seeded drinks")
	}
	name := "Test " + time.Now().Format("150405.000")
	d := call[api.DrinkType](t, ts, "POST", "/drinks",
		api.NewDrinkType{Name: name, VolumeMl: 500, AlcoholPercent: 10, UnitName: "x"}, 200)
	if d.BeerEquivalent != 2 {
		t.Fatalf("beer equivalent = %v, want 2", d.BeerEquivalent)
	}
	call[api.Error](t, ts, "POST", "/drinks",
		api.NewDrinkType{Name: name, VolumeMl: 500, AlcoholPercent: 10, UnitName: "x"}, 409)
	call[[]api.User](t, ts, "GET", "/users", nil, 200)
	call[api.Error](t, ts, "GET", "/nope", nil, 404)
}

func TestStream(t *testing.T) {
	ts := newTestServer(t)
	user := "streamer-" + time.Now().Format("150405.000")
	call[api.User](t, ts, "POST", "/login", api.LoginRequest{Username: user}, 200)
	p := call[api.Party](t, ts, "POST", "/parties", api.CreatePartyRequest{Attendees: &[]string{user}}, 200)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", ts.URL+"/api/parties/"+p.Code+"/stream?username="+user, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("content-type %q", ct)
	}

	events := make(chan string, 16)
	go func() {
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			if line := sc.Text(); strings.HasPrefix(line, "event: ") {
				events <- strings.TrimPrefix(line, "event: ")
			}
		}
	}()
	want := func(name string) {
		t.Helper()
		select {
		case got := <-events:
			if got != name {
				t.Fatalf("got event %q, want %q", got, name)
			}
		case <-ctx.Done():
			t.Fatalf("timed out waiting for %q", name)
		}
	}
	want("party_updated") // snapshot on connect
	want("online_users")

	call[struct{}](t, ts, "POST", "/parties/"+p.Code+"/options", api.PartyOption{Name: "Live"}, 204)
	want("party_updated")

	call[api.Error](t, ts, "GET", "/parties/NOPE00/stream?username="+user, nil, 404)
}
