package utils

import (
	"strconv"
	"strings"
)

// MultiError multiple errors
type MultiError []error

// Errors merge multiple errors.
func Errors(errs []error) error {
	count := len(errs)
	if count == 0 {
		return nil
	}
	multiError := make(MultiError, 0, count)
	for _, err := range errs {
		if err == nil {
			continue
		}
		multiError = append(multiError, err)
	}
	if len(multiError) == 0 {
		return nil
	}
	return multiError
}

func (m MultiError) Error() string {
	var errMsg = "MultiError:\n"
	for i, err := range m {
		errMsg += strconv.Itoa(i+1) + ". " + strings.Trim(err.Error(), "\n") + "\n"
	}
	return errMsg
}
