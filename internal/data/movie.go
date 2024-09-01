package data

import (
	"database/sql"
	"errors"
	"time"

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
