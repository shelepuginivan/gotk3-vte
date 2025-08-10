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

// CopyClipboardFormat copies selected text to clipboard in the specified
// format.
func (t *Terminal) CopyClipboardFormat(format Format) {
	C.vte_terminal_copy_clipboard_format(t.ptr, C.VteFormat(format))
}

// CopyPrimary copies selected text in the primary selection.
func (t *Terminal) CopyPrimary() {
	C.vte_terminal_copy_primary(t.ptr)
}

// PasteClipboard pastes contents of clipboard to the terminal.
func (t *Terminal) PasteClipboard() {
	C.vte_terminal_paste_clipboard(t.ptr)
}

// PastePrimary pastes contents of the primary selection to the terminal.
func (t *Terminal) PastePrimary() {
	C.vte_terminal_paste_primary(t.ptr)
}

// PasteText pastes text to the terminal.
func (t *Terminal) PasteText(text string) {
	s := C.CString(text)
	C.vte_terminal_paste_text(t.ptr, s)
	C.free(unsafe.Pointer(s))
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
