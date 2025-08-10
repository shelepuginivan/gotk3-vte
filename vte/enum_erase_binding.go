package vte

// #include <vte/vte.h>
import "C"

// EraseBinding is an enumeration type that can be used to indicate which
// string the terminal should send to an application when the user presses the
// Delete or Backspace keys.
type EraseBinding int

const (
	// For backspace, attempt to determine the right value from the terminal's IO
	// settings. For delete, use the control sequence.
	ERASE_AUTO EraseBinding = C.VTE_ERASE_AUTO

	// Send an ASCII backspace character (0x08).
	ERASE_ASCII_BACKSPACE EraseBinding = C.VTE_ERASE_ASCII_BACKSPACE

	// Send an ASCII delete character (0x7F).
	ERASE_ASCII_DELETE EraseBinding = C.VTE_ERASE_ASCII_DELETE

	// Send the "@7" control sequence.
	ERASE_DELETE_SEQUENCE EraseBinding = C.VTE_ERASE_DELETE_SEQUENCE

	// Send terminal's "erase" setting.
	ERASE_TTY EraseBinding = C.VTE_ERASE_TTY
)
