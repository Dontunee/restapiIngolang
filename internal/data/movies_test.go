package data

import (
	"greenlight.alexedwards.net/internal/validator"
	"testing"
)

func Test_ValidateMovie_Should_Make_Validator_Invalid_When_Movie_Input_Is_Invalid(t *testing.T) {
	//Arrange
	v := validator.New()
	invalidInput := &Movie{
		Title:  "",
		Year:   1000,
		Genres: []string{"sci", "sci"},
	}

	//Act
	ValidateMovie(v, invalidInput)

	//Assert
	if v.Valid() {
		t.Error("Invalid movie validation while testing for input validation")
	}
}
