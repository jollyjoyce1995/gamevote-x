package service

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"

	surrealmodels "github.com/surrealdb/surrealdb.go/pkg/models"
)

type PollService struct {
	PollRepo *storage.PollRepository
	VoteRepo *storage.VoteRepository
}

func NewPollService(pollRepo *storage.PollRepository, voteRepo *storage.VoteRepository) *PollService {
	return &PollService{PollRepo: pollRepo, VoteRepo: voteRepo}
}

func (s *PollService) Upsert(poll *models.Poll) (*models.Poll, error) {
	err := s.PollRepo.Save(poll)
	return poll, err
}

func (s *PollService) GetPollsByPartyIdAndAttendee(id surrealmodels.RecordID, attendee surrealmodels.RecordID, round int) (*models.Poll, error) {
	return s.PollRepo.FindByPartyIdAndAttendee(id, attendee, round)
}

func (s *PollService) MarkAllInProgressAsCompleted(partyId surrealmodels.RecordID) error {
	return s.PollRepo.MarkAllInProgressAsCompleted(partyId)
}

func (s *PollService) GetMaxRoundByPartyId(partyId surrealmodels.RecordID) (int, error) {
	return s.PollRepo.FindMaxRoundByPartyId(partyId)
}

func (s *PollService) GetVotedUsernamesByPartyId(partyId surrealmodels.RecordID) ([]string, error) {
	return s.PollRepo.FindVotedUsernamesByPartyId(partyId)
}

func (s *PollService) GetResultsByPartyId(partyId surrealmodels.RecordID) ([]models.Poll, error) {
	polls, err := s.PollRepo.FindAllByPartyId(partyId)
	if err != nil {
		return nil, err
	}

	return polls, nil
}
