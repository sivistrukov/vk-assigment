package auth

import (
	"context"

	"github.com/sivistrukov/vk-assigment/internal/models"
)

type UserRepo interface {
	GetByUsername(context.Context, string) (models.User, error)
}

type Service struct {
	userRepo UserRepo
}

func NewService(userRepo UserRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) Authenticate(username string, password string) (models.User, error) {
	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return models.User{}, ErrNotAuthorized
	}

	if !ComparePasswordAndHash(password, user.Password) {
		return models.User{}, ErrNotAuthorized
	}

	return user, nil
}
