package models

import "time"

type User struct {
	ID        string    `json:"id,omitempty" surreal:"id,omitempty"`
	Username  string    `json:"username" surreal:"username"`
	CreatedAt time.Time `json:"createdAt" surreal:"createdAt"`
	LastLogin time.Time `json:"lastLogin" surreal:"lastLogin"`
}
