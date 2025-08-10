package vte

// #include <vte/vte.h>
import "C"

// CursorShape is an enumeration type that can be used to indicate what should
// the terminal draw at the cursor position.
type CursorShape int

const (
	// Draw a block cursor. This is the default.
	CURSOR_SHAPE_BLOCK CursorShape = C.VTE_CURSOR_SHAPE_BLOCK

	// Draw a vertical bar on the left side of character.
	// This is similar to the default cursor for other GTK+ widgets.
	CURSOR_SHAPE_IBEAM CursorShape = C.VTE_CURSOR_SHAPE_IBEAM

	// Draw a horizontal bar below the character.
	CURSOR_SHAPE_UNDERLINE CursorShape = C.VTE_CURSOR_SHAPE_UNDERLINE
)
