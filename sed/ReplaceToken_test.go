package sed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveReplaceToken(t *testing.T) {
	assert.Equal(t, ReplaceToken{Source: "", Target: ""}, ResolveReplaceToken(">>>"))
	assert.Equal(t, ReplaceToken{Source: "a", Target: ""}, ResolveReplaceToken("a>>>"))
	assert.Equal(t, ReplaceToken{Source: "", Target: "a"}, ResolveReplaceToken(">>>a"))
	assert.Equal(t, ReplaceToken{Source: "a", Target: "b"}, ResolveReplaceToken("a>>>b"))
}
