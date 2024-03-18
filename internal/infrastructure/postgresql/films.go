package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
	"github.com/sivistrukov/vk-assigment/internal/services/text"
)

type FilmRepo struct {
	db *sql.DB
}

func NewFilmRepo(db *sql.DB) *FilmRepo {
	return &FilmRepo{
		db: db,
	}
}

func (r *FilmRepo) Create(
	_ context.Context, film *models.Film, actorsIds ...uint,
) error {
	var err error
	tx, _ := r.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	stmt := `
	INSERT INTO films (title, description, release_date, rating) 
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`
	err = tx.QueryRow(
		stmt,
		film.Title,
		film.Description,
		film.ReleaseDate,
		film.Rating,
	).Scan(&film.ID)
	if err != nil {
		return err
	}

	for _, actorId := range actorsIds {
		stmt = `
		INSERT INTO actors_and_films (film_id, actor_id) 
		VALUES ($1, $2);
		`
		_, err = tx.Exec(stmt, film.ID, actorId)
		if err != nil {
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				return &ErrRecordNotFound{
					tableName: "actors",
					identity:  fmt.Sprintf("%d", actorId),
				}
			}
			return err
		}
	}

	return nil
}

func (r *FilmRepo) Update(
	ctx context.Context, id uint, updates map[string]any,
) error {
	var err error
	tx, _ := r.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	err = r.updateFilm(ctx, tx, id, updates)
	if err != nil {
		return err
	}

	if value, ok := updates["actors_ids"]; ok &&
		reflect.TypeOf(value).Kind() == reflect.Slice {

		actorsIds := value.([]uint)
		err = r.updateFilmActors(ctx, tx, id, actorsIds...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FilmRepo) updateFilm(
	_ context.Context, tx *sql.Tx, id uint, updates map[string]any,
) error {
	builder := strings.Builder{}
	builder.WriteString("UPDATE films SET ")

	values := make([]any, 0, len(updates)+1)
	i := 0
	for field, value := range updates {
		if field == "actors_ids" {
			continue
		}
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

	result, err := tx.Exec(stmt, values...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &ErrRecordNotFound{
			tableName: "films",
			identity:  fmt.Sprintf("%d", id),
		}
	}

	return nil
}

func (r *FilmRepo) updateFilmActors(
	_ context.Context, tx *sql.Tx, filmId uint, actorsIds ...uint,
) error {
	stmt := `
	SELECT actor_id FROM actors_and_films
	WHERE film_id = $1
	`
	rows, err := tx.Query(stmt, filmId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var oldIds []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return err
		}

		oldIds = append(oldIds, id)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	var newIds, deletedIds []uint
	for _, v := range actorsIds {
		if !slices.Contains(oldIds, v) {
			newIds = append(newIds, v)
		}
	}

	for _, v := range oldIds {
		if !slices.Contains(actorsIds, v) {
			deletedIds = append(deletedIds, v)
		}
	}

	for _, v := range newIds {
		stmt = `
		INSERT INTO actors_and_films (film_id, actor_id)
		VALUES ($1, $2);
		`

		_, err = tx.Exec(stmt, filmId, v)
		if err != nil {
			return err
		}
	}

	for _, v := range deletedIds {
		stmt = `
        DELETE FROM actors_and_films
        WHERE film_id = $1 AND actor_id = $2;
        `

		_, err = tx.Exec(stmt, filmId, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FilmRepo) Remove(_ context.Context, id uint) error {
	stmt, err := r.db.Prepare(`
	DELETE FROM films WHERE id = $1
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
			tableName: "films",
			identity:  fmt.Sprintf("%d", id),
		}
	}

	return nil
}

func (r *FilmRepo) GetFilmsWithActors(
	ctx context.Context, search string, sortBy string,
) ([]schemas.FilmWithActorsResponse, error) {
	builder := strings.Builder{}
	builder.WriteString(`
	SELECT films.id, films.title, films.description, films.release_date, films.rating
	FROM films
	INNER JOIN actors_and_films AS aaf ON films.id = aaf.film_id
	INNER JOIN actors ON aaf.actor_id = actors.id
	`)

	if len(search) > 0 {
		builder.WriteString("WHERE films.title ILIKE '%" + search + "%'")
		builder.WriteString("OR actors.first_name ILIKE '%" + search + "%'")
		builder.WriteString("OR actors.last_name ILIKE '%" + search + "%'")
		builder.WriteString("OR actors.middle_name ILIKE '%" + search + "%'")
	}

	builder.WriteString("GROUP BY films.id")

	if len(sortBy) > 0 {
		builder.WriteString(" ORDER BY ")
		params := strings.Split(sortBy, ",")
		for i, v := range params {
			if i > 0 {
				builder.WriteString(", ")
			}

			sortParam := text.CamelToSnake(v)
			if sortParam[0] == '-' {
				builder.WriteString(fmt.Sprintf("%s DESC", sortParam[1:]))
			} else {
				builder.WriteString(fmt.Sprintf("%s ASC", sortParam))
			}
		}
	} else {
		builder.WriteString(" ORDER BY rating DESC")
	}

	stmt := builder.String()
	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	films := make([]schemas.FilmWithActorsResponse, 0)
	for rows.Next() {
		var film schemas.FilmWithActorsResponse
		var date time.Time
		err = rows.Scan(
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

		stmt = `
		SELECT actors.id, actors.first_name, actors.last_name, actors.middle_name, actors.sex, actors.birthday
		FROM actors 
		INNER JOIN actors_and_films ON actors.id = actors_and_films.actor_id
		WHERE actors_and_films.film_id = $1
		`
		rows2, err := r.db.Query(stmt, film.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		actors := make([]schemas.ActorInfo, 0)
		for rows2.Next() {
			var date time.Time
			var actor schemas.ActorInfo
			err = rows2.Scan(
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

			actors = append(actors, actor)
		}

		film.Actors = actors

		films = append(films, film)
	}

	return films, nil
}
