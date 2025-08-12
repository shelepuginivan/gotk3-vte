// Package vte provides Go bindings for VTE.
package vte

// #include <vte/vte.h>
import "C"
import "unsafe"

// GetMajorVersion returns the major version of VTE.
func GetMajorVersion() uint {
	return uint(C.vte_get_major_version())
}

// GetMinorVersion returns the minor version of VTE.
func GetMinorVersion() uint {
	return uint(C.vte_get_minor_version())
}

// GetMicroVersion returns the micro version of VTE.
func GetMicroVersion() uint {
	return uint(C.vte_get_micro_version())
}

// GetUserShell returns user's shell. If it is unable to determine the user
// shell, "/bin/sh" is returned.
func GetUserShell() string {
	cstr := C.vte_get_user_shell()
	if cstr == nil {
		return "/bin/sh"
	}
	defer C.free(unsafe.Pointer(cstr))
	return goString(cstr)
}
