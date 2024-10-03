package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "string", messageType[string]())
	assert.Equal(t, "int", messageType[int]())
	assert.Equal(t, "int", messageType[*int]())
	assert.Equal(t, "int", messageType[**int]())
	assert.Equal(t, "struct {}", messageType[struct{}]())
	assert.Equal(t, "struct {}", messageType[*struct{}]())
	assert.Equal(t, "error", messageType[error]())
	assert.Equal(t, "error", messageType[*error]())
	assert.Equal(t, "kafka.S", messageType[S]())
	assert.Equal(t, "kafka.S", messageType[*S]())
	assert.Equal(t, "kafka.G[int]", messageType[G[int]]())
	assert.Equal(t, "kafka.G[int]", messageType[*G[int]]())
	assert.Equal(t, "kafka.I", messageType[I]())
	assert.Equal(t, "kafka.I", messageType[*I]())
}

type S struct{}
type G[T any] struct{}
type I interface{}
