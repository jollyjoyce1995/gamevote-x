package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

type BeerRepository struct{}

func (r *BeerRepository) Save(beer *models.Beer) error {
	ctx := context.Background()
	if beer.ID == "" {
		slog.Debug("Creating new beer record in DB", "partyId", beer.PartyID, "attendee", beer.Attendee)
		res, err := surrealdb.Create[models.Beer](ctx, DB, "beers", beer)
		if err == nil && res != nil {
			beer.ID = res.ID
		}
		return err
	}
	_, err := surrealdb.Update[models.Beer](ctx, DB, beer.ID, beer)
	return err
}

func (r *BeerRepository) FindByPartyID(partyID string) ([]models.Beer, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.Beer](ctx, DB, "SELECT * FROM beers WHERE partyId = $partyId", map[string]interface{}{"partyId": partyID})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.Beer{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *BeerRepository) InitTable() error {
	ctx := context.Background()
	query := `
		IF (SELECT VALUE id FROM (INFO FOR DB).tables.beers) == NONE {
			DEFINE TABLE beers SCHEMAFULL;
			DEFINE FIELD partyId ON TABLE beers TYPE string;
			DEFINE FIELD attendee ON TABLE beers TYPE string;
		};
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

func (r *BeerRepository) DeleteAll() error {
	ctx := context.Background()
	_, err := surrealdb.Query[interface{}](ctx, DB, "DELETE beers", nil)
	return err
}
