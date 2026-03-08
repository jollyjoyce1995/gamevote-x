package models

type PollStatus string

const (
	PollStatusInProgress PollStatus = "IN_PROGRESS"
	PollStatusCompleted  PollStatus = "COMPLETED"
)

type Poll struct {
	ID        string        `json:"id,omitempty" surreal:"id,omitempty"`
	Options   []PartyOption `json:"options" surreal:"options"`
	Attendees []string      `json:"attendees" surreal:"attendees"`
	Status    PollStatus    `json:"status" surreal:"status"`
}
