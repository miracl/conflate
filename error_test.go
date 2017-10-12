package conflate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeError_NoArgs(t *testing.T) {
	err := makeError("err msg")
	assert.Equal(t, err.Error(), "err msg")
}

func TestMakeError_WithArgs(t *testing.T) {
	err := makeError("err %v", "msg")
	assert.Equal(t, err.Error(), "err msg")
}

func TestMakeContextError_NoArgs(t *testing.T) {
	err := makeContextError("ctx", "err msg")
	assert.Equal(t, err.Error(), "err msg (ctx)")
}

func TestMakeContextError_WithArgs(t *testing.T) {
	err := makeContextError("ctx", "err %v", "msg")
	assert.Equal(t, err.Error(), "err msg (ctx)")
}

func TestWrapError_Nil(t *testing.T) {
	err := wrapError(nil, "an error")
	assert.Nil(t, err)
}

func TestWrapError_NoArgs(t *testing.T) {
	err := wrapError(makeError("err2"), "err1")
	assert.Equal(t, err.Error(), "err1 : err2")
}

func TestWrapError_WithArgs(t *testing.T) {
	err := wrapError(makeError("err2"), "err%v", 1)
	assert.Equal(t, err.Error(), "err1 : err2")
}
