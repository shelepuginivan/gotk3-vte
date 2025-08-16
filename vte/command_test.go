package vte_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/shelepuginivan/gotk3-vte/vte"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := vte.CommandNew([]string{"/usr/bin/zsh", "-c", "echo something"})

	cwd, _ := os.Getwd()

	// Test default values
	assert.Equal(t, []string{"/usr/bin/zsh", "-c", "echo something"}, cmd.Args)
	assert.Equal(t, os.Environ(), cmd.Env)
	assert.Equal(t, cwd, cmd.Dir)
	assert.Equal(t, vte.SPAWN_DEFAULT, cmd.SpawnFlags)
	assert.Equal(t, vte.PTY_DEFAULT, cmd.PtyFlags)
	assert.Equal(t, vte.DEFAULT_TIMEOUT, cmd.Timeout)
	assert.NotNil(t, cmd.Cancellable)
	assert.Nil(t, cmd.OnSpawn)
}

func TestCommandWithEnv(t *testing.T) {
	cmd := vte.CommandNew(
		nil,
		vte.CommandWithEnv("X_VTE_TEST", "TestCommandWithEnv"),
		vte.CommandWithEnv("X_VTE_ANOTHER", "something"),
	)

	expected := []string{
		"X_VTE_TEST=TestCommandWithEnv",
		"X_VTE_ANOTHER=something",
	}

	assert.Equal(t, expected, cmd.Env)
}

func TestCommandWithWorkdir(t *testing.T) {
	cmd := vte.CommandNew(
		nil,
		vte.CommandWithWorkdir("/tmp"),
	)

	assert.Equal(t, "/tmp", cmd.Dir)
}

func TestCommandWithSpawnFlags(t *testing.T) {
	cmd := vte.CommandNew(
		nil,
		vte.CommandWithSpawnFlags(vte.SPAWN_CHILD_INHERITS_STDERR),
		vte.CommandWithSpawnFlags(vte.SPAWN_DO_NOT_REAP_CHILD),
		vte.CommandWithSpawnFlags(vte.SPAWN_CLOEXEC_PIPES|vte.SPAWN_STDOUT_TO_DEV_NULL),
	)

	expected := vte.SPAWN_DEFAULT |
		vte.SPAWN_CHILD_INHERITS_STDERR |
		vte.SPAWN_DO_NOT_REAP_CHILD |
		vte.SPAWN_CLOEXEC_PIPES |
		vte.SPAWN_STDOUT_TO_DEV_NULL

	assert.Equal(t, expected, cmd.SpawnFlags)
}

func TestCommandWithPtyFlags(t *testing.T) {
	cmd := vte.CommandNew(
		nil,
		vte.CommandWithPtyFlags(vte.PTY_NO_CTTY),
		vte.CommandWithPtyFlags(vte.PTY_NO_SESSION),
		vte.CommandWithPtyFlags(vte.PTY_NO_UTMP|vte.PTY_NO_WTMP),
	)

	expected := vte.PTY_DEFAULT |
		vte.PTY_NO_CTTY |
		vte.PTY_NO_SESSION |
		vte.PTY_NO_UTMP |
		vte.PTY_NO_WTMP

	assert.Equal(t, expected, cmd.PtyFlags)
}

func TestCommandWithTimeout(t *testing.T) {
	cmd := vte.CommandNew(
		nil,
		vte.CommandWithTimeout(time.Hour),
	)

	assert.Equal(t, time.Hour, cmd.Timeout)
}

func TestCommandWithCancellable(t *testing.T) {
	cancellable, _ := glib.CancellableNew()

	cmd := vte.CommandNew(
		nil,
		vte.CommandWithCancellable(cancellable),
	)

	assert.Equal(t, cancellable, cmd.Cancellable)

}

func TestCommandWithOnSpawn(t *testing.T) {
	cmd := vte.CommandNew(
		nil,
		vte.CommandWithOnSpawn(func(pid int, err error) {
			fmt.Println(pid, err)
		}),
	)

	assert.NotNil(t, cmd.OnSpawn)
}
