package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

type GameRepository struct{}

func (r *GameRepository) Count() (int, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]struct {
		Count int `json:"count"`
	}](ctx, DB, "SELECT count() FROM games GROUP ALL", nil)
	if err != nil {
		return 0, err
	}
	if res == nil || len(*res) == 0 || len((*res)[0].Result) == 0 {
		return 0, nil
	}
	return (*res)[0].Result[0].Count, nil
}

func (r *GameRepository) BulkInsert(games []models.Game) error {
	ctx := context.Background()
	slog.Info("Starting bulk insert of games", "count", len(games))
	for _, g := range games {
		if _, err := surrealdb.Create[models.Game](ctx, DB, "games", &g); err != nil {
			// skip duplicates / errors silently during bulk insert
			continue
		}
	}
	slog.Info("Completed bulk insert of games")
	return nil
}

func (r *GameRepository) Search(query string, limit int) ([]models.Game, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Game](ctx, DB, "SELECT * FROM games WHERE string::lowercase(name) CONTAINS string::lowercase($q) LIMIT $limit", map[string]interface{}{"q": query, "limit": limit})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.Game{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *GameRepository) InitTable() error {
	ctx := context.Background()
	query := `
		IF (SELECT VALUE id FROM (INFO FOR DB).tables.games) == NONE {
			DEFINE TABLE games SCHEMAFULL;
			DEFINE FIELD appId ON TABLE games TYPE int;
			DEFINE FIELD name ON TABLE games TYPE string;
			DEFINE FIELD imageUrl ON TABLE games TYPE option<string>;
		};
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

func (r *GameRepository) DeleteAll() error {
	ctx := context.Background()
	// Check if table exists before deleting or handle error
	_, err := surrealdb.Query[interface{}](ctx, DB, "DELETE games", nil)
	if err != nil {
		// If error contains "does not exist", we can ignore it for a clear operation
		// However, with SCHEMAFULL and InitDB, it SHOULD exist.
		// We'll keep it simple for now as InitDB will ensure it exists.
		return err
	}
	return nil
}
