package vte

// RegexMatchFlags is a bitfield type that represents PCRE2 match flags.
//
// See man:pcre2api(3) for more information about every flag.
type RegexMatchFlags uint

const (
	REGEX_MATCH_FLAGS_ANCHORED                  RegexMatchFlags = 0x80000000
	REGEX_MATCH_FLAGS_NO_UTF_CHECK              RegexMatchFlags = 0x40000000
	REGEX_MATCH_FLAGS_ENDANCHORED               RegexMatchFlags = 0x20000000
	REGEX_MATCH_FLAGS_NOTBOL                    RegexMatchFlags = 0x00000001
	REGEX_MATCH_FLAGS_NOTEOL                    RegexMatchFlags = 0x00000002
	REGEX_MATCH_FLAGS_NOTEMPTY                  RegexMatchFlags = 0x00000004
	REGEX_MATCH_FLAGS_NOTEMPTY_ATSTART          RegexMatchFlags = 0x00000008
	REGEX_MATCH_FLAGS_PARTIAL_SOFT              RegexMatchFlags = 0x00000010
	REGEX_MATCH_FLAGS_PARTIAL_HARD              RegexMatchFlags = 0x00000020
	REGEX_MATCH_FLAGS_COPY_MATCHED_SUBJECT      RegexMatchFlags = 0x00004000
	REGEX_MATCH_FLAGS_DISABLE_RECURSELOOP_CHECK RegexMatchFlags = 0x00040000
)
