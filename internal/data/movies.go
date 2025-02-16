package data

import (
	"database/sql"
	"time"

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
	return nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "title cannot be empty")
	v.Check(len(movie.Title) < 500, "title", "title must be under 500 characters")

	v.Check(movie.Year != 0, "year", "must provide year")
	v.Check(movie.Year > 1888 && movie.Year <= int32(time.Now().Year()), "year", "incorrect year")

	v.Check(movie.Runtime > 0, "runtime", "runtime must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "genres cannot be empty")
	v.Check(validator.CheckForEmptyStrings(movie.Genres), "genres", "genres cannot be empty")
	v.Check(len(movie.Genres) > 0 && len(movie.Genres) <= 5, "genres", "movie must have between 1 and 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "genres must be unique")
}
