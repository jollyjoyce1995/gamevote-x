package service

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"log/slog"
	"time"

	surrealmodels "github.com/surrealdb/surrealdb.go/pkg/models"
)

type UserService struct {
	UserRepo *storage.UserRepository
}

func NewUserService(userRepo *storage.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

type UserDTO struct {
	ID        *surrealmodels.RecordID `json:"-"`
	Username  string                  `json:"username"`
	CreatedAt time.Time               `json:"createdAt"`
	LastLogin time.Time               `json:"lastLogin"`
}

func (s *UserService) LoginOrRegister(username string) (*UserDTO, error) {
	slog.Debug("Login or register attempt", "username", username)
	user, err := s.UserRepo.Upsert(username)
	if err != nil {
		return nil, err
	}
	return s.ToDTO(user), nil
}

func (s *UserService) FindAll() ([]UserDTO, error) {
	users, err := s.UserRepo.FindAll()
	if err != nil {
		return nil, err
	}
	results := make([]UserDTO, 0, len(users))
	for _, n := range users {
		results = append(results, *s.ToDTO(&n))
	}

	return results, nil
}

func (s *UserService) FindByUsername(username string) (*UserDTO, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.ToDTO(user), nil
}

func (s *UserService) ToDTO(party *models.User) *UserDTO {
	return &UserDTO{
		ID:        party.ID,
		Username:  party.Username,
		CreatedAt: party.CreatedAt,
		LastLogin: party.LastLogin,
	}
}
