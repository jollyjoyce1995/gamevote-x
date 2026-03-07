package models

type PartyStatus string

const (
	PartyStatusNomination PartyStatus = "NOMINATION"
	PartyStatusVoting     PartyStatus = "VOTING"
	PartyStatusResults    PartyStatus = "RESULTS"
)

type Party struct {
	ID        string         `json:"id,omitempty" surreal:"id,omitempty"`
	Code      string         `json:"code" surreal:"code"`
	Attendees []string       `json:"attendees" surreal:"attendees"`
	Options   []string       `json:"options" surreal:"options"`
	Status    PartyStatus    `json:"status" surreal:"status"`
	Results   map[string]int `json:"results,omitempty" surreal:"results"`
	PollID    string         `json:"pollId,omitempty" surreal:"pollId"`
}
