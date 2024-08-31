package data

import "time"

type Movie struct {
	ID 					int     	`json:"id"`
	Title 			string 		`json:"title"`
	Year 				int32			`json:"year,omitempty"`
	Runtime 		Runtime		`json:"runtime,omitempty"`
	Genres 			[]string  `json:"genres,omitempty"`
	Version     int32     `json:"version"`
	Created_At  time.Time	`json:"-"`
}