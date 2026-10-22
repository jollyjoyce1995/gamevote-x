package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gamevote-api-go/api"
	"gamevote-api-go/internal/store"
)

// changed reloads the party after a mutation, pushes it to everybody watching
// and returns it. Pass wantResult=false to skip the reload if nobody listens.
func (s *Server) changed(ctx context.Context, code string, wantResult bool) (*api.Party, error) {
	if !wantResult && !s.broker.HasListeners(code) {
		return nil, nil
	}
	p, err := s.store.GetParty(ctx, code)
	if err != nil {
		return nil, err
	}
	s.broker.Publish(code, "party_updated", p)
	return p, nil
}

func (s *Server) ListParties(ctx context.Context, _ api.ListPartiesRequestObject) (api.ListPartiesResponseObject, error) {
	parties, err := s.store.ListParties(ctx)
	if err != nil {
		return nil, err
	}
	return api.ListParties200JSONResponse(parties), nil
}

func (s *Server) CreateParty(ctx context.Context, req api.CreatePartyRequestObject) (api.CreatePartyResponseObject, error) {
	var attendees []string
	if req.Body.Attendees != nil {
		for _, a := range *req.Body.Attendees {
			name, err := cleanUsername(a)
			if err != nil {
				return nil, err
			}
			attendees = append(attendees, name)
		}
	}
	code, err := s.store.CreateParty(ctx, attendees)
	if err != nil {
		return nil, err
	}
	p, err := s.store.GetParty(ctx, code)
	if err != nil {
		return nil, err
	}
	return api.CreateParty200JSONResponse(*p), nil
}

func (s *Server) GetParty(ctx context.Context, req api.GetPartyRequestObject) (api.GetPartyResponseObject, error) {
	p, err := s.store.GetParty(ctx, cleanCode(req.Code))
	if err != nil {
		return nil, err
	}
	return api.GetParty200JSONResponse(*p), nil
}

func (s *Server) SetPartyStatus(ctx context.Context, req api.SetPartyStatusRequestObject) (api.SetPartyStatusResponseObject, error) {
	if !req.Body.Status.Valid() {
		return nil, store.Invalid("unknown status %q", req.Body.Status)
	}
	code := cleanCode(req.Code)
	if err := s.store.SetStatus(ctx, code, req.Body.Status); err != nil {
		return nil, err
	}
	p, err := s.changed(ctx, code, true)
	if err != nil {
		return nil, err
	}
	return api.SetPartyStatus200JSONResponse(*p), nil
}

func (s *Server) JoinParty(ctx context.Context, req api.JoinPartyRequestObject) (api.JoinPartyResponseObject, error) {
	name, err := cleanUsername(req.Username)
	if err != nil {
		return nil, err
	}
	code := cleanCode(req.Code)
	if err := s.store.AddAttendee(ctx, code, name); err != nil {
		return nil, err
	}
	p, err := s.changed(ctx, code, true)
	if err != nil {
		return nil, err
	}
	return api.JoinParty200JSONResponse(*p), nil
}

func (s *Server) AddOption(ctx context.Context, req api.AddOptionRequestObject) (api.AddOptionResponseObject, error) {
	opt := *req.Body
	opt.Name = strings.TrimSpace(opt.Name)
	if opt.Name == "" || len(opt.Name) > 200 {
		return nil, store.Invalid("name must be 1-200 characters")
	}
	code := cleanCode(req.Code)
	if err := s.store.AddOption(ctx, code, opt); err != nil {
		return nil, err
	}
	if _, err := s.changed(ctx, code, false); err != nil {
		return nil, err
	}
	return api.AddOption204Response{}, nil
}

func (s *Server) RemoveOption(ctx context.Context, req api.RemoveOptionRequestObject) (api.RemoveOptionResponseObject, error) {
	code := cleanCode(req.Code)
	if err := s.store.RemoveOption(ctx, code, req.Name); err != nil {
		return nil, err
	}
	if _, err := s.changed(ctx, code, false); err != nil {
		return nil, err
	}
	return api.RemoveOption204Response{}, nil
}

func (s *Server) SubmitVotes(ctx context.Context, req api.SubmitVotesRequestObject) (api.SubmitVotesResponseObject, error) {
	name, err := cleanUsername(req.Username)
	if err != nil {
		return nil, err
	}
	code := cleanCode(req.Code)
	if err := s.store.SubmitVotes(ctx, code, name, *req.Body); err != nil {
		return nil, err
	}
	if _, err := s.changed(ctx, code, false); err != nil {
		return nil, err
	}
	return api.SubmitVotes204Response{}, nil
}

func (s *Server) AddBeer(ctx context.Context, req api.AddBeerRequestObject) (api.AddBeerResponseObject, error) {
	code := cleanCode(req.Code)
	drinkName := "Beer"
	if req.Body.DrinkName != nil && strings.TrimSpace(*req.Body.DrinkName) != "" {
		drinkName = strings.TrimSpace(*req.Body.DrinkName)
	}
	var beerEq float32 = 1.0
	if req.Body.BeerEquivalent != nil && *req.Body.BeerEquivalent > 0 {
		beerEq = *req.Body.BeerEquivalent
	}
	if err := s.store.AddBeer(ctx, code, strings.TrimSpace(req.Body.Attendee), drinkName, beerEq); err != nil {
		return nil, err
	}
	if _, err := s.changed(ctx, code, false); err != nil {
		return nil, err
	}
	return api.AddBeer204Response{}, nil
}

// ── live updates ────────────────────────────────────────────────────────

func (s *Server) StreamParty(ctx context.Context, req api.StreamPartyRequestObject) (api.StreamPartyResponseObject, error) {
	user, err := cleanUsername(req.Params.Username)
	if err != nil {
		return nil, err
	}
	code := cleanCode(req.Code)
	// Fail with a normal 404 before switching to streaming.
	if _, err := s.store.GetParty(ctx, code); err != nil {
		return nil, err
	}
	return sseStream{ctx: ctx, s: s, code: code, user: user}, nil
}

// sseStream is a response that stays open until the client goes away.
type sseStream struct {
	ctx  context.Context
	s    *Server
	code string
	user string
}

func (st sseStream) VisitStreamPartyResponse(w http.ResponseWriter) error {
	rc := http.NewResponseController(w)
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("X-Accel-Buffering", "no") // don't let nginx buffer the stream
	w.WriteHeader(http.StatusOK)

	// Join first so no update can fall between the snapshot and the subscription.
	c := st.s.broker.Join(st.code, st.user)
	defer st.s.broker.Leave(st.code, c)

	p, err := st.s.store.GetParty(st.ctx, st.code)
	if err != nil {
		return nil // headers are out already; the client will reconnect
	}
	if err := writeEvent(w, rc, "party_updated", p); err != nil {
		return nil
	}

	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()
	for {
		select {
		case msg := <-c.ch:
			if _, err := w.Write(msg); err != nil {
				return nil
			}
		case <-ping.C:
			if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
				return nil
			}
		case <-st.ctx.Done():
			return nil
		case <-st.s.broker.closed:
			return nil
		}
		if err := rc.Flush(); err != nil {
			return nil
		}
	}
}
