package storage

import (
	"context"
	"fmt"

	"gamevote-api-go/internal/helpers"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
	surrealmodels "github.com/surrealdb/surrealdb.go/pkg/models"
)

type PartyRepository struct{}

func (r *PartyRepository) Save(party *models.Party) error {
	ctx := context.Background()
	slog.Info("Before updating party in DB", "id", party)

	if party.ID == nil {
		slog.Debug("Creating new party in DB", "code", party.Code)
		res, err := surrealdb.Create[models.Party](ctx, DB, "parties", party)
		if err == nil && res != nil {
			party.ID = res.ID
			return nil
		}
		return err
	}

	slog.Info("Updating party in DB", "id", party)

	_, err := surrealdb.Update[models.Party](ctx, DB, *party.ID, party)
	return err
}

func (r *PartyRepository) FindByID(id string) (*models.Party, error) {
	recordID, err := helpers.ToRecordID(id)
	if err != nil {
		return nil, err
	}
	return r.FindBySurrealID(*recordID)
}

func (r *PartyRepository) FindBySurrealID(id surrealmodels.RecordID) (*models.Party, error) {
	ctx := context.Background()
	res, err := surrealdb.Select[models.Party](ctx, DB, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("party not found")
	}
	return res, nil
}

func (r *PartyRepository) FindAll() ([]models.Party, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Party](ctx, DB, "SELECT * FROM parties ORDER BY id", nil)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.Party{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *PartyRepository) FindByCode(code string) (*models.Party, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Party](ctx, DB, "SELECT * FROM parties WHERE code = $code LIMIT 1", map[string]interface{}{"code": code})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 || len((*res)[0].Result) == 0 {
		return nil, fmt.Errorf("party not found by code")
	}
	party := (*res)[0].Result[0]
	return &party, nil
}

func (r *PartyRepository) ExistsByCode(code string) bool {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Party](ctx, DB, "SELECT * FROM parties WHERE code = $code LIMIT 1", map[string]interface{}{"code": code})
	if err == nil && res != nil && len(*res) > 0 && len((*res)[0].Result) > 0 {
		return true
	}
	return false
}

func (r *PartyRepository) InitTable() error {
	ctx := context.Background()
	query := `
		DEFINE TABLE IF NOT EXISTS parties SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS code ON TABLE parties TYPE string;
		DEFINE FIELD IF NOT EXISTS attendees ON TABLE parties TYPE array<string>;
		DEFINE FIELD IF NOT EXISTS options ON TABLE parties TYPE array<{name:string, appId: int, imageUrl: string}>;
		DEFINE FIELD IF NOT EXISTS status ON TABLE parties TYPE string ASSERT $value INSIDE ['NOMINATION', 'VOTING', 'RESULTS'];
		DEFINE FIELD IF NOT EXISTS results ON TABLE parties TYPE option<object>;
		DEFINE FIELD IF NOT EXISTS pollId ON TABLE parties TYPE option<string>;
	`
	slog.Debug("Initializing parties table")
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

func (r *PartyRepository) DeleteAll() error {
	ctx := context.Background()
	_, err := surrealdb.Query[interface{}](ctx, DB, "DELETE parties", nil)
	return err
}
