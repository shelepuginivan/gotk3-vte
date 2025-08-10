package vte

// #cgo pkg-config: gtk+-3.0 vte-2.91
// #include <glib.h>
import "C"
import "time"

const (
	DEFAULT_TIMEOUT    = -1 * time.Millisecond
	INDEFINITE_TIMEOUT = time.Duration(C.G_MAXINT) * time.Millisecond
)
