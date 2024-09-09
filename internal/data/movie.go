package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gaba-bouliva/movie-api/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID 					int     	`json:"id"`
	Title 			string 		`json:"title"`
	Year 				int32			`json:"year,omitempty"`
	Runtime 		Runtime		`json:"runtime,omitempty"`
	Genres 			[]string  `json:"genres,omitempty"`
	Version     int32     `json:"version"`
	Created_At  time.Time	`json:"-"`
}

type MovieModel struct {
	DB 	*sql.DB
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime.Value != 0, "runtime", "must be provided")
	v.Check(movie.Runtime.Value > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

func (m MovieModel) Get(id int) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1
	`
	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.Runtime.Value,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch{
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Create(movie *Movie) error {
	query := `
	INSERT INTO movies (title, year, runtime, genres) 
	VALUES ($1, $2, $3, $4)
	RETURNING id,version
	`
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime.Value,
		pq.Array(movie.Genres),
	}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.Version)
} 

func (m MovieModel) Update(movie *Movie) error {
	query := `
	UPDATE movies
	SET	title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5
	RETURNING version
	`

	args := []any{movie.Title, movie.Year, movie.Runtime.Value, pq.Array(movie.Genres), movie.ID}

	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int) error{
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM movies
	WHERE id = $1
	`

	result, err := m.DB.Exec(query,  id)
	if err != nil {
		return ErrRecordNotFound	
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrRecordNotFound
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *MovieModel) GetAll(title string, genres[]string, filters Filters) ([]*Movie, error) {
	query := `
	SELECT id, title, year, runtime, genres, version
	FROM movies
	ORDER BY id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	movies := []*Movie{}

	var movie Movie

	for rows.Next() {
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Year,
			&movie.Runtime.Value,
			pq.Array(&movie.Genres),
			&movie.Version )
		
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies,nil
}