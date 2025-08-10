package vte

// #include <gtk/gtk.h>
// #include <pango/pango.h>
// #include <vte/vte.h>
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/pango"
)

func goString(cstr *C.gchar) string {
	return C.GoString((*C.char)(cstr))
}

func goBool(b C.gboolean) bool {
	return b != C.FALSE
}

func gboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func cStringArr(s []string) []*C.char {
	cArr := make([]*C.char, len(s)+1)
	for i, s := range s {
		cArr[i] = C.CString(s)
	}

	cArr[len(s)] = (*C.char)(nil)

	return cArr
}

func cStringArrFree(cArr []*C.char) {
	for i, cstr := range cArr {
		if i != len(cArr)-1 {
			C.free(unsafe.Pointer(cstr))
		}
	}
}

func wrapPangoFontDescription(desc *C.PangoFontDescription) *pango.FontDescription {
	fd := &pango.FontDescription{}

	// We need to set `fd.pangoFontDescription`, but it is private.
	// Hence, we create a pointer to this field and write desc to it.
	ptrFd := unsafe.Pointer(fd)
	ptrFdNative := (**C.PangoFontDescription)(ptrFd)

	*ptrFdNative = desc
	return fd
}

func unwrapPangoFontDescription(desc *pango.FontDescription) *C.PangoFontDescription {
	ptrFd := unsafe.Pointer(desc)
	ptrFdNative := (**C.PangoFontDescription)(ptrFd)

	return *ptrFdNative
}
