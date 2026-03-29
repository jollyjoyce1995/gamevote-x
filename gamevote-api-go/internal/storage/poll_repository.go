package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
	surrealmodels "github.com/surrealdb/surrealdb.go/pkg/models"
)

type PollRepository struct{}

func (r *PollRepository) Save(poll *models.Poll) error {
	ctx := context.Background()
	slog.Info("Before updating poll in DB", "id", poll)
	if poll.ID == nil {
		slog.Debug("Creating new poll in DB")
		res, err := surrealdb.Create[models.Poll](ctx, DB, "polls", poll)
		if err == nil && res != nil {
			poll.ID = res.ID
		}
		return err
	}

	slog.Debug("Updating poll in DB", "id", poll.ID)

	_, err := surrealdb.Update[models.Poll](ctx, DB, *poll.ID, poll)
	return err
}


func (r *PollRepository) FindByPartyIdAndAttendee(partyId surrealmodels.RecordID, attendee surrealmodels.RecordID, round int) (*models.Poll, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Poll](ctx, DB, "SELECT * FROM polls where party.id = $partyId and attendee.id = $attendee and round = $round", map[string]interface{}{"partyId": partyId, "attendee": attendee, "round": round})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 || len((*res)[0].Result) == 0 {
		return nil, nil
	}
	return &(*res)[0].Result[0], nil
}


func (r *PollRepository) FindAllByPartyId(partyId surrealmodels.RecordID) ([]models.Poll, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Poll](ctx, DB, "SELECT * FROM polls WHERE party = $partyId", map[string]interface{}{"partyId": partyId})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.Poll{}, nil
	}
	return (*res)[0].Result, nil
}

type MaxRoundResult struct {
	MaxRound int `json:"maxRound"`
}

func (r *PollRepository) FindMaxRoundByPartyId(partyId surrealmodels.RecordID) (int, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]MaxRoundResult](ctx, DB, "SELECT math::max(round) AS maxRound FROM polls WHERE party = $partyId AND status = 'COMPLETED' GROUP ALL", map[string]interface{}{"partyId": partyId})
	if err != nil {
		return 0, err
	}
	if res == nil || len(*res) == 0 || len((*res)[0].Result) == 0 {
		return 0, nil
	}
	return (*res)[0].Result[0].MaxRound, nil
}

func (r *PollRepository) MarkAllInProgressAsCompleted(partyId surrealmodels.RecordID) error {
	ctx := context.Background()
	_, err := surrealdb.Query[interface{}](ctx, DB, "UPDATE polls SET status = 'COMPLETED' WHERE party = $partyId AND status = 'IN_PROGRESS'", map[string]interface{}{"partyId": partyId})
	return err
}

type VotedUsernameResult struct {
	Username string `json:"username"`
}

func (r *PollRepository) FindVotedUsernamesByPartyIdAndRound(partyId surrealmodels.RecordID, round int) ([]string, error) {
	ctx := context.Background()
	query := "SELECT attendee.username AS username FROM polls WHERE party = $partyId AND round = $round"

	res, err := surrealdb.Query[[]VotedUsernameResult](ctx, DB, query, map[string]interface{}{"partyId": partyId, "round": round})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []string{}, nil
	}

	usernames := make([]string, 0, len((*res)[0].Result))
	for _, row := range (*res)[0].Result {
		usernames = append(usernames, row.Username)
	}

	return usernames, nil
}

func (r *PollRepository) InitTable() error {
	ctx := context.Background()
	query := `
		DEFINE TABLE IF NOT EXISTS polls SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS options ON TABLE polls TYPE array<{name:string, appId: option<int>, imageUrl: option<string>, vote: int}>;
		DEFINE FIELD IF NOT EXISTS attendee ON TABLE polls TYPE record<users>;
		DEFINE FIELD IF NOT EXISTS party ON TABLE polls TYPE record<parties>;
		DEFINE FIELD IF NOT EXISTS status ON TABLE polls TYPE string ASSERT $value INSIDE ['IN_PROGRESS', 'COMPLETED'];
		DEFINE FIELD IF NOT EXISTS round ON TABLE polls TYPE int;
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

