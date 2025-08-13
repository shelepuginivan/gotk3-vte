package vte

import (
	"fmt"
	"os"
	"time"

	"github.com/gotk3/gotk3/glib"
)

// CommandOption allows to configure [Command].
type CommandOption func(*Command)

// Command represents a command executed in the pseudo-terminal.
type Command struct {
	// Command line arguments, including the command itself as Args[0].
	Args []string

	// Environment of the command. Each entry is of the form "NAME=VALUE".
	// If nil, command will inherit environment from the current process.
	Env []string

	// Working directory of the command.
	Dir string

	// Command spawn flags.
	SpawnFlags SpawnFlags

	// Command PTY flags.
	PtyFlags PtyFlags

	// Timeout of the command.
	Timeout time.Duration

	// GLib Cancellable object.
	Cancellable *glib.Cancellable

	// OnSpawn is a callback that runs when command is spawned.
	// The second argument indicates whether there was an error.
	OnSpawn func(pid int, err error)
}

// WithEnv appends variable to the command environment.
//
// Can be used multiple times.
func WithEnv(name string, value any) CommandOption {
	return func(c *Command) {
		c.Env = append(c.Env, fmt.Sprintf("%s=%v", name, value))
	}
}

// WithWorkdir sets command working directory.
func WithWorkdir(workdir string) CommandOption {
	return func(c *Command) {
		c.Dir = workdir
	}
}

// WithSpawnFlags appends command spawn flags.
//
// Can be used multiple times.
func WithSpawnFlags(flags SpawnFlags) CommandOption {
	return func(c *Command) {
		c.SpawnFlags |= flags
	}
}

// WithPtyFlags appends command PTY flags.
//
// Can be used multiple times.
func WithPtyFlags(flags PtyFlags) CommandOption {
	return func(c *Command) {
		c.PtyFlags |= flags
	}
}

// WithTimeout sets command timeout.
func WithTimeout(timeout time.Duration) CommandOption {
	return func(c *Command) {
		c.Timeout = timeout
	}
}

// WithCancellable sets command cancellable.
func WithCancellable(cancellable *glib.Cancellable) CommandOption {
	return func(c *Command) {
		c.Cancellable = cancellable
	}
}

// WithOnSpawn sets callback that runs when command starts or fails to start.
func WithOnSpawn(callback func(pid int, err error)) CommandOption {
	return func(c *Command) {
		c.OnSpawn = callback
	}
}

// NewCommand creates a new [Command].
func NewCommand(args []string, options ...CommandOption) *Command {
	glib.CancellableNew()

	c := &Command{
		Args: args,
	}

	for _, option := range options {
		option(c)
	}

	commandSetDefaults(c)

	return c
}

// commandSetDefaults populates empty fields of [Command] with default values.
func commandSetDefaults(c *Command) {
	if c.Dir == "" {
		c.Dir, _ = os.Getwd()
	}

	if c.Env == nil {
		c.Env = os.Environ()
	}

	if c.Cancellable == nil {
		c.Cancellable, _ = glib.CancellableNew()
	}

	if c.Timeout == 0 {
		c.Timeout = DEFAULT_TIMEOUT
	}
}
