package data

import "greenlight.alexedwards.net/internal/validator"

type Filters struct {
	Page int
	PageSize int
	Sort string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, filters Filters){
	//Check that the page and page_size parameters contain sensible values
	v.Check(filters.Page > 0 , "page", "must be greater than zero")
	v.Check(filters.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(filters.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(filters.PageSize <= 100, "page_size", "must be a maximum of 100")

	//Check that the sort parameter matches a value in the safe list
	v.Check(validator.In(filters.Sort,filters.SortSafeList...), "sort", "invalid sort value")
}


