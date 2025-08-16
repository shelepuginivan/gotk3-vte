package vte

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestRegexNew(t *testing.T) {
	reg, err := RegexNew(".")

	assert.NoError(t, err)

	// Test default values
	assert.Equal(t, REGEX_PURPOSE_SEARCH, reg.purpose)
	assert.Equal(t, REGEX_COMPILE_FLAGS_MULTILINE, reg.flags)
	assert.Equal(t, RegexCompileExtraFlags(0), reg.extraFlags)

	t.Run("Invalid PCRE2 regular expression pattern", func(t *testing.T) {
		reg, err := RegexNew(`.{2,1}`)
		assert.Nil(t, reg)
		assert.Error(t, err)
	})

	t.Run("Finalizer", func(t *testing.T) {
		RegexNew(".*")
		runtime.GC()
	})
}

func TestRegex_RefCount(t *testing.T) {
	reg, err := RegexNew(".")
	assert.NoError(t, err)

	reg.Ref()
	reg.Unref()
	reg.Unref()

	// It is not possible to test, but uncommenting the line below should result
	// in SIGABRT, since `reg.ptr` is already freed at this point due to reference
	// count reached 0.
	//
	//reg.Unref()
}

func TestRegexWithPurpose(t *testing.T) {
	reg, err := RegexNew(".", RegexWithPurpose(REGEX_PURPOSE_MATCH))
	assert.NoError(t, err)
	assert.Equal(t, REGEX_PURPOSE_MATCH, reg.purpose)
}

func TestRegexWithCompileFlags(t *testing.T) {
	reg, err := RegexNew(
		".",
		RegexWithCompileFlags(REGEX_COMPILE_FLAGS_ANCHORED),
		RegexWithCompileFlags(REGEX_COMPILE_FLAGS_ALT_BSUX|REGEX_COMPILE_FLAGS_ALT_CIRCUMFLEX),
	)

	expected := REGEX_COMPILE_FLAGS_MULTILINE |
		REGEX_COMPILE_FLAGS_ANCHORED |
		REGEX_COMPILE_FLAGS_ALT_BSUX |
		REGEX_COMPILE_FLAGS_ALT_CIRCUMFLEX

	assert.NoError(t, err)
	assert.Equal(t, expected, reg.flags)
}

func TestRegexWithCompileExtraFlags(t *testing.T) {
	reg, err := RegexNew(
		".",
		RegexWithCompileExtraFlags(REGEX_COMPILE_EXTRA_FLAGS_ALT_BSUX),
		RegexWithCompileExtraFlags(REGEX_COMPILE_EXTRA_FLAGS_ASCII_BSW|REGEX_COMPILE_EXTRA_FLAGS_PYTHON_OCTAL),
	)

	expected := REGEX_COMPILE_EXTRA_FLAGS_ALT_BSUX |
		REGEX_COMPILE_EXTRA_FLAGS_ASCII_BSW |
		REGEX_COMPILE_EXTRA_FLAGS_PYTHON_OCTAL

	assert.NoError(t, err)
	assert.Equal(t, expected, reg.extraFlags)
}

func TestRegex_Native(t *testing.T) {
	reg, err := RegexNew(".")
	assert.NoError(t, err)
	assert.NotEqual(t, uintptr(unsafe.Pointer(nil)), reg.Native())
}
