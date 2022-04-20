package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (runtime Runtime) MarshalJSON() ([]byte, error) {
	//Generate a string containing the movie runtime in the required format
	jsonValue := fmt.Sprintf("%d mins", runtime)

	//Use the strconv.Quote() function on the string to wrap it in double quotes. it needs
	//to be surrounded by double quotes in order to be a valid *JSON string*.

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}