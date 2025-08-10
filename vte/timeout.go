package vte

// #include <glib.h>
import "C"
import "time"

const (
	DEFAULT_TIMEOUT    = -1 * time.Millisecond
	INDEFINITE_TIMEOUT = time.Duration(C.G_MAXINT) * time.Millisecond
)
