package vte

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrCgoCall(t *testing.T) {
	err := errCgoCall{
		Function: "some_func",
		Detail:   "lorem ipsum",
	}

	assert.Equal(t, "vte: some_func() lorem ipsum", err.Error())
}

func Test_errFailed(t *testing.T) {
	err := errFailed("some_func")
	assert.Equal(t, "vte: some_func() failed", err.Error())
}

func Test_errNilPointer(t *testing.T) {
	err := errNilPointer("some_func")
	assert.Equal(t, "vte: some_func() returned nil pointer", err.Error())
}
