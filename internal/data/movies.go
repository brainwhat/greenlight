package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"greenlight.brainwhat/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"` // This just changes output names
	CreatedAt time.Time `json:"-"`  // "-" doen't show field in json response
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"` // omitempty doesn't show field if it's not defined/zero/""/false/etc
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	stmt := `INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	// pq.Array returns pq.StringArray type that implements the driver.Valuer and sql.Scanner interfaces
	// That are neccessary to translate []string to postgres text[] array
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(stmt, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	// Not necessary, just to not make db call
	if id < 0 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `UPDATE movies
	SET title=$1, year=$2, runtime=$3, genres=$4, version = version +1
	where id=$5
	RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "cannot be empty")
	v.Check(len(movie.Title) < 500, "title", "must be under 500 characters")

	v.Check(movie.Year != 0, "year", "cannpt be empty")
	v.Check(movie.Year > 1888 && movie.Year <= int32(time.Now().Year()), "year", "must be between 1888 and today")

	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "cannot be empty")
	v.Check(validator.CheckForEmptyStrings(movie.Genres), "genres", "cannot be empty")
	v.Check(len(movie.Genres) > 0 && len(movie.Genres) <= 5, "genres", "must have between 1 and 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must be unique")
}
