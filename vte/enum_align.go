package vte

// #include <vte/vte.h>
import "C"

// Align is an enumeration type that can be used to specify how the terminal
// uses extra allocated space.
type Align int

const (
	// Align to left/top.
	ALIGN_START Align = C.VTE_ALIGN_START

	// Align to center.
	ALIGN_CENTER Align = C.VTE_ALIGN_CENTER

	// Align to right/bottom.
	ALIGN_END Align = C.VTE_ALIGN_END
)
