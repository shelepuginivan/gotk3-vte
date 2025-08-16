package vte

// #include <glib.h>
import "C"
import "fmt"

type errCgoCall struct {
	Function string
	Detail   string
}

func (err errCgoCall) Error() string {
	return fmt.Sprintf("vte: %s() %s", err.Function, err.Detail)
}

func errFromGError(function string, gerr *C.GError) error {
	return errCgoCall{
		Function: function,
		Detail:   goString(gerr.message),
	}
}

func errNilPointer(function string) error {
	return errCgoCall{
		Function: function,
		Detail:   "returned nil pointer",
	}
}

func errFailed(function string) error {
	return errCgoCall{
		Function: function,
		Detail:   "failed",
	}
}
