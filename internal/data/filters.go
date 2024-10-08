package data

import (
	"github.com/gaba-bouliva/movie-api/internal/validator"
)

type Filters struct {
	Page 					int
	PageSize 			int						 
	Sort 					string
	SortSafeList 	[]string
}

func ValidateFilters(f Filters, v *validator.Validator) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

