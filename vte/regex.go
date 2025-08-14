package vte

// #include <vte/vte.h>
// #include <glib.h>
//
// #include "glib.go.h"
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// RegexOption allows to configure [Regex].
type RegexOption func(*Regex)

// Regex represents PCRE2 regular expression used for matching and searching
// text in [Terminal].
type Regex struct {
	ptr *C.VteRegex

	pattern    string
	purpose    RegexPurpose
	flags      RegexCompileFlags
	extraFlags RegexCompileExtraFlags
}

// RegexWithPurpose sets regex purpose.
func RegexWithPurpose(purpose RegexPurpose) RegexOption {
	return func(r *Regex) {
		r.purpose = purpose
	}
}

// RegexWithCompileFlags appends flags used for compilation to regex.
//
// Can be used multiple times.
func RegexWithCompileFlags(flags RegexCompileFlags) RegexOption {
	return func(r *Regex) {
		r.flags |= flags
	}
}

// RegexWithCompileExtraFlags appends extra flags used for compilation to
// regex.
//
// Can be used multiple times.
func RegexWithCompileExtraFlags(extraFlags RegexCompileExtraFlags) RegexOption {
	return func(r *Regex) {
		r.extraFlags |= extraFlags
	}
}

// RegexNew returns a new [Regex].
func RegexNew(pattern string, options ...RegexOption) (*Regex, error) {
	r := &Regex{
		pattern: pattern,
		purpose: REGEX_PURPOSE_SEARCH,

		// NOTE: both vte_terminal_match_add_regex and vte_terminal_search_add_regex
		// require this flag
		flags: REGEX_COMPILE_FLAGS_MULTILINE,
	}

	for _, option := range options {
		option(r)
	}

	var (
		gerr        *C.GError
		cPattern    = C.CString(pattern)
		cLength     = C.intToGssize(C.int(len(pattern)))
		cFlags      = C.uint(r.flags)
		cExtraFlags = C.uint(r.extraFlags)
		cOffset     = C.uintToGsize(0)

		function string
	)

	switch r.purpose {
	case REGEX_PURPOSE_MATCH:
		r.ptr = C.vte_regex_new_for_match_full(cPattern, cLength, cFlags, cExtraFlags, &cOffset, &gerr)
		function = "vte_regex_new_for_match_full"
	case REGEX_PURPOSE_SEARCH:
		r.ptr = C.vte_regex_new_for_search_full(cPattern, cLength, cFlags, cExtraFlags, &cOffset, &gerr)
		function = "vte_regex_new_for_search_full"
	}

	if gerr != nil {
		defer C.g_error_free(gerr)
		return nil, errFromGError(function, gerr)
	}

	if r.ptr == nil {
		return nil, errNilPointer(function)
	}

	runtime.SetFinalizer(r, func(r *Regex) { glib.FinalizerStrategy(r.Unref) })

	return r, nil
}

func wrapRegex(ptr *C.VteRegex) *Regex {
	r := &Regex{ptr: ptr}
	runtime.SetFinalizer(r, func(r *Regex) { glib.FinalizerStrategy(r.Unref) })
	return r
}

// String returns the source pattern used to compile the regular expression.
func (r *Regex) String() string {
	return r.pattern
}

// Ref increases the reference count of regex by 1.
//
// Reference counting is handled at package level. Most applications should
// not need to call this.
func (r *Regex) Ref() {
	C.vte_regex_ref(r.ptr)
}

// Unref decreases the reference count of regex by 1, and frees it whenever
// the count reaches 0.
//
// Reference counting is handled at package level. Most applications should
// not need to call this.
func (r *Regex) Unref() {
	C.vte_regex_unref(r.ptr)
}

// Native returns a pointer to the underlying VteRegex.
func (r *Regex) Native() uintptr {
	return uintptr(unsafe.Pointer(r.ptr))
}
