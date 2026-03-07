package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"github.com/surrealdb/surrealdb.go"
)

type VoteRepository struct{}

func (r *VoteRepository) Save(vote *models.Vote) error {
	ctx := context.Background()
	if vote.ID == "" {
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
