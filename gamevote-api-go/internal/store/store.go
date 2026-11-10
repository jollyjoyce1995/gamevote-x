// Package store is the Postgres access layer.
package store

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gamevote-api-go/api"
)

//go:embed schema.sql
var schema string

type Store struct {
	pool *pgxpool.Pool
}

// Open connects (retrying while the database starts up) and applies the schema.
// The pool is deliberately tiny: a handful of users never needs more.
func Open(ctx context.Context, url string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	cfg.MaxConns = 4
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		if err = pool.Ping(ctx); err == nil {
			break
		}
		if time.Now().After(deadline) || ctx.Err() != nil {
			pool.Close()
			return nil, fmt.Errorf("database not reachable: %w", err)
		}
		slog.Info("waiting for database", "error", err)
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			pool.Close()
			return nil, ctx.Err()
		}
	}

	if err := applySchema(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return &Store{pool: pool}, nil
}

func applySchema(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	// The simple protocol allows several statements in one call.
	_, err = conn.Conn().PgConn().Exec(ctx, schema).ReadAll()
	return err
}

func (s *Store) Close() { s.pool.Close() }

// querier is satisfied by both *pgxpool.Pool and pgx.Tx.
type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func pgCode(err error) string {
	var pe *pgconn.PgError
	if errors.As(err, &pe) {
		return pe.Code
	}
	return ""
}

const (
	codeUniqueViolation = "23505"
	codeForeignKey      = "23503"
)

// ── users ───────────────────────────────────────────────────────────────

func (s *Store) Login(ctx context.Context, username string) (api.User, error) {
	var u api.User
	err := s.pool.QueryRow(ctx, `
		INSERT INTO users (username) VALUES ($1)
		ON CONFLICT (username) DO UPDATE SET last_login = now()
		RETURNING username, created_at, last_login`, username).
		Scan(&u.Username, &u.CreatedAt, &u.LastLogin)
	return u, err
}

func (s *Store) ListUsers(ctx context.Context) ([]api.User, error) {
	rows, err := s.pool.Query(ctx, `SELECT username, created_at, last_login FROM users ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []api.User{}
	for rows.Next() {
		var u api.User
		if err := rows.Scan(&u.Username, &u.CreatedAt, &u.LastLogin); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// ── drinks ──────────────────────────────────────────────────────────────

func (s *Store) ListDrinks(ctx context.Context) ([]api.DrinkType, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, volume_ml, alcohol_percent, beer_equivalent, unit_name
		FROM drink_types ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	drinks := []api.DrinkType{}
	for rows.Next() {
		var d api.DrinkType
		if err := rows.Scan(&d.Id, &d.Name, &d.VolumeMl, &d.AlcoholPercent, &d.BeerEquivalent, &d.UnitName); err != nil {
			return nil, err
		}
		drinks = append(drinks, d)
	}
	return drinks, rows.Err()
}

func (s *Store) CreateDrink(ctx context.Context, d api.NewDrinkType, beerEquivalent float32) (api.DrinkType, error) {
	var out api.DrinkType
	err := s.pool.QueryRow(ctx, `
		INSERT INTO drink_types (name, volume_ml, alcohol_percent, beer_equivalent, unit_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, volume_ml, alcohol_percent, beer_equivalent, unit_name`,
		d.Name, d.VolumeMl, d.AlcoholPercent, beerEquivalent, d.UnitName).
		Scan(&out.Id, &out.Name, &out.VolumeMl, &out.AlcoholPercent, &out.BeerEquivalent, &out.UnitName)
	if pgCode(err) == codeUniqueViolation {
		return out, Conflict("a drink named %q already exists", d.Name)
	}
	return out, err
}
