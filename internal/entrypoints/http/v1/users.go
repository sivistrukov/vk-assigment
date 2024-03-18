package v1

import (
	"context"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

type UserService interface {
	CreateUser(context.Context, schemas.CreateUserRequest) (models.User, error)
}

type UsersHandler struct {
	service  UserService
	validate *validator.Validate
}

func NewUsersHandler(service UserService, validate *validator.Validate) *UsersHandler {
	return &UsersHandler{
		service:  service,
		validate: validate,
	}
}

// Create godoc
//
//	@Summary		Create user
//	@Description	Create new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		schemas.CreateUserRequest	true "New user"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/v1/users [post]
func (h *UsersHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var schema schemas.CreateUserRequest
		err := validateRequestBody(r, h.validate, &schema)
		if err != nil {
			resp := schemas.ErrorResponse{Error: err.Error()}
			_ = writeJson(w, resp, http.StatusBadRequest)
			return
		}

		user, err := h.service.CreateUser(r.Context(), schema)
		if err != nil {
			internalError(w)
			return
		}

		err = writeJson(w, user, http.StatusCreated)
		if err != nil {
			internalError(w)
			return
		}
	})
}
