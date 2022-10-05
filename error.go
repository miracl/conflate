package conflate

import (
	"fmt"
)

type context string

type errWithContext struct {
	msg     string
	context context
}

func (e errWithContext) Error() string {
	return fmt.Sprintf("%v (%v)", e.msg, e.context)
}
