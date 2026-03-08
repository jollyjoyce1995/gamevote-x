package service

import (
	"errors"
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"log/slog"
	"math/rand"
	"time"

	surrealmodels "github.com/surrealdb/surrealdb.go/pkg/models"
)

type PartyService struct {
	PartyRepo   *storage.PartyRepository
	BeerRepo    *storage.BeerRepository
	PollService *PollService
	Broker      *SSEBroker
}

func NewPartyService(partyRepo *storage.PartyRepository, beerRepo *storage.BeerRepository, pollService *PollService, broker *SSEBroker) *PartyService {
	rand.Seed(time.Now().UnixNano())
	return &PartyService{
		PartyRepo:   partyRepo,
		BeerRepo:    beerRepo,
		PollService: pollService,
		Broker:      broker,
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
		slog.Warn("Random code collision, retrying", "code", randomCode)
	}
	return randomCode
}

func (s *PartyService) CreateParty(party *models.Party) (*models.Party, error) {
	party.Code = s.createCodeForParty()
	if party.Attendees == nil {
		party.Attendees = []string{}
	}
	if party.Options == nil {
		party.Options = []models.PartyOption{}
	}
	party.Status = models.PartyStatusNomination
	slog.Info("Saving new party to database", "code", party.Code)
	err := s.PartyRepo.Save(party)
	return party, err
}

func (s *PartyService) GetParty(id string) (*models.Party, error) {
	return s.PartyRepo.FindByID(id)
}

func (s *PartyService) GetParties() ([]*PartyDTO, error) {
	parties, err := s.PartyRepo.FindAll()
	if err != nil {
		return nil, err
	}
	dtos := make([]*PartyDTO, 0, len(parties))
	for i := range parties {
		dto, err := s.ToDTO(&parties[i])
		if err != nil {
			continue
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (s *PartyService) GetPartyByCode(code string) (*PartyDTO, error) {
	party, err := s.PartyRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	return s.ToDTO(party)
}

func (s *PartyService) AllowedTransitions(id surrealmodels.RecordID) (map[models.PartyStatus]bool, error) {
	party, err := s.PartyRepo.FindBySurrealID(id)
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

func (s *PartyService) PatchParty(id surrealmodels.RecordID, toStatus models.PartyStatus) (*models.Party, error) {
	party, err := s.PartyRepo.FindBySurrealID(id)
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
	slog.Info("Committing party status change", "code", party.Code, "from", fromStatus, "to", toStatus)
	err = s.PartyRepo.Save(party)
	if err != nil {
		slog.Error("Failed to save party status change", "code", party.Code, "error", err)
		return nil, err
	}

	return party, nil
}

func (s *PartyService) AddOption(id string, option models.PartyOption) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}
	slog.Info("Adding option to party", "code", party.Code, "option", option.Name)
	for _, o := range party.Options {
		if o.Name == option.Name {
			slog.Warn("Option already exists in party", "code", party.Code, "option", option.Name)
			return errors.New("bad request: game already nominated")
		}
	}

	party.Options = append(party.Options, option)
	slog.Info("Saving party with new option", "code", party.Code, "option", option.Name)
	if err := s.PartyRepo.Save(party); err != nil {
		return err
	}
	// Broadcast updated state
	if s.Broker != nil {
		dto, _ := s.ToDTO(party)
		s.Broker.Broadcast(party.Code, "party_updated", dto)
	}
	return nil
}

func (s *PartyService) AddAttendee(id surrealmodels.RecordID, value string) error {
	party, err := s.PartyRepo.FindBySurrealID(id)
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

	if err := s.PartyRepo.Save(party); err != nil {
		return err
	}
	// Broadcast updated state
	if s.Broker != nil {
		dto, _ := s.ToDTO(party)
		s.Broker.Broadcast(party.Code, "party_updated", dto)
	}
	return nil
}

func (s *PartyService) DeleteAttendee(id surrealmodels.RecordID, index int) error {
	party, err := s.PartyRepo.FindBySurrealID(id)
	if err != nil {
		return err
	}

	if index >= 0 && index < len(party.Attendees) {
		party.Attendees = append(party.Attendees[:index], party.Attendees[index+1:]...)
		return s.PartyRepo.Save(party)
	}
	return errors.New("not found: index out of bounds")
}

func (s *PartyService) DeleteOption(id string, name string) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}

	found := false
	for i, o := range party.Options {
		if o.Name == name {
			party.Options = append(party.Options[:i], party.Options[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("not found: game not nominated")
	}

	if err := s.PartyRepo.Save(party); err != nil {
		return err
	}
	// Broadcast updated state
	if s.Broker != nil {
		dto, _ := s.ToDTO(party)
		s.Broker.Broadcast(party.Code, "party_updated", dto)
	}
	return nil
}

func (s *PartyService) PostBeer(id string, attendee string) error {
	party, err := s.PartyRepo.FindByID(id)
	if err != nil {
		return err
	}

	beer := &models.Beer{
		PartyID:  party.ID.String(),
		Attendee: attendee,
	}
	return s.BeerRepo.Save(beer)
}

// Helper models for the responses
type PartyDTO struct {
	ID              string               `json:"id"`
	Attendees       []string             `json:"attendees"`
	Options         []models.PartyOption `json:"options"`
	Status          string               `json:"status"`
	Results         map[string]int       `json:"results,omitempty"`
	Code            string               `json:"code,omitempty"`
	Links           map[string]Link      `json:"_links"`
	BeerCount       int                  `json:"beerCount"`
	BeerPerAttendee map[string]int       `json:"beerPerAttendee"`
}

type Link struct {
	Href string `json:"href"`
}

func (s *PartyService) ToDTO(party *models.Party) (*PartyDTO, error) {
	beers, err := s.BeerRepo.FindByPartyID(party.ID.String())
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
		"self": {Href: "/parties/" + party.ID.String()},
	}
	if party.PollID != "" {
		links["poll"] = Link{Href: "/polls/" + party.PollID}
	}

	return &PartyDTO{
		ID:              party.ID.String(),
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

func (dto *PartyDTO) ToDomain() (*models.Party, error) {
	var recordID *surrealmodels.RecordID // Defaults to nil

	if dto.ID != "" {
		// Create the ID and assign the address to your pointer
		id := surrealmodels.NewRecordID("parties", dto.ID)
		recordID = &id
	}
	options := make([]models.PartyOption, len(dto.Options))
	for i, opt := range dto.Options {
		options[i] = models.PartyOption{
			Name:     opt.Name,
			AppID:    opt.AppID,
			ImageURL: opt.ImageURL,
		}
	}

	return &models.Party{
		ID:        recordID,
		Code:      dto.Code,
		Attendees: dto.Attendees,
		Options:   options,
		Status:    models.PartyStatus(dto.Status),
		Results:   dto.Results,
	}, nil
}
