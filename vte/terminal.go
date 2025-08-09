// Package vte provides Go bindings for Vte.
package vte

// #cgo pkg-config: gtk+-3.0 vte-2.91
// #include <gtk/gtk.h>
// #include <vte/vte.h>
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

	return &Terminal{widget}, nil
}
