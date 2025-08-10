package vte

// #include <vte/vte.h>
import "C"

// PtyFlags is a bitfield type that represents [Pty] flags.
type PtyFlags int

const (
	PTY_DEFAULT     PtyFlags = C.VTE_PTY_DEFAULT
	PTY_NO_LASTLOG  PtyFlags = C.VTE_PTY_NO_LASTLOG
	PTY_NO_UTMP     PtyFlags = C.VTE_PTY_NO_UTMP
	PTY_NO_WTMP     PtyFlags = C.VTE_PTY_NO_WTMP
	PTY_NO_HELPER   PtyFlags = C.VTE_PTY_NO_HELPER
	PTY_NO_FALLBACK PtyFlags = C.VTE_PTY_NO_FALLBACK
	PTY_NO_SESSION  PtyFlags = C.VTE_PTY_NO_SESSION
	PTY_NO_CTTY     PtyFlags = C.VTE_PTY_NO_CTTY
)
