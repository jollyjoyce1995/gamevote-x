package models

type Vote struct {
	ID       string         `json:"id,omitempty" surreal:"id,omitempty"`
	PollID   string         `json:"pollId" surreal:"pollId"`
	Attendee string         `json:"attendee" surreal:"attendee"`
	Choices  map[string]int `json:"choices" surreal:"choices"`
}
