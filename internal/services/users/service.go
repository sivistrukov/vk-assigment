package users

import (
	"context"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
	"github.com/sivistrukov/vk-assigment/internal/services/auth"
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
	password, err := auth.HashPassword(request.Password)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		Username: request.Username,
		Password: password,
		IsAdmin:  request.IsAdmin,
	}

	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
