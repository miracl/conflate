package conflate

import (
	"fmt"
)

type context string

func makeError(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func makeContextError(ctx context, msg string, args ...interface{}) error {
	return makeError("%v (%v)", makeError(msg, args...), ctx)
}

func wrapError(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return makeError("%v : %v", makeError(msg, args...), err)
}

func detailError(err error, msg string, args ...interface{}) error {
	return makeError("%v. %v", err, makeError(msg, args...))
}
