package validator

import "regexp"

// EmailRX Regular expression for sanity checking the format of email address
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator A new Validator type which contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

//New is a helper which creates a new validator instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

//Valid returns true if the errors map doesnt contain any entries
func (validator *Validator) Valid() bool {
	return len(validator.Errors) == 0
}

//AddError adds an error message to the map (so long as no entry already exists for the given key )
func (validator *Validator) AddError(key, message string) {
	if _, exists := validator.Errors[key]; !exists {
		validator.Errors[key] = message
	}
}

//Check adds an error message to the map only if a validation is not 'ok'
func (validator *Validator) Check(ok bool, key, message string) {
	if !ok {
		validator.AddError(key, message)
	}
}

//In returns true of a specific value is in a list of strings.
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

//Matches returns true if a string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

//Unique returns true if all string values in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
