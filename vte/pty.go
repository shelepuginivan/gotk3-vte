package vte

// #cgo pkg-config: gtk+-3.0 vte-2.91
// #include <glib.h>
// #include <gtk/gtk.h>
// #include <vte/vte.h>
// #include "glib.go.h"
import "C"
import (
	"errors"
	"os"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// PtyFlags represents VtePtyFlags.
type PtyFlags int

const (
	PTY_DEFAULT     PtyFlags = C.VTE_PTY_DEFAULT
	PTY_NO_LASTLOG  PtyFlags = C.VTE_PTY_NO_LASTLOG
	PTY_NO_UTMP     PtyFlags = C.VTE_PTY_NO_UTMP
	PTY_NO_WTMP     PtyFlags = C.VTE_PTY_NO_WTMP
	PTY_NO_HELPER   PtyFlags = C.VTE_PTY_NO_HELPER
	PTY_NO_FALLBACK PtyFlags = C.VTE_PTY_NO_FALLBACK
	PTY_NO_SESSION  PtyFlags = C.VTE_PTY_NO_SESSION
	PTY_NO_CTTY     PtyFlags = C.VTE_PTY_NO_CTTY
)

// Pty is a wrapper around VtePty.
type Pty struct {
	ptr *C.VtePty
}

// PtyNewSync allocates a new pseudo-terminal.
func PtyNewSync(flags PtyFlags, cancellable *glib.Cancellable) (*Pty, error) {
	var (
		gerr *C.GError
		c    = C.toCancellable(unsafe.Pointer(cancellable.GObject))
		f    = C.VtePtyFlags(flags)
	)

	pty := C.vte_pty_new_sync(f, c, &gerr)
	if pty == nil {
		defer C.g_error_free(gerr)
		return nil, errors.New(goString(gerr.message))
	}

	return &Pty{pty}, nil
}

// PtyNewForeignSync creates a new [Pty] for the PTY master file. Newly created
// [PTY] will take ownership of file's descriptor and close it on finalize.
func PtyNewForeignSync(file *os.File, cancellable *glib.Cancellable) (*Pty, error) {
	var (
		gerr *C.GError
		c    = C.toCancellable(unsafe.Pointer(cancellable.GObject))
		fd   = C.int(file.Fd())
	)

	pty := C.vte_pty_new_foreign_sync(fd, c, &gerr)
	if pty == nil {
		defer C.g_error_free(gerr)
		return nil, errors.New(goString(gerr.message))
	}

	return &Pty{pty}, nil
}
