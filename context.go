package conflate

import (
	"fmt"
	"path"
)

func rootContext() context {
	return context("#")
}

func (c context) String() string {
	return string(c)
}

func (c context) add(s ...string) context {
	return context(path.Join(c.String(), path.Join(s...)))
}

func (c context) addInt(i int) context {
	return context(fmt.Sprintf("%v[%v]", c.String(), i))
}
