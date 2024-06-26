package conflate

import (
	"fmt"
)

type context string

type contextError struct {
	msg     string
	context context
}

func (e contextError) Error() string {
	return fmt.Sprintf("%v (%v)", e.msg, e.context)
}
