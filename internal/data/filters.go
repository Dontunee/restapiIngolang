package data

import (
	"greenlight.alexedwards.net/internal/validator"
	"math"
	"strings"
)

// Metadata Define a new Metadata struct for holding the pagination metadata
type Metadata struct {
 CurrentPage int `json:"current_page,omitempty"`
 PageSize int    `json:"page_size,omitempty"`
 FirstPage int `json:"first_page,omitempty"`
 LastPage int `json:"last_page,omitempty"`
 TotalRecords int `json:"total_records,omitempty"`
}

//The calculateMetadata() function calculates the appropriate metadata values
// given the total number of records, current page, and page size values. Note
// that the last page is calculated using the math.ceil() function, which rounds
// up a float to the nearest integer . eg if there are 12 records in total and
// a page size of 5 , the last page would be math.ceil(12/5) = 3.
func calculateMetadata(totalRecords, page, pageSize int) Metadata{
	if totalRecords == 0 {
		//return an empty meta data struct if there are no records
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords)/ float64(pageSize))),
		TotalRecords: totalRecords,
	}
}




type Filters struct {
	Page int
	PageSize int
	Sort string
	SortSafeList []string
}

// Allows us to set the maximum number of records that a SQL query should return
func (filters Filters) limit() int{
	return filters.PageSize
}

//Allows us to skip a specific number of rows before starting to return records from the query
func (filters Filters) offset() int {
	return (filters.Page - 1) * filters.PageSize
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


