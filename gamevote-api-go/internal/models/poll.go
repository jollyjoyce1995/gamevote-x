package models

import surrealmodels "github.com/surrealdb/surrealdb.go/pkg/models"

type PollStatus string

const (
	PollStatusInProgress PollStatus = "IN_PROGRESS"
	PollStatusCompleted  PollStatus = "COMPLETED"
)

type PartyOptionWithVote struct {
	PartyOption
	Vote int `json:"vote" surreal:"vote"`
}

type Poll struct {
	ID       *surrealmodels.RecordID `json:"id,omitempty" surreal:"id,omitempty"`
	Options  []PartyOptionWithVote   `json:"options" surreal:"options"`
	Attendee *surrealmodels.RecordID `json:"attendee" surreal:"attendee"`
	Party    *surrealmodels.RecordID `json:"party" surreal:"party"`
	Status   PollStatus              `json:"status" surreal:"status"`
}
