package vte

// #include <vte/vte.h>
// #include <glib.h>
//
// #include "glib.go.h"
import "C"

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

	return r, nil
}

// String returns the source pattern used to compile the regular expression.
func (r *Regex) String() string {
	return r.pattern
}
