package sed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveReplaceToken(t *testing.T) {
	assert.Equal(t, ReplaceToken{Source: "", Target: ""}, ResolveReplaceToken(">>>"))
	assert.Equal(t, ReplaceToken{Source: "a", Target: ""}, ResolveReplaceToken("a>>>"))
	assert.Equal(t, ReplaceToken{Source: "", Target: "a"}, ResolveReplaceToken(">>>a"))
	assert.Equal(t, ReplaceToken{Source: "a", Target: "b"}, ResolveReplaceToken("a>>>b"))
}
