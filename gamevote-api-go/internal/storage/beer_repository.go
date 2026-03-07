package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"github.com/surrealdb/surrealdb.go"
)

type BeerRepository struct{}

func (r *BeerRepository) Save(beer *models.Beer) error {
	ctx := context.Background()
	if beer.ID == "" {
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
