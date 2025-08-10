package vte

// #include <vte/vte.h>
import "C"

// Format is an enumeration type that can be used to specify the format the
// selection should be copied to the clipboard in.
type Format int

const (
	// Export as plain text.
	FORMAT_TEXT Format = C.VTE_FORMAT_TEXT

	// Export as HTML formatted text.
	FORMAT_HTML Format = C.VTE_FORMAT_HTML
)
