package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

type VoteRepository struct{}

func (r *VoteRepository) Save(vote *models.Vote) error {
	ctx := context.Background()
	if vote.ID == "" {
		slog.Debug("Creating new vote in DB", "pollId", vote.PollID, "attendee", vote.Attendee)
		res, err := surrealdb.Create[models.Vote](ctx, DB, "votes", vote)
		if err == nil && res != nil {
			vote.ID = res.ID
		}
		return err
	}
	_, err := surrealdb.Update[models.Vote](ctx, DB, vote.ID, vote)
	return err
}

func (r *VoteRepository) FindByPollID(pollID string) ([]models.Vote, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Vote](ctx, DB, "SELECT * FROM votes WHERE pollId = $pollId", map[string]interface{}{"pollId": pollID})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.Vote{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *VoteRepository) InitTable() error {
	ctx := context.Background()
	query := `
		DEFINE TABLE IF NOT EXISTS votes SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS pollId ON TABLE votes TYPE string;
		DEFINE FIELD IF NOT EXISTS attendee ON TABLE votes TYPE string;
		DEFINE FIELD IF NOT EXISTS choices ON TABLE votes TYPE object;
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

