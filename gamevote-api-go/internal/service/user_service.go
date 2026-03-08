package service

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"log/slog"
)

type UserService struct {
	UserRepo *storage.UserRepository
}

func NewUserService(userRepo *storage.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) LoginOrRegister(username string) (*models.User, error) {
	slog.Debug("Login or register attempt", "username", username)
	return s.UserRepo.Upsert(username)
}

func (s *UserService) FindAll() ([]models.User, error) {
	return s.UserRepo.FindAll()
}
