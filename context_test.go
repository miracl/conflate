package conflate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootContext_String(t *testing.T) {
	ctx := rootContext()
	assert.Equal(t, "#", ctx.String())
}

func TestMakeContext_Add(t *testing.T) {
	ctx := rootContext()
	ctx2 := ctx.add("parent").add("child")
	assert.Equal(t, "#", ctx.String())
	assert.Equal(t, "#/parent/child", ctx2.String())
}

func TestMakeContext_AddInt(t *testing.T) {
	ctx := rootContext()
	ctx2 := ctx.add("parent").addInt(3)
	assert.Equal(t, "#", ctx.String())
	assert.Equal(t, "#/parent[3]", ctx2.String())
}
