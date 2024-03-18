package actors

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
	"github.com/sivistrukov/vk-assigment/internal/services/text"
)

type actorRepo interface {
	Create(context.Context, *models.Actor) error
	Update(context.Context, uint, map[string]any) error
	Remove(context.Context, uint) error
	GetListWithFilms(context.Context) ([]schemas.ActorWithFilmsResponse, error)
}

type Service struct {
	actorRepo actorRepo
}

func NewService(actorRepo actorRepo) *Service {
	return &Service{
		actorRepo: actorRepo,
	}
}

func (s *Service) AddActor(
	ctx context.Context, request schemas.AddActorRequest,
) (models.Actor, error) {
	actor := models.Actor{
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		MiddleName: request.MiddleName,
		Sex:        request.Sex,
		Birthday:   request.Birthday.ToTime(),
	}

	err := s.actorRepo.Create(ctx, &actor)
	if err != nil {
		return models.Actor{}, fmt.Errorf("error creating actor: %v", err)
	}

	return actor, nil
}

func (s *Service) UpdateActor(
	ctx context.Context, id uint, request schemas.UpdateActorRequest,
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

	err := s.actorRepo.Update(ctx, id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) PartialUpdateActor(
	ctx context.Context, id uint, request schemas.PartialUpdateActorRequest,
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

	err := s.actorRepo.Update(ctx, id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveActor(ctx context.Context, id uint) error {
	return s.actorRepo.Remove(ctx, id)
}

func (s *Service) GetActorsWithFilms(
	ctx context.Context,
) ([]schemas.ActorWithFilmsResponse, error) {
	return s.actorRepo.GetListWithFilms(ctx)
}
