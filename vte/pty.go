package vte

// #include <glib.h>
// #include <gtk/gtk.h>
// #include <vte/vte.h>
// #include "exec.go.h"
// #include "glib.go.h"
// #include "vte.go.h"
import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// PtySize represents size of [Pty].
type PtySize struct {
	Rows    int
	Columns int
}

// Pty is a wrapper around VtePty.
type Pty struct {
	*glib.Object
}

func wrapPty(obj *glib.Object) *Pty {
	if obj == nil {
		return nil
	}
	return &Pty{obj}
}

func (pty *Pty) native() *C.VtePty {
	if pty == nil || pty.GObject == nil {
		return nil
	}
	return C.toPty(unsafe.Pointer(pty.GObject))
}

// PtyNewSync allocates a new pseudo-terminal.
func PtyNewSync(flags PtyFlags, cancellable *glib.Cancellable) (*Pty, error) {
	if cancellable == nil {
		return nil, fmt.Errorf("cancellable must not be nil")
	}

	var (
		gerr *C.GError
		c    = C.toCancellable(unsafe.Pointer(cancellable.GObject))
		f    = C.VtePtyFlags(flags)
	)

	pty := C.vte_pty_new_sync(f, c, &gerr)
	if pty == nil {
		if gerr == nil {
			return nil, errNilPointer("vte_pty_new_sync")
		}

		defer C.g_error_free(gerr)
		return nil, errFromGError("vte_pty_new_sync", gerr)
	}

	return wrapPty(glib.Take(unsafe.Pointer(pty))), nil
}

// PtyNewForeignSync creates a new [Pty] for the PTY master file. Newly created
// [Pty] will take ownership of PTY file descriptor and close it on finalize.
func PtyNewForeignSync(ptmx *os.File, cancellable *glib.Cancellable) (*Pty, error) {
	if cancellable == nil {
		return nil, fmt.Errorf("cancellable must not be nil")
	}

	var (
		gerr *C.GError
		c    = C.toCancellable(unsafe.Pointer(cancellable.GObject))
		fd   = C.int(ptmx.Fd())
	)

	pty := C.vte_pty_new_foreign_sync(fd, c, &gerr)
	if pty == nil {
		if gerr == nil {
			return nil, errNilPointer("vte_pty_new_foreign_sync")
		}

		defer C.g_error_free(gerr)
		return nil, errFromGError("vte_pty_new_foreign_sync", gerr)
	}

	return wrapPty(glib.Take(unsafe.Pointer(pty))), nil
}

// GetFd return the file descriptor of the PTY master in pty. The file
// descriptor belongs to pty and must not be closed or have its flags changed.
func (pty *Pty) GetFd() uintptr {
	return uintptr(C.vte_pty_get_fd(pty.native()))
}

// GetSize returns size of the pseudo terminal.
func (pty *Pty) GetSize() (*PtySize, error) {
	var (
		gerr    *C.GError
		rows    C.int
		columns C.int
	)

	success := C.vte_pty_get_size(pty.native(), &rows, &columns, &gerr)

	if !goBool(success) {
		if gerr == nil {
			return nil, errFailed("vte_pty_get_size")
		}

		defer C.g_error_free(gerr)
		return nil, errFromGError("vte_pty_get_size", gerr)
	}

	return &PtySize{
		Rows:    int(rows),
		Columns: int(columns),
	}, nil
}

// SetSize attempts to resize the pseudo terminal's window size. If successful,
// the OS kernel will send SIGWINCH to the child process group.
func (pty *Pty) SetSize(size *PtySize) error {
	var (
		gerr    *C.GError
		rows    = C.int(size.Rows)
		columns = C.int(size.Columns)
	)

	success := C.vte_pty_set_size(pty.native(), rows, columns, &gerr)

	if !goBool(success) {
		if gerr == nil {
			return errFailed("vte_pty_set_size")
		}

		defer C.g_error_free(gerr)
		return errFromGError("vte_pty_get_size", gerr)
	}

	return nil
}

// SetUTF8 tells the kernel whether the terminal is UTF-8 or not, in case it can
// make use of the info.
func (pty *Pty) SetUTF8(v bool) error {
	var (
		gerr *C.GError
		utf8 = gboolean(v)
	)

	success := C.vte_pty_set_utf8(pty.native(), utf8, &gerr)

	if !goBool(success) {
		if gerr == nil {
			return errFailed("vte_pty_set_utf8")
		}

		defer C.g_error_free(gerr)
		return errFromGError("vte_pty_set_utf8", gerr)
	}

	return nil
}

// Spawn starts the specified command under the pseudo-terminal pty.
//
// The command is spawned asynchronously. When cmd is executed or execution
// fails, its field Callback is called (if it is not nil) with the
// corresponding PID and occurred error.
//
// The TERM environment variable is automatically set to a default value,
// but can be overridden from cmd.
//
// [SPAWN_STDOUT_TO_DEV_NULL], [SPAWN_STDERR_TO_DEV_NULL], and
// [SPAWN_CHILD_INHERITS_STDIN] are not supported in flags, since stdin, stdout
// and stderr of the child process will always be connected to the PTY. Also
// [SPAWN_LEAVE_DESCRIPTORS_OPEN] is not supported; and
// [SPAWN_DO_NOT_REAP_CHILD] will always be added to spawn_flags.
func (pty *Pty) Spawn(cmd *Command) {
	var ccallID C.gpointer
	if cmd.OnSpawn != nil {
		callID := assignCallID(cmd)
		ccallID = C.uintToGpointer(C.uint(callID))
	}

	var (
		workdir               = C.CString(cmd.Dir)
		argv                  = cStringArr(cmd.Args)
		envv                  = cStringArr(cmd.Env)
		spawnFlags            = C.GSpawnFlags(cmd.SpawnFlags)
		childSetup            C.GSpawnChildSetupFunc
		childSetupData        C.gpointer
		cTimeout              = C.int(cmd.Timeout.Milliseconds())
		cCancellable          = C.toCancellable(unsafe.Pointer(cmd.Cancellable.GObject))
		childSetupDataDestroy C.GDestroyNotify
		callback              = C.GAsyncReadyCallback(C.ptySpawnAsyncCallback)
		userData              = ccallID
	)

	defer C.free(unsafe.Pointer(workdir))
	defer cStringArrFree(argv)
	defer cStringArrFree(envv)

	C.vte_pty_spawn_async(
		pty.native(),
		workdir,
		&argv[0],
		&envv[0],
		spawnFlags,
		childSetup,
		childSetupData,
		childSetupDataDestroy,
		cTimeout,
		cCancellable,
		callback,
		userData,
	)
}

func (pty *Pty) spawnFinish(res *C.GAsyncResult) (int, error) {
	var (
		gerr *C.GError
		pid  C.GPid
	)

	success := C.vte_pty_spawn_finish(pty.native(), res, &pid, &gerr)

	if !goBool(success) {
		if gerr == nil {
			return 0, errFailed("vte_pty_spawn_finish")
		}

		defer C.g_error_free(gerr)
		return 0, errFromGError("vte_pty_spawn_finish", gerr)
	}

	return int(pid), nil
}
