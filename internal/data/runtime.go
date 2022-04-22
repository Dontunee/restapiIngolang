package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidRunTimeFormat Define an error that our UnmarshalJSON() can return if we are
// unable to parse the JSON string successfully
var ErrInvalidRunTimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (runtime Runtime) MarshalJSON() ([]byte, error) {
	//Generate a string containing the movie runtime in the required format
	jsonValue := fmt.Sprintf("%d mins", runtime)

	//Use the strconv.Quote() function on the string to wrap it in double quotes. it needs
	//to be surrounded by double quotes in order to be a valid *JSON string*.

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// UnMarshalJSON Implement a UnMarshalJSON() method on the RuntimeType so that it satisfies the
// json.UnMarshaler interface. IMPORTANT: Because UnMarshalJSON() needs to modify the receiver
// (our runtime type), we must use a pointer for this to work correctly.
// Otherwise , we will be modifying a copy (which is then discarded when this method returns).
func (runtime *Runtime) UnMarshalJSON(jsonValue []byte) error {

	//we expect that the incoming JSON value will be a string in the format
	// "<runtime> mins" , and the first thing we need to do is remove the
	//surrounding double-quotes from the string. if we cant unquote it, then we return the
	//ErrInvalidRunTimeFormat error
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRunTimeFormat
	}

	//Split the string to isolate the part containing the number
	parts := strings.Split(unquotedJSONValue, " ")

	//Sanity check the parts of the string to make sure it was in the expected format
	//if it isn't, we return the ErrInvalidRunTimeFormat error again
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRunTimeFormat
	}

	//Otherwise, parse the string containing the number into an int32, Again if this
	//fails return the ErrInvalidRunTimeFormat error again
	number, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRunTimeFormat
	}

	//Convert the int32 to a Runtime type and assign this to the receiver. Note that we use
	//the * operator to deference the receiver (which is a pointer to a Runtime type)
	//in order to set the underlying value of the pointer
	*runtime = Runtime(number)

	return nil
}
