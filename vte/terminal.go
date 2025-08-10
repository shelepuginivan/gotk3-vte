// Package vte provides Go bindings for Vte.
package vte

// #cgo pkg-config: gtk+-3.0 vte-2.91
// #include <gtk/gtk.h>
// #include <vte/vte.h>
// #include "vte.go.h"
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Terminal is a wrapper around VteTerminal.
type Terminal struct {
	gtk.Widget

	ptr *C.VteTerminal
}

// TerminalNew creates a new instance of [*Terminal].
func TerminalNew() (*Terminal, error) {
	ptr := C.vte_terminal_new()
	if ptr == nil {
		return nil, fmt.Errorf("vte_terminal_new returned nil pointer")
	}

	cGObject := glib.ToGObject(unsafe.Pointer(ptr))

	gObject := glib.Object{
		GObject: cGObject,
	}

	initiallyUnowned := glib.InitiallyUnowned{
		Object: &gObject,
	}

	widget := gtk.Widget{
		InitiallyUnowned: initiallyUnowned,
	}

	return &Terminal{
		Widget: widget,

		ptr: C.toTerminal(unsafe.Pointer(ptr)),
	}, nil
}

// GetPty returns [Pty] associated with the terminal.
func (t *Terminal) GetPty() *Pty {
	pty := C.vte_terminal_get_pty(t.ptr)
	return &Pty{pty}
}

// SetPty sets [Pty] to use in terminal. Use nil to unset the pty.
func (t *Terminal) SetPty(pty *Pty) {
	if pty == nil {
		pty = &Pty{}
	}

	C.vte_terminal_set_pty(t.ptr, pty.ptr)
}
