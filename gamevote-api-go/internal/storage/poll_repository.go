package storage

import (
	"context"
	"fmt"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

type PollRepository struct{}

func (r *PollRepository) Save(poll *models.Poll) error {
	ctx := context.Background()
	if poll.ID == "" {
		slog.Debug("Creating new poll in DB")
		res, err := surrealdb.Create[models.Poll](ctx, DB, "polls", poll)
		if err == nil && res != nil {
			poll.ID = res.ID
		}
		return err
	}

	slog.Debug("Updating poll in DB", "id", poll.ID)

	_, err := surrealdb.Update[models.Poll](ctx, DB, poll.ID, poll)
	return err
}

func (r *PollRepository) FindAll() ([]models.Poll, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Poll](ctx, DB, "SELECT * FROM polls", nil)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.Poll{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *PollRepository) FindByID(id string) (*models.Poll, error) {
	ctx := context.Background()
	res, err := surrealdb.Select[models.Poll](ctx, DB, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("poll not found")
	}
	return res, nil
}

func (r *PollRepository) InitTable() error {
	ctx := context.Background()
	query := `
		DEFINE TABLE IF NOT EXISTS polls SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS options ON TABLE polls TYPE array<{name:string, appId: option<int>, imageUrl: option<string>}>;
		DEFINE FIELD IF NOT EXISTS attendees ON TABLE polls TYPE array<string>;
		DEFINE FIELD IF NOT EXISTS status ON TABLE polls TYPE string ASSERT $value INSIDE ['IN_PROGRESS', 'COMPLETED'];
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

func (r *PollRepository) DeleteAll() error {
	ctx := context.Background()
	_, err := surrealdb.Query[interface{}](ctx, DB, "DELETE polls", nil)
	return err
}
