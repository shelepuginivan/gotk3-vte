package vte

// RegexCompileExtraFlags is a bitfield type that represents PCRE2 extra
// compilation flags.
//
// See man:pcre2api(3) for more information about every flag.
type RegexCompileExtraFlags uint

const (
	REGEX_COMPILE_EXTRA_FLAGS_ALLOW_SURROGATE_ESCAPES RegexCompileExtraFlags = 0x00000001
	REGEX_COMPILE_EXTRA_FLAGS_BAD_ESCAPE_IS_LITERAL   RegexCompileExtraFlags = 0x00000002
	REGEX_COMPILE_EXTRA_FLAGS_MATCH_WORD              RegexCompileExtraFlags = 0x00000004
	REGEX_COMPILE_EXTRA_FLAGS_MATCH_LINE              RegexCompileExtraFlags = 0x00000008
	REGEX_COMPILE_EXTRA_FLAGS_ESCAPED_CR_IS_LF        RegexCompileExtraFlags = 0x00000010
	REGEX_COMPILE_EXTRA_FLAGS_ALT_BSUX                RegexCompileExtraFlags = 0x00000020
	REGEX_COMPILE_EXTRA_FLAGS_ALLOW_LOOKAROUND_BSK    RegexCompileExtraFlags = 0x00000040
	REGEX_COMPILE_EXTRA_FLAGS_CASELESS_RESTRICT       RegexCompileExtraFlags = 0x00000080
	REGEX_COMPILE_EXTRA_FLAGS_ASCII_BSD               RegexCompileExtraFlags = 0x00000100
	REGEX_COMPILE_EXTRA_FLAGS_ASCII_BSS               RegexCompileExtraFlags = 0x00000200
	REGEX_COMPILE_EXTRA_FLAGS_ASCII_BSW               RegexCompileExtraFlags = 0x00000400
	REGEX_COMPILE_EXTRA_FLAGS_ASCII_POSIX             RegexCompileExtraFlags = 0x00000800
	REGEX_COMPILE_EXTRA_FLAGS_ASCII_DIGIT             RegexCompileExtraFlags = 0x00001000
	REGEX_COMPILE_EXTRA_FLAGS_PYTHON_OCTAL            RegexCompileExtraFlags = 0x00002000
	REGEX_COMPILE_EXTRA_FLAGS_NO_BS0                  RegexCompileExtraFlags = 0x00004000
	REGEX_COMPILE_EXTRA_FLAGS_NEVER_CALLOUT           RegexCompileExtraFlags = 0x00008000
	REGEX_COMPILE_EXTRA_FLAGS_TURKISH_CASING          RegexCompileExtraFlags = 0x00010000
)
