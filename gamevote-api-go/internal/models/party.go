package models

import (
	"time"

	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type PartyStatus string

const (
	PartyStatusNomination PartyStatus = "NOMINATION"
	PartyStatusVoting     PartyStatus = "VOTING"
	PartyStatusResults    PartyStatus = "RESULTS"
)

type PartyOption struct {
	Name     string `json:"name" surreal:"name"`
	AppID    int    `json:"appId,omitempty" surreal:"appId"`
	ImageURL string `json:"imageUrl,omitempty" surreal:"imageUrl"`
}

type Party struct {
	ID        *models.RecordID `json:"id,omitempty"`
	Code      string           `json:"code" surreal:"code"`
	Attendees []string         `json:"attendees" surreal:"attendees"`
	Options   []PartyOption    `json:"options" surreal:"options"`
	Status    PartyStatus      `json:"status" surreal:"status"`
	CreatedAt time.Time        `json:"createdAt" surreal:"createdAt"`
}
