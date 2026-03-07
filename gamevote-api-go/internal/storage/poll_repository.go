package storage

import (
	"context"
	"fmt"
	"gamevote-api-go/internal/models"

	"github.com/surrealdb/surrealdb.go"
)

type PollRepository struct{}

func (r *PollRepository) Save(poll *models.Poll) error {
	ctx := context.Background()
	if poll.ID == "" {
		res, err := surrealdb.Create[models.Poll](ctx, DB, "polls", poll)
		if err == nil && res != nil {
			poll.ID = res.ID
		}
		return err
	}

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
