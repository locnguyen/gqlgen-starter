package directives

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDirectiveRoot(t *testing.T) {
	root := NewDirectiveRoot()
	assert.NotNil(t, root, "directive root should not be nil")
}
