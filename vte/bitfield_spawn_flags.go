package vte

// #include <glib.h>
import "C"

// SpawnFlags is a bitfield type that represents flags passed to spawn
// functions.
type SpawnFlags int

const (
	// No flags, default behaviour.
	SPAWN_DEFAULT SpawnFlags = C.G_SPAWN_DEFAULT

	// The parent's open file descriptors will be inherited by the child;
	// otherwise all descriptors except stdin, stdout and stderr will be closed
	// before calling exec() in the child.
	SPAWN_LEAVE_DESCRIPTORS_OPEN SpawnFlags = C.G_SPAWN_LEAVE_DESCRIPTORS_OPEN

	// The child will not be automatically reaped; you must handle SIGCHLD
	// yourself, or the child will become a zombie.
	SPAWN_DO_NOT_REAP_CHILD SpawnFlags = C.G_SPAWN_DO_NOT_REAP_CHILD

	// Search for the executable (argv[0]) in the user's PATH.
	SPAWN_SEARCH_PATH SpawnFlags = C.G_SPAWN_SEARCH_PATH

	// The child's standard output will be discarded, instead of going to the
	// same location as the parent's standard output.
	SPAWN_STDOUT_TO_DEV_NULL SpawnFlags = C.G_SPAWN_STDOUT_TO_DEV_NULL

	// The child's standard error will be discarded, instead of going to the
	// same location as the parent's standard error.
	SPAWN_STDERR_TO_DEV_NULL SpawnFlags = C.G_SPAWN_STDERR_TO_DEV_NULL

	// The child will inherit the parent's standard input
	// (by default, the child's standard input is attached to /dev/null).
	SPAWN_CHILD_INHERITS_STDIN SpawnFlags = C.G_SPAWN_CHILD_INHERITS_STDIN

	// The first element of argv is the file to execute, while the remaining
	// elements are the actual argument vector to pass to the file.
	SPAWN_FILE_AND_ARGV_ZERO SpawnFlags = C.G_SPAWN_FILE_AND_ARGV_ZERO

	// If executable (argv[0]) is not an absolute path, it will be looked for in
	// the PATH from the passed child environment.
	SPAWN_SEARCH_PATH_FROM_ENVP SpawnFlags = C.G_SPAWN_SEARCH_PATH_FROM_ENVP

	// Create all pipes with the O_CLOEXEC flag set.
	SPAWN_CLOEXEC_PIPES SpawnFlags = C.G_SPAWN_CLOEXEC_PIPES

	// The child will inherit the parent's standard output.
	SPAWN_CHILD_INHERITS_STDOUT SpawnFlags = C.G_SPAWN_CHILD_INHERITS_STDOUT

	// The child will inherit the parent's standard error.
	SPAWN_CHILD_INHERITS_STDERR SpawnFlags = C.G_SPAWN_CHILD_INHERITS_STDERR

	// The child's standard input is attached to /dev/null.
	SPAWN_STDIN_FROM_DEV_NULL SpawnFlags = C.G_SPAWN_STDIN_FROM_DEV_NULL
)
