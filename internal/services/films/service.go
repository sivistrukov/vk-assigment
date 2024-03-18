package films

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/infrastructure/postgresql"
	"github.com/sivistrukov/vk-assigment/internal/models"
	"github.com/sivistrukov/vk-assigment/internal/services/text"
)

type filmRepo interface {
	Create(context.Context, *models.Film, ...uint) error
	Update(context.Context, uint, map[string]any) error
	Remove(context.Context, uint) error
	GetFilmsWithActors(context.Context, string, string) ([]schemas.FilmWithActorsResponse, error)
}

type Service struct {
	filmRepo filmRepo
}

func NewService(filmRepo filmRepo) *Service {
	return &Service{
		filmRepo: filmRepo,
	}
}

func (s *Service) AddFilm(
	ctx context.Context, request schemas.AddFilmRequest,
) (models.Film, error) {
	film := models.Film{
		Title:       request.Title,
		Description: request.Description,
		ReleaseDate: request.ReleaseDate.ToTime(),
		Rating:      request.Rating,
	}

	err := s.filmRepo.Create(ctx, &film, request.ActorsIDs...)
	if err != nil {
		var notFoundErr *postgresql.ErrRecordNotFound
		if errors.As(err, &notFoundErr) {
			return models.Film{}, err
		}
		return models.Film{}, fmt.Errorf("error creating film: %v", err)
	}

	return film, nil
}

func (s *Service) UpdateFilm(
	ctx context.Context, id uint, request schemas.UpdateFilmRequest,
) error {
	reqType := reflect.TypeOf(request)
	reqValues := reflect.ValueOf(request)

	var updates = make(map[string]any, reqType.NumField())
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		value := reqValues.Field(i)

		var val any
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				val = nil
			} else {
				val = value.Elem().Interface()
			}
		} else {
			val = value.Interface()
		}

		if _, ok := val.(schemas.Date); ok {
			val = val.(schemas.Date).ToTime()
		}

		updates[text.CamelToSnake(field.Name)] = val
	}

	err := s.filmRepo.Update(ctx, id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) PartialUpdateFilm(
	ctx context.Context, id uint, request schemas.PartialUpdateFilmRequest,
) error {
	reqType := reflect.TypeOf(request)
	reqValues := reflect.ValueOf(request)

	var updates = make(map[string]any, reqType.NumField())
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		value := reqValues.Field(i)

		var val any
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			val = value.Elem().Interface()
		} else {
			val = value.Interface()
		}

		if _, ok := val.(schemas.Date); ok {
			val = val.(schemas.Date).ToTime()
		}

		updates[text.CamelToSnake(field.Name)] = val
	}

	err := s.filmRepo.Update(ctx, id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveFilm(ctx context.Context, filmId uint) error {
	return s.filmRepo.Remove(ctx, filmId)
}

func (s *Service) GetFilmsWithActors(
	ctx context.Context, search string, sortBy string,
) ([]schemas.FilmWithActorsResponse, error) {
	return s.filmRepo.GetFilmsWithActors(ctx, search, sortBy)
}
