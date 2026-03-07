package storage

import (
	"context"
	"fmt"

	"gamevote-api-go/internal/models"

	"github.com/surrealdb/surrealdb.go"
)

type PartyRepository struct{}

func (r *PartyRepository) Save(party *models.Party) error {
	ctx := context.Background()
	if party.ID == "" {
		res, err := surrealdb.Create[models.Party](ctx, DB, "parties", party)
		if err == nil && res != nil {
			party.ID = res.ID
		}
		return err
	}

	_, err := surrealdb.Update[models.Party](ctx, DB, party.ID, party)
	return err
}

func (r *PartyRepository) FindByID(id string) (*models.Party, error) {
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
