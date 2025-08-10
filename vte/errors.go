package vte

// #include <glib.h>
import "C"
import "fmt"

type ErrCgoCall struct {
	Function string
	Detail   string
}

func (err ErrCgoCall) Error() string {
	return fmt.Sprintf("vte: %s() %s", err.Function, err.Detail)
}

func errFromGError(function string, gerr *C.GError) error {
	return ErrCgoCall{
		Function: function,
		Detail:   goString(gerr.message),
	}
}

func errNilPointer(function string) error {
	return ErrCgoCall{
		Function: function,
		Detail:   "returned a nil pointer",
	}
}

func errFailed(function string) error {
	return ErrCgoCall{
		Function: function,
		Detail:   "failed",
	}
}
