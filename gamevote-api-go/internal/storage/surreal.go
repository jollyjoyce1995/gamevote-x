package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func InitDB() error {
	wsUrl := os.Getenv("SURREAL_WS")
	if wsUrl == "" {
		wsUrl = "ws://localhost:8000/rpc"
	}

	db, err := surrealdb.New(wsUrl)
	if err != nil {
		return fmt.Errorf("failed to open surreal connection: %w", err)
	}

	DB = db

	user := os.Getenv("SURREAL_USER")
	pass := os.Getenv("SURREAL_PASS")
	if user != "" && pass != "" {
		authData := map[string]interface{}{
			"user": user,
			"pass": pass,
		}
		if _, err = DB.SignIn(context.Background(), authData); err != nil {
			return fmt.Errorf("failed to signin: %w", err)
		}
	}

	ns := os.Getenv("SURREAL_NS")
	if ns == "" {
		ns = "gamevote"
	}
	dbName := os.Getenv("SURREAL_DB")
	if dbName == "" {
		dbName = "gamevote"
	}

	if err = DB.Use(context.Background(), ns, dbName); err != nil {
		return fmt.Errorf("failed to use namespace/db: %w", err)
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close(context.Background())
	}
}
