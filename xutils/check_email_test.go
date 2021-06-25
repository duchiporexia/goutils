package xutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEmailValid(t *testing.T) {
	assert.False(t, IsEmailValid("xvv"))
	assert.False(t, IsEmailValid("xvv@g"))
	assert.False(t, IsEmailValid("xvv@g."))
	assert.True(t, IsEmailValid("xvv@g.cc"))
	assert.True(t, IsEmailValid("xvv@gmail.com"))
}
