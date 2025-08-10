package vte

// #include <gtk/gtk.h>
// #include <vte/vte.h>
import "C"
import (
	"unsafe"
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
