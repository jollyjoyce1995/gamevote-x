// Package server implements the HTTP API defined in api/openapi.yaml.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"gamevote-api-go/api"
	"gamevote-api-go/internal/steam"
	"gamevote-api-go/internal/store"
)

type Server struct {
	store  *store.Store
	broker *Broker
	steam  *steam.Client
}

var _ api.StrictServerInterface = (*Server)(nil)

func New(st *store.Store, broker *Broker, sc *steam.Client) *Server {
	return &Server{store: st, broker: broker, steam: sc}
}

// Handler returns the API (under /api) plus, if staticDir is set, the built frontend.
func (s *Server) Handler(staticDir string) http.Handler {
	mux := http.NewServeMux()
	strict := api.NewStrictHandlerWithOptions(s, nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  func(w http.ResponseWriter, _ *http.Request, err error) { writeError(w, store.Invalid("%v", err)) },
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) { writeError(w, err) },
	})
	api.HandlerWithOptions(strict, api.StdHTTPServerOptions{
		BaseRouter:       mux,
		BaseURL:          "/api",
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) { writeError(w, store.Invalid("%v", err)) },
	})
	mux.HandleFunc("/api/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, store.NotFound("no such endpoint"))
	})
	if staticDir != "" {
		mux.Handle("/", staticHandler(staticDir))
	}
	return limitBodies(mux)
}

func limitBodies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		}
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, err error) {
	status, msg := http.StatusInternalServerError, "internal error"
	var se *store.Error
	if errors.As(err, &se) {
		msg = se.Msg
		switch se.Kind {
		case store.KindInvalid:
			status = http.StatusBadRequest
		case store.KindNotFound:
			status = http.StatusNotFound
		case store.KindConflict:
			status = http.StatusConflict
		case store.KindUnavailable:
			status = http.StatusBadGateway
		}
	} else {
		slog.Error("request failed", "error", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(api.Error{Error: msg})
}

// ── input cleaning ──────────────────────────────────────────────────────

func cleanUsername(s string) (string, error) {
	s = strings.TrimSpace(s)
	if n := utf8.RuneCountInString(s); n == 0 || n > 32 {
		return "", store.Invalid("username must be 1-32 characters")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", store.Invalid("username contains invalid characters")
		}
	}
	return s, nil
}

func cleanCode(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

// ── users ───────────────────────────────────────────────────────────────

func (s *Server) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	name, err := cleanUsername(req.Body.Username)
	if err != nil {
		return nil, err
	}
	u, err := s.store.Login(ctx, name)
	if err != nil {
		return nil, err
	}
	return api.Login200JSONResponse(u), nil
}

func (s *Server) ListUsers(ctx context.Context, _ api.ListUsersRequestObject) (api.ListUsersResponseObject, error) {
	users, err := s.store.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return api.ListUsers200JSONResponse(users), nil
}

// ── drinks ──────────────────────────────────────────────────────────────

func (s *Server) ListDrinks(ctx context.Context, _ api.ListDrinksRequestObject) (api.ListDrinksResponseObject, error) {
	drinks, err := s.store.ListDrinks(ctx)
	if err != nil {
		return nil, err
	}
	return api.ListDrinks200JSONResponse(drinks), nil
}

func (s *Server) CreateDrink(ctx context.Context, req api.CreateDrinkRequestObject) (api.CreateDrinkResponseObject, error) {
	d := *req.Body
	d.Name = strings.TrimSpace(d.Name)
	d.UnitName = strings.TrimSpace(d.UnitName)
	if d.Name == "" || d.UnitName == "" {
		return nil, store.Invalid("name and unitName are required")
	}
	if d.VolumeMl <= 0 || d.AlcoholPercent < 0 || d.AlcoholPercent > 100 {
		return nil, store.Invalid("volumeMl must be positive and alcoholPercent between 0 and 100")
	}
	// One 0.5l beer at 5% (25 ml of pure alcohol) is the reference.
	var beerEq float32
	if d.BeerEquivalent != nil && *d.BeerEquivalent > 0 {
		beerEq = *d.BeerEquivalent
	} else {
		beerEq = float32(d.VolumeMl) * d.AlcoholPercent / 100 / 25
	}
	created, err := s.store.CreateDrink(ctx, d, beerEq)
	if err != nil {
		return nil, err
	}
	return api.CreateDrink200JSONResponse(created), nil
}

// ── games ───────────────────────────────────────────────────────────────

func (s *Server) SearchGames(ctx context.Context, req api.SearchGamesRequestObject) (api.SearchGamesResponseObject, error) {
	q := strings.TrimSpace(req.Params.Q)
	if q == "" || len(q) > 100 {
		return nil, store.Invalid("q must be 1-100 characters")
	}
	games, err := s.steam.Search(ctx, q)
	if err != nil {
		return nil, err
	}
	return api.SearchGames200JSONResponse(games), nil
}
