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

type FilmService interface {
	AddFilm(context.Context, schemas.AddFilmRequest) (models.Film, error)
	UpdateFilm(context.Context, uint, schemas.UpdateFilmRequest) error
	PartialUpdateFilm(context.Context, uint, schemas.PartialUpdateFilmRequest) error
	RemoveFilm(context.Context, uint) error
	GetFilmsWithActors(context.Context, string, string) ([]schemas.FilmWithActorsResponse, error)
}

type FilmsHandler struct {
	service  FilmService
	validate *validator.Validate
}

func NewFilmsHandler(service FilmService, validate *validator.Validate) *FilmsHandler {
	return &FilmsHandler{
		service:  service,
		validate: validate,
	}
}

// Add godoc
//
//	@Summary		Add film
//	@Description	Add film to database
//	@Security		BasicAuth
//	@Tags			films
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		schemas.AddFilmRequest	true	"New film"
//	@Success		201		{object}	schemas.ActorInfo
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		403		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/films [post]
func (h *FilmsHandler) Add() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var schema schemas.AddFilmRequest
		err := validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		film, err := h.service.AddFilm(r.Context(), schema)
		if err != nil {
			var notFoundErr *postgresql.ErrRecordNotFound
			if errors.As(err, &notFoundErr) {
				resp := schemas.ErrorResponse{Error: "invalid actor id"}
				_ = writeJson(w, resp, http.StatusBadRequest)
				return
			}
			internalError(w)
			return
		}

		resp := schemas.NewFilmInfo(film)
		err = writeJson(w, resp, http.StatusCreated)
		if err != nil {
			internalError(w)
			return
		}
	})
}

// Update godoc
//
//	@Summary		Update film
//	@Description	Update film
//	@Security		BasicAuth
//	@Tags			films
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		schemas.UpdateFilmRequest	true	"Update film"
//	@Param			id		path		int							true	"Film id"
//	@Success		204		{object}	nil
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		403		{object}	schemas.ErrorResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/films/{id} [put]
func (h *FilmsHandler) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			resp := schemas.ErrorResponse{Error: "invalid path parameter: id"}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		var schema schemas.UpdateFilmRequest
		err = validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		err = h.service.UpdateFilm(r.Context(), uint(id), schema)
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
//	@Summary		Partial update film
//	@Description	Partial update film
//	@Security		BasicAuth
//	@Tags			films
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		schemas.PartialUpdateFilmRequest	true	"Update film"
//	@Param			id		path		int									true	"Film id"
//	@Success		204		{object}	nil
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		403		{object}	schemas.ErrorResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/films/{id} [patch]
func (h *FilmsHandler) PartialUpdate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			resp := schemas.ErrorResponse{Error: "invalid path parameter: id"}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		var schema schemas.PartialUpdateFilmRequest
		err = validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		err = h.service.PartialUpdateFilm(r.Context(), uint(id), schema)
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
//	@Summary		Remove film
//	@Description	Remove film from database
//	@Security		BasicAuth
//	@Tags			films
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Film id"
//	@Success		204	{object}	nil
//	@Failure		401	{object}	schemas.ErrorResponse
//	@Failure		403	{object}	schemas.ErrorResponse
//	@Failure		404	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/v1/films/{id} [delete]
func (h *FilmsHandler) Remove() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			resp := schemas.ErrorResponse{Error: "invalid path parameter: id"}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return

		}
		err = h.service.RemoveFilm(r.Context(), uint(id))
		if err != nil {
			var notFoundErr *postgresql.ErrRecordNotFound
			if errors.As(err, &notFoundErr) {
				resp := schemas.ErrorResponse{Error: "film not found"}
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
//	@Summary		List films
//	@Description	Get list of films
//	@Security		BasicAuth
//	@Tags			films
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string	false	"search by films title and actors names"
//	@Param			sortBy	query		string	false	"sorting by field. Format: orderBy=field1,-field2"
//	@Success		200		{array}		schemas.FilmWithActorsResponse
//	@Failure		401		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/films [get]
func (h *FilmsHandler) GetList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		films, err := h.service.GetFilmsWithActors(
			r.Context(),
			r.URL.Query().Get("search"),
			r.URL.Query().Get("sortBy"),
		)
		if err != nil {
			internalError(w)
			return
		}

		err = writeJson(w, films, http.StatusOK)
		if err != nil {
			internalError(w)
			return
		}
	})
}
