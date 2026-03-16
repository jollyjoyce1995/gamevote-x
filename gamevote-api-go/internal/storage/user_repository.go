package storage

import (
	"context"
	"fmt"
	"gamevote-api-go/internal/models"
	"time"

	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

type UserRepository struct{}

func (r *UserRepository) Save(user *models.User) error {
	ctx := context.Background()
	if user.ID == nil {
		res, err := surrealdb.Create[models.User](ctx, DB, "users", user)
		if err == nil && res != nil {
			user.ID = res.ID
		}
		return err
	}

	_, err := surrealdb.Update[models.User](ctx, DB, *user.ID, user)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.User](ctx, DB, "SELECT * FROM users WHERE username = $username LIMIT 1", map[string]interface{}{"username": username})
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 || len((*res)[0].Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	user := (*res)[0].Result[0]
	return &user, nil
}

func (r *UserRepository) Upsert(username string) (*models.User, error) {
	slog.Debug("Upserting user", "username", username)
	user, err := r.FindByUsername(username)
	slog.Info("User found", "user", user)
	now := time.Now()

	if err != nil {
		// Does not exist
		user = &models.User{
			Username:  username,
			CreatedAt: now,
			LastLogin: now,
		}
		err = r.Save(user)
		return user, err
	}

	// Exists, update last login
	user.LastLogin = now
	err = r.Save(user)
	return user, err
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	ctx := context.Background()
	res, err := surrealdb.Query[[]models.User](ctx, DB, "SELECT * FROM users", nil)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return []models.User{}, nil
	}
	return (*res)[0].Result, nil
}

func (r *UserRepository) InitTable() error {
	ctx := context.Background()
	query := `
		DEFINE TABLE IF NOT EXISTS users SCHEMAFULL;
		DEFINE FIELD IF NOT EXISTS username ON TABLE users TYPE string;
		DEFINE FIELD IF NOT EXISTS createdAt ON TABLE users TYPE datetime;
		DEFINE FIELD IF NOT EXISTS lastLogin ON TABLE users TYPE datetime;
	`
	_, err := surrealdb.Query[interface{}](ctx, DB, query, nil)
	return err
}

