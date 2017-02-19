package multipartbuilder

import (
	"errors"
	"strings"
)

func multiError(errs []error) error {
	n := len(errs)
	if n == 0 {
		return nil
	}

	ss := make([]string, n)
	for i, err := range errs {
		ss[i] = err.Error()
	}
	return errors.New(strings.Join(ss, "; "))
}
