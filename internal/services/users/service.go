package users

import (
	"context"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

type userRepo interface {
	Create(context.Context, *models.User) error
}

type Service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) CreateUser(
	ctx context.Context, request schemas.CreateUserRequest,
) (models.User, error) {
	user := models.User{
		Username: request.Username,
		Password: request.Password,
		IsAdmin:  request.IsAdmin,
	}
	err := s.userRepo.Create(ctx, &user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
