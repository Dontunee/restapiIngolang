package data

import (
	"greenlight.alexedwards.net/internal/validator"
	"strings"
)

type Filters struct {
	Page int
	PageSize int
	Sort string
	SortSafeList []string
}

// Check that the client-provided sort field matches one of the entries in our safe list
// and if it does, extract the column name from the sort field by stripping the leading
// hyphen character (if one exists)
func (filters Filters) sortColumn() string {
	for _, safeValue := range filters.SortSafeList {
		if filters.Sort == safeValue {
			return strings.TrimPrefix(filters.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + filters.Sort)
}


// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field
func (filters Filters) sortDirection() string {
	if strings.HasPrefix(filters.Sort, "-") {
		return "DESC"
	}

	return "ASC"
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


