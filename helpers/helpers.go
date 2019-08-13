package helpers

import (
	"fmt"
	"github.com/springernature/halfpipe/linters/errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func AssertMissingField(t *testing.T, name string, err error) {
	t.Helper()

	mf, ok := err.(errors.MissingFieldError)
	if !ok {
		assert.Fail(t, "error is not a MissingField", err)
	} else {
		assert.Equal(t, name, mf.Name)
	}
}

func AssertMissingFieldInErrors(t *testing.T, name string, errs []error) {
	t.Helper()

	for _, err := range errs {
		mf, ok := err.(errors.MissingFieldError)
		if ok {
			if strings.Contains(mf.Name, name) {
				return
			}
		}
	}
	assert.Fail(t, fmt.Sprintf("Could not find invalid field error for '%s' in %s", name, errs))
}

func AssertInvalidFieldInErrors(t *testing.T, name string, errs []error) {
	t.Helper()

	for _, err := range errs {
		mf, ok := err.(errors.InvalidFieldError)
		if ok {
			if strings.Contains(mf.Name, name) {
				return
			}
		}
	}
	assert.Fail(t, fmt.Sprintf("Could not find invalid field error for '%s' in %s", name, errs))
}

func AssertFileErrorInErrors(t *testing.T, path string, errs []error) {
	t.Helper()

	for _, err := range errs {
		e, ok := err.(errors.FileError)
		if ok {
			if e.Path == path {
				return
			}
		}
	}
	assert.Fail(t, fmt.Sprintf("Could not find FileError for path '%s' in %s", path, errs))

}
