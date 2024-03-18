package models

import "time"

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type Actor struct {
	ID         uint
	FirstName  string
	LastName   string
	MiddleName *string
	Sex        Sex
	Birthday   time.Time
}

type Film struct {
	ID          uint
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      uint8
}
