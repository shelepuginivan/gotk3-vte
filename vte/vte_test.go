package vte_test

import (
	"bytes"
	"os/exec"
	"os/user"
	"testing"

	"github.com/shelepuginivan/gotk3-vte/vte"
	"github.com/stretchr/testify/assert"
)

func TestGetUserShell(t *testing.T) {
	actual := vte.GetUserShell()

	user, err := user.Current()
	if err != nil {
		assert.Equal(t, "/bin/sh", actual)
		return
	}

	out, err := exec.Command("getent", "passwd", user.Uid).Output()
	if err != nil {
		assert.Equal(t, "/bin/sh", actual)
		return
	}

	expected := string(bytes.Split(bytes.TrimSpace(out), []byte{':'})[6])
	assert.Equal(t, expected, actual)
}
