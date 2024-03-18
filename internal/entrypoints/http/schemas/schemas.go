package schemas

import (
	"github.com/sivistrukov/vk-assigment/internal/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required" example:"user"`
	Password string `json:"password" validate:"required" example:"<PASSWORD>"`
	IsAdmin  bool   `json:"isAdmin" example:"false"`
}

type AddActorRequest struct {
	FirstName  string     `json:"firstName" validate:"required"`
	LastName   string     `json:"lastName" validate:"required"`
	MiddleName *string    `json:"middleName,omitempty"`
	Sex        models.Sex `json:"sex" validate:"required,sexValidation"`
	Birthday   Date       `json:"birthday" validate:"required,dateValidation" example:"02-01-2006"`
}

type UpdateActorRequest struct {
	FirstName  string     `json:"firstName" validate:"required"`
	LastName   string     `json:"lastName" validate:"required"`
	MiddleName *string    `json:"middleName"`
	Sex        models.Sex `json:"sex" validate:"required,sexValidation"`
	Birthday   Date       `json:"birthday" validate:"required,dateValidation" example:"02-01-2006"`
}

type PartialUpdateActorRequest struct {
	FirstName  *string     `json:"firstName" validate:"omitempty"`
	LastName   *string     `json:"lastName" validate:"omitempty"`
	MiddleName *string     `json:"middleName" validate:"omitempty"`
	Sex        *models.Sex `json:"sex" validate:"omitempty,sexValidation"`
	Birthday   *Date       `json:"birthday" validate:"omitempty,dateValidation" example:"02-01-2006"`
}

type AddFilmRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"max=1000"`
	ReleaseDate Date   `json:"releaseDate" validate:"required,dateValidation" example:"02-01-2006"`
	Rating      uint8  `json:"rating" validate:"min=0,max=10"`
	ActorsIDs   []uint `json:"actorsIds" validate:"required"`
}

type UpdateFilmRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"required,max=1000"`
	ReleaseDate Date   `json:"releaseDate" validate:"required,dateValidation" example:"02-01-2006"`
	Rating      uint8  `json:"rating" validate:"min=0,max=10"`
	ActorsIds   []uint `json:"actorsIds" validate:"required"`
}

type PartialUpdateFilmRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=150"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	ReleaseDate *Date   `json:"releaseDate" validate:"omitempty,dateValidation" example:"02-01-2006"`
	Rating      *uint8  `json:"rating" validate:"omitempty,min=0,max=10"`
	ActorsIDs   *[]uint `json:"actorsIds" validate:"omitempty"`
}

type ActorWithFilmsResponse struct {
	ID         uint       `json:"id"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	MiddleName *string    `json:"middleName"`
	Sex        models.Sex `json:"sex"`
	Birthday   Date       `json:"birthday" example:"02-01-2006"`
	Films      []FilmInfo `json:"films"`
}

type FilmInfo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate Date   `json:"releaseDate" example:"02-01-2006"`
	Rating      uint8  `json:"rating"`
}

func NewFilmInfo(film models.Film) FilmInfo {
	return FilmInfo{
		ID:          film.ID,
		Title:       film.Title,
		Description: film.Description,
		ReleaseDate: NewDate(film.ReleaseDate),
		Rating:      film.Rating,
	}
}

type FilmWithActorsResponse struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ReleaseDate Date        `json:"releaseDate" example:"02-01-2006"`
	Rating      uint8       `json:"rating"`
	Actors      []ActorInfo `json:"actors"`
}

type ActorInfo struct {
	ID         uint       `json:"id"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	MiddleName *string    `json:"middleName"`
	Sex        models.Sex `json:"sex"`
	Birthday   Date       `json:"birthday" example:"02-01-2006"`
}

func NewActorInfo(actor models.Actor) ActorInfo {
	return ActorInfo{
		ID:         actor.ID,
		FirstName:  actor.FirstName,
		LastName:   actor.LastName,
		MiddleName: actor.MiddleName,
		Sex:        actor.Sex,
		Birthday:   NewDate(actor.Birthday),
	}
}
