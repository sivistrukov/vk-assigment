package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

type ActorRepo struct {
	db *sql.DB
}

func NewActorRepo(db *sql.DB) *ActorRepo {
	return &ActorRepo{
		db: db,
	}
}

func (r *ActorRepo) Create(_ context.Context, actor *models.Actor) error {
	stmt, err := r.db.Prepare(`
	INSERT INTO actors (first_name, last_name, middle_name, sex, birthday) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		actor.FirstName,
		actor.LastName,
		actor.MiddleName,
		actor.Sex,
		actor.Birthday,
	).Scan(&actor.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ActorRepo) Update(
	_ context.Context, id uint, updates map[string]any,
) error {
	builder := strings.Builder{}
	builder.WriteString("UPDATE actors SET ")

	values := make([]any, 0, len(updates)+1)
	i := 0
	for field, value := range updates {
		if i > 0 {
			builder.WriteString(", ")
		}
		i++

		builder.WriteString(fmt.Sprintf("%s = $%v", field, i))
		values = append(values, value)
	}
	if len(values) == 0 {
		return nil
	}

	builder.WriteString(fmt.Sprintf(" WHERE id = $%v", i+1))
	values = append(values, id)

	stmt := builder.String()

	result, err := r.db.Exec(stmt, values...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &ErrRecordNotFound{
			tableName: "actors",
			identity:  fmt.Sprintf("%d", id),
		}
	}

	return nil
}

func (r *ActorRepo) Remove(_ context.Context, id uint) error {
	stmt, err := r.db.Prepare(`
	DELETE FROM actors WHERE id = $1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &ErrRecordNotFound{
			tableName: "actors",
			identity:  fmt.Sprintf("%d", id),
		}
	}

	return nil
}

func (r *ActorRepo) GetListWithFilms(_ context.Context) ([]schemas.ActorWithFilmsResponse, error) {
	stmt := `
	SELECT id, first_name, last_name, middle_name, sex, birthday 
	FROM actors 
	`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actors := make([]schemas.ActorWithFilmsResponse, 0)
	for rows.Next() {
		var actor schemas.ActorWithFilmsResponse
		var date time.Time
		err = rows.Scan(
			&actor.ID,
			&actor.FirstName,
			&actor.LastName,
			&actor.MiddleName,
			&actor.Sex,
			&date,
		)
		if err != nil {
			return nil, err
		}
		actor.Birthday = schemas.NewDate(date)

		stmt = `
		SELECT films.id, films.title, films.description, films.release_date, films.rating 
		FROM films 
		INNER JOIN actors_and_films ON films.id = actors_and_films.film_id
		WHERE actors_and_films.actor_id = $1
		`
		rows2, err := r.db.Query(stmt, actor.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		films := make([]schemas.FilmInfo, 0)
		for rows2.Next() {
			var date time.Time
			var film schemas.FilmInfo
			err = rows2.Scan(
				&film.ID,
				&film.Title,
				&film.Description,
				&date,
				&film.Rating,
			)
			if err != nil {
				return nil, err
			}
			film.ReleaseDate = schemas.NewDate(date)

			films = append(films, film)
		}

		actor.Films = films
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}
