package models

import (
	"time"

	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type User struct {
	ID        *models.RecordID `json:"id,omitempty" surreal:"id,omitempty"`
	Username  string           `json:"username" surreal:"username"`
	CreatedAt time.Time        `json:"createdAt" surreal:"createdAt"`
	LastLogin time.Time        `json:"lastLogin" surreal:"lastLogin"`
}
