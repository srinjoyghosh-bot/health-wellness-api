package utils

import (
	"strconv"
)

// ParseUint converts a string to a uint, returning an error if the string is not a valid unsigned integer.
func ParseUint(s string) (uint, error) {
	parsed, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, NewBadRequestError("Invalid unsigned integer format")
	}
	return uint(parsed), nil
}
