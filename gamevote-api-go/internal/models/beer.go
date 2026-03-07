package models

type Beer struct {
	ID       string `json:"id,omitempty" surreal:"id,omitempty"`
	PartyID  string `json:"partyId" surreal:"partyId"`
	Attendee string `json:"attendee" surreal:"attendee"`
}
