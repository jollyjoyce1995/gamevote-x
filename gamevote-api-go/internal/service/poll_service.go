package service

import (
	"errors"
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"sort"
)

type PollService struct {
	PollRepo *storage.PollRepository
	VoteRepo *storage.VoteRepository
}

func NewPollService(pollRepo *storage.PollRepository, voteRepo *storage.VoteRepository) *PollService {
	return &PollService{PollRepo: pollRepo, VoteRepo: voteRepo}
}

func (s *PollService) Create(poll *models.Poll) (*models.Poll, error) {
	poll.Status = models.PollStatusInProgress
	err := s.PollRepo.Save(poll)
	return poll, err
}

func (s *PollService) GetPolls() ([]models.Poll, error) {
	return s.PollRepo.FindAll()
}

func (s *PollService) GetPoll(id string) (*models.Poll, error) {
	return s.PollRepo.FindByID(id)
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

func (s *PollService) AddAttendee(id string, attendee string) error {
	poll, err := s.PollRepo.FindByID(id)
	if err != nil {
		return err
	}
	poll.Attendees = append(poll.Attendees, attendee)
	poll.Status = models.PollStatusInProgress
	return s.PollRepo.Save(poll)
}

func (s *PollService) AddVote(id string, attendee string, choices map[string]int) (map[string]int, error) {
	// Validate choices
	for _, v := range choices {
		if v < -1 || v > 1 {
			return nil, errors.New("bad request: vote value must be between -1 and 1")
		}
	}

	poll, err := s.PollRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	attendeeExists := false
	for _, a := range poll.Attendees {
		if a == attendee {
			attendeeExists = true
			break
		}
	}
	if !attendeeExists {
		return nil, errors.New("bad request: not an attendee")
	}

	if poll.Status != models.PollStatusInProgress {
		return nil, errors.New("bad request: poll is completed")
	}

	// check options
	for choiceKey := range choices {
		validOption := false
		for _, o := range poll.Options {
			if o == choiceKey {
				validOption = true
				break
			}
		}
		if !validOption {
			return nil, errors.New("bad request: some options are not valid")
		}
	}

	votes, err := s.VoteRepo.FindByPollID(id)
	if err != nil {
		return nil, err
	}

	for _, v := range votes {
		if v.Attendee == attendee {
			return nil, errors.New("bad request: cannot override vote")
		}
	}

	normalizedChoices := make(map[string]int)
	for _, option := range poll.Options {
		if val, exists := choices[option]; exists {
			normalizedChoices[option] = val
		} else {
			normalizedChoices[option] = 0
		}
	}

	newVote := &models.Vote{
		PollID:   id,
		Attendee: attendee,
		Choices:  normalizedChoices,
	}
	err = s.VoteRepo.Save(newVote)
	if err != nil {
		return nil, err
	}

	votes = append(votes, *newVote)

	// Check if all attendees have voted
	allVoted := true
	for _, a := range poll.Attendees {
		hasVoted := false
		for _, v := range votes {
			if v.Attendee == a {
				hasVoted = true
				break
			}
		}
		if !hasVoted {
			allVoted = false
			break
		}
	}

	if allVoted {
		poll.Status = models.PollStatusCompleted
		err = s.PollRepo.Save(poll)
		if err != nil {
			return nil, err
		}
	}

	return normalizedChoices, nil
}

func (s *PollService) GetResults(id string) (map[string]int, error) {
	poll, err := s.PollRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	votesMap, err := s.GetVotes(id)
	if err != nil {
		return nil, err
	}

	type optionResult struct {
		name  string
		score int
	}
	var results []optionResult

	for _, option := range poll.Options {
		sum := 0
		for _, attendeeVotes := range votesMap {
			sum += attendeeVotes[option]
		}
		results = append(results, optionResult{name: option, score: sum})
	}

	// Sort by descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	sortedMap := make(map[string]int)
	for _, r := range results {
		sortedMap[r.name] = r.score
	}

	return sortedMap, nil
}

func (s *PollService) GetOutstanding(id string) ([]string, error) {
	poll, err := s.PollRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	votesMap, err := s.GetVotes(id)
	if err != nil {
		return nil, err
	}

	var outstanding []string
	for _, a := range poll.Attendees {
		if _, exists := votesMap[a]; !exists {
			outstanding = append(outstanding, a)
		}
	}

	return outstanding, nil
}
