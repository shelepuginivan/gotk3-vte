package vte

// #cgo pkg-config: gtk+-3.0 vte-2.91
// #include <gtk/gtk.h>
// #include <vte/vte.h>
import "C"

func goString(cstr *C.gchar) string {
	return C.GoString((*C.char)(cstr))
}
