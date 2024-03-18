package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/infrastructure/postgresql"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

type ActorService interface {
	AddActor(context.Context, schemas.AddActorRequest) (models.Actor, error)
	UpdateActor(context.Context, uint, schemas.UpdateActorRequest) error
	PartialUpdateActor(context.Context, uint, schemas.PartialUpdateActorRequest) error
	RemoveActor(context.Context, uint) error
	GetActorsWithFilms(context.Context) ([]schemas.ActorWithFilmsResponse, error)
}

type ActorHandler struct {
	service  ActorService
	validate *validator.Validate
}

func NewActorHandler(service ActorService, validate *validator.Validate) *ActorHandler {
	return &ActorHandler{
		service:  service,
		validate: validate,
	}
}

// Add godoc
//
//	@Summary		Add actor
//	@Description	Add actor to database
//	@Security		BasicAuth
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		schemas.AddActorRequest	true	"New actor"
//	@Success		201		{object}	schemas.ActorInfo
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		403		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/actors [post]
func (h *ActorHandler) Add() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var schema schemas.AddActorRequest
		err := validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		actor, err := h.service.AddActor(r.Context(), schema)
		if err != nil {
			internalError(w)
			return
		}

		resp := schemas.NewActorInfo(actor)
		err = writeJson(w, resp, http.StatusCreated)
		if err != nil {
			internalError(w)
			return
		}
	})
}

// Update godoc
//
//	@Summary		Update actor
//	@Description	Update actor
//	@Security		BasicAuth
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		schemas.UpdateActorRequest	true	"Update actor"
//	@Param			id		path		int							true	"Actor id"
//	@Success		204		{object}	nil
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		403		{object}	schemas.ErrorResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/actors/{id} [put]
func (h *ActorHandler) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			resp := schemas.ErrorResponse{Error: "invalid path parameter: id"}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		var schema schemas.UpdateActorRequest
		err = validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		err = h.service.UpdateActor(r.Context(), uint(id), schema)
		if err != nil {
			var notFoundErr *postgresql.ErrRecordNotFound
			if errors.As(err, &notFoundErr) {
				resp := schemas.ErrorResponse{Error: "actor not found"}
				_ = writeJson(w, resp, http.StatusNotFound)
				return
			}
			internalError(w)
			return
		}

		err = writeJson(w, nil, http.StatusNoContent)
		if err != nil {
			internalError(w)
			return
		}
	})
}

// PartialUpdate godoc
//
//	@Summary		Partial update actor
//	@Description	Partial update actor
//	@Security		BasicAuth
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		schemas.PartialUpdateActorRequest	true	"Update actor"
//	@Param			id		path		int									true	"Actor id"
//	@Success		204		{object}	nil
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		403		{object}	schemas.ErrorResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/actors/{id} [patch]
func (h *ActorHandler) PartialUpdate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			resp := schemas.ErrorResponse{Error: "invalid path parameter: id"}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		var schema schemas.PartialUpdateActorRequest
		err = validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		err = h.service.PartialUpdateActor(r.Context(), uint(id), schema)
		if err != nil {
			var notFoundErr *postgresql.ErrRecordNotFound
			if errors.As(err, &notFoundErr) {
				resp := schemas.ErrorResponse{Error: "actor not found"}
				_ = writeJson(w, resp, http.StatusNotFound)
				return
			}
			internalError(w)
			return
		}

		err = writeJson(w, nil, http.StatusNoContent)
		if err != nil {
			internalError(w)
			return
		}
	})
}

// Remove godoc
//
//	@Summary		Remove actor
//	@Description	Remove actor from database
//	@Security		BasicAuth
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Actor id"
//	@Success		204	{object}	nil
//	@Failure		401	{object}	schemas.ErrorResponse
//	@Failure		403	{object}	schemas.ErrorResponse
//	@Failure		404	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/v1/actors/{id} [delete]
func (h *ActorHandler) Remove() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			resp := schemas.ErrorResponse{Error: "invalid path parameter: id"}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return

		}
		err = h.service.RemoveActor(r.Context(), uint(id))
		if err != nil {
			var notFoundErr *postgresql.ErrRecordNotFound
			if errors.As(err, &notFoundErr) {
				resp := schemas.ErrorResponse{Error: "actor not found"}
				_ = writeJson(w, resp, http.StatusNotFound)
				return
			}
			internalError(w)
			return
		}

		err = writeJson(w, nil, http.StatusNoContent)
		if err != nil {
			internalError(w)
			return
		}
	})
}

// GetList godoc
//
//	@Summary		List actors
//	@Description	Get list of actors
//	@Security		BasicAuth
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		schemas.ActorWithFilmsResponse
//	@Failure		401	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/v1/actors [get]
func (h *ActorHandler) GetList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actors, err := h.service.GetActorsWithFilms(r.Context())
		if err != nil {
			internalError(w)
			return
		}

		err = writeJson(w, actors, http.StatusOK)
		if err != nil {
			internalError(w)
			return
		}
	})
}
