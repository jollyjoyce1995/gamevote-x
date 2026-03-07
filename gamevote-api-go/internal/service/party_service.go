package service

import (
	"errors"
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"math/rand"
	"time"
)

type PartyService struct {
	PartyRepo   *storage.PartyRepository
	BeerRepo    *storage.BeerRepository
	PollService *PollService
}

func NewPartyService(partyRepo *storage.PartyRepository, beerRepo *storage.BeerRepository, pollService *PollService) *PartyService {
	rand.Seed(time.Now().UnixNano())
	return &PartyService{
		PartyRepo:   partyRepo,
		BeerRepo:    beerRepo,
		PollService: pollService,
	}
}

func (s *PartyService) createRandomCode() string {
	possibleCharacters := "ABCDEFGHJKLMNPQRSTUVQXYZ23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = possibleCharacters[rand.Intn(len(possibleCharacters))]
	}
	return string(b)
}

func (s *PartyService) createCodeForParty() string {
	var randomCode string
	for {
		randomCode = s.createRandomCode()
		if !s.PartyRepo.ExistsByCode(randomCode) {
			break
		}
	}
	return randomCode
}

func (s *PartyService) CreateParty(party *models.Party) (*models.Party, error) {
	party.Code = s.createCodeForParty()
	if party.Attendees == nil {
		party.Attendees = []string{}
	}
	if party.Options == nil {
		party.Options = []string{}
	}
	party.Status = models.PartyStatusNomination
	err := s.PartyRepo.Save(party)
	return party, err
}

func (s *PartyService) GetParty(id string) (*models.Party, error) {
	return s.PartyRepo.FindByID(id)
}

func (s *PartyService) GetIdForCode(code string) (string, error) {
	party, err := s.PartyRepo.FindByCode(code)
	if err != nil {
		return "", err
	}
	return party.ID, nil
}

func (s *PartyService) AllowedTransitions(id string) (map[models.PartyStatus]bool, error) {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	transitions := make(map[models.PartyStatus]bool)
	switch party.Status {
	case models.PartyStatusNomination:
		if len(party.Options) > 0 {
			transitions[models.PartyStatusVoting] = true
		}
	case models.PartyStatusVoting:
		transitions[models.PartyStatusNomination] = true
		transitions[models.PartyStatusResults] = true
	case models.PartyStatusResults:
		transitions[models.PartyStatusNomination] = true
		transitions[models.PartyStatusVoting] = true
	}
	return transitions, nil
}

func (s *PartyService) PatchParty(id string, toStatus models.PartyStatus) (*models.Party, error) {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	fromStatus := party.Status
	if fromStatus == toStatus {
		return party, nil
	}

	transitions, err := s.AllowedTransitions(id)
	if err != nil {
		return nil, err
	}

	if !transitions[toStatus] {
		return nil, errors.New("bad request: illegal transition")
	}

	if toStatus == models.PartyStatusVoting || toStatus == models.PartyStatusNomination {
		if party.PollID != "" {
			poll, err := s.PollService.GetPoll(party.PollID)
			if err == nil {
				poll.Status = models.PollStatusCompleted
				s.PollService.UpdatePoll(poll)
			}
			party.PollID = ""
			party.Results = map[string]int{}
		}
	}

	if toStatus == models.PartyStatusVoting {
		poll := &models.Poll{
			Options:   party.Options,
			Attendees: party.Attendees,
		}
		createdPoll, err := s.PollService.Create(poll)
		if err != nil {
			return nil, err
		}
		party.PollID = createdPoll.ID
	} else if toStatus == models.PartyStatusResults {
		if party.PollID != "" {
			poll, err := s.PollService.GetPoll(party.PollID)
			if err == nil {
				results, err := s.PollService.GetResults(poll.ID)
				if err == nil {
					party.Results = results
				}
				poll.Status = models.PollStatusCompleted
				s.PollService.UpdatePoll(poll)
			}
		}
	}

	party.Status = toStatus
	err = s.PartyRepo.Save(party)
	if err != nil {
		return nil, err
	}

	return party, nil
}

func (s *PartyService) AddOption(id string, value string) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}
	party.Options = append(party.Options, value)
	return s.PartyRepo.Save(party)
}

func (s *PartyService) AddAttendee(id string, value string) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}

	for _, a := range party.Attendees {
		if a == value {
			return errors.New("bad request: attendee already exists")
		}
	}

	party.Attendees = append(party.Attendees, value)
	if party.PollID != "" {
		s.PollService.AddAttendee(party.PollID, value)
	}

	return s.PartyRepo.Save(party)
}

func (s *PartyService) DeleteAttendee(id string, index int) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}

	if index >= 0 && index < len(party.Attendees) {
		party.Attendees = append(party.Attendees[:index], party.Attendees[index+1:]...)
		return s.PartyRepo.Save(party)
	}
	return errors.New("not found: index out of bounds")
}

func (s *PartyService) DeleteOption(id string, index int) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}

	if index >= 0 && index < len(party.Options) {
		party.Options = append(party.Options[:index], party.Options[index+1:]...)
		return s.PartyRepo.Save(party)
	}
	return errors.New("not found: index out of bounds")
}

func (s *PartyService) PostBeer(id string, attendee string) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}

	beer := &models.Beer{
		PartyID:  party.ID,
		Attendee: attendee,
	}
	return s.BeerRepo.Save(beer)
}

// Helper models for the responses
type PartyDTO struct {
	ID              string          `json:"id"`
	Attendees       []string        `json:"attendees"`
	Options         []string        `json:"options"`
	Status          string          `json:"status"`
	Results         map[string]int  `json:"results,omitempty"`
	Code            string          `json:"code,omitempty"`
	Links           map[string]Link `json:"_links"`
	BeerCount       int             `json:"beerCount"`
	BeerPerAttendee map[string]int  `json:"beerPerAttendee"`
}

type Link struct {
	Href string `json:"href"`
}

func (s *PartyService) ToDTO(party *models.Party) (*PartyDTO, error) {
	beers, err := s.BeerRepo.FindByPartyID(party.ID)
	if err != nil {
		return nil, err
	}

	beerPerAttendee := make(map[string]int)
	for _, a := range party.Attendees {
		count := 0
		for _, b := range beers {
			if b.Attendee == a {
				count++
			}
		}
		beerPerAttendee[a] = count
	}

	links := map[string]Link{
		"self": {Href: "/parties/" + party.ID},
	}
	if party.PollID != "" {
		links["poll"] = Link{Href: "/polls/" + party.PollID}
	}

	return &PartyDTO{
		ID:              party.ID,
		Attendees:       party.Attendees,
		Options:         party.Options,
		Status:          string(party.Status),
		Results:         party.Results,
		Code:            party.Code,
		Links:           links,
		BeerCount:       len(beers),
		BeerPerAttendee: beerPerAttendee,
	}, nil
}
