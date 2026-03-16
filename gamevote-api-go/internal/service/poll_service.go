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



func (s *PollService) GetPollsByPartyIdAndAttendee(id surrealmodels.RecordID, attendee surrealmodels.RecordID) (*models.Poll, error) {
	return s.PollRepo.FindByPartyIdAndAttendee(id, attendee)
}

func (s *PollService) GetVotedUsernamesByPartyId(partyId surrealmodels.RecordID) ([]string, error) {
	return s.PollRepo.FindVotedUsernamesByPartyId(partyId)
}



func (s *PollService) GetResultsByPartyId(partyId surrealmodels.RecordID) (map[string]int, error) {
	polls, err := s.PollRepo.FindAllByPartyId(partyId)
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)
	
	// Initialize all options to 0 based on the first poll (or any poll) to ensure all games are present
	if len(polls) > 0 {
		for _, opt := range polls[0].Options {
			results[opt.Name] = 0
		}
	}

	// Sum votes for all options across all polls
	for _, poll := range polls {
		for _, opt := range poll.Options {
			results[opt.Name] += opt.Vote
		}
	}

	return results, nil
}
