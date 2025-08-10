package vte

// #include <vte/vte.h>
import "C"

// CursorBlinkMode is an enumeration type that can be used to indicate the
// cursor blink mode for the terminal.
type CursorBlinkMode int

const (
	// Follow GTK+ settings for cursor blinking.
	CURSOR_BLINK_SYSTEM CursorBlinkMode = C.VTE_CURSOR_BLINK_SYSTEM

	// Cursor blinks.
	CURSOR_BLINK_ON CursorBlinkMode = C.VTE_CURSOR_BLINK_ON

	// Cursor does not blink.
	CURSOR_BLINK_OFF CursorBlinkMode = C.VTE_CURSOR_BLINK_ON
)
