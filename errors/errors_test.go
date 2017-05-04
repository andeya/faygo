package errors

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	errs := []error{
		errors.New("error_text1"),
		errors.New("error_text2"),
		nil,
		errors.New("error_text4"),
		errors.New("error_text5"),
		nil,
		errors.New("error_text7"),
	}
	t.Log(Errors(errs))
	t.Log(Errors(nil))
}
