package storage

import (
	"context"
	"gamevote-api-go/internal/models"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

type DrinkTypeRepository struct{}

func (r *DrinkTypeRepository) Save(dt *models.DrinkType) error {
	ctx := context.Background()
	if dt.ID == "" {
		slog.Debug("Creating new drink type in DB", "name", dt.Name)
		res, err := surrealdb.Create[models.DrinkType](ctx, DB, "drink_types", dt)
		if err == nil && res != nil {
			dt.ID = res.ID
		}
		return err
	}

	_, err := surrealdb.Update[models.DrinkType](ctx, DB, dt.ID, dt)
	return err
}

func (r *DrinkTypeRepository) FindAll() ([]models.DrinkType, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.DrinkType](ctx, DB, "SELECT * FROM drink_types", nil)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.DrinkType{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *DrinkTypeRepository) InitTable() error {
	ctx := context.Background()
	query := `
		DEFINE TABLE IF NOT EXISTS drink_types SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS name ON TABLE drink_types TYPE string;
		DEFINE FIELD IF NOT EXISTS volumeMl ON TABLE drink_types TYPE int;
		DEFINE FIELD IF NOT EXISTS alcoholPercent ON TABLE drink_types TYPE float;
		DEFINE FIELD IF NOT EXISTS beerEquivalent ON TABLE drink_types TYPE float;
		DEFINE FIELD IF NOT EXISTS unitName ON TABLE drink_types TYPE string;
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

