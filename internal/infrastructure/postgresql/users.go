package postgresql

import (
	"context"
	"database/sql"

	"github.com/sivistrukov/vk-assigment/internal/models"
	"github.com/sivistrukov/vk-assigment/internal/services/auth"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(_ context.Context, user *models.User) error {
	stmt, err := r.db.Prepare(`
	INSERT INTO users (username, password, is_admin) 
	VALUES ($1, $2, $3)
	RETURNING id;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	password, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = stmt.QueryRow(
		user.Username,
		password,
		user.IsAdmin,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetByUsername(_ context.Context, username string) (models.User, error) {
	stmt := `
	SELECT id, username, password, is_admin 
	FROM users WHERE username = $1
	`
	row := r.db.QueryRow(stmt, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, &ErrRecordNotFound{
				tableName: "users",
				identity:  username,
			}
		}
		return user, err
	}

	return user, nil
}
