package service

import (
	"errors"
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

func (s *PollService) GetPolls() ([]models.Poll, error) {
	return s.PollRepo.FindAll()
}

func (s *PollService) GetPoll(id *surrealmodels.RecordID) (*models.Poll, error) {
	return s.PollRepo.FindByID(id)
}

func (s *PollService) GetPollsByPartyIdAndAttendee(id surrealmodels.RecordID, attendee surrealmodels.RecordID) (*models.Poll, error) {
	return s.PollRepo.FindByPartyIdAndAttendee(id, attendee)
}

func (s *PollService) GetVotedUsernamesByPartyId(partyId surrealmodels.RecordID) ([]string, error) {
	return s.PollRepo.FindVotedUsernamesByPartyId(partyId)
}

func (s *PollService) UpdatePoll(poll *models.Poll) (*models.Poll, error) {
	currentPoll, err := s.PollRepo.FindByID(poll.ID)
	if err != nil {
		return nil, err
	}

	if currentPoll.Status == models.PollStatusCompleted && poll.Status == models.PollStatusInProgress {
		return nil, errors.New("can't reactivate completed polls")
	}

	if currentPoll.Status == models.PollStatusInProgress && poll.Status == models.PollStatusCompleted {
		currentPoll.Status = models.PollStatusCompleted
		err = s.PollRepo.Save(currentPoll)
		if err != nil {
			return nil, err
		}
	}

	return currentPoll, nil
}

func (s *PollService) GetVotes(id string) (map[string]map[string]int, error) {
	votes, err := s.VoteRepo.FindByPollID(id)
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]int)
	for _, v := range votes {
		result[v.Attendee] = v.Choices
	}
	return result, nil
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
