package vte

// #include <vte/vte.h>
import "C"

// TextBlinkMode is an enumeration type that can be used to indicate whether
// the terminal allows the text contents to be blinked.
type TextBlinkMode int

const (
	// Do not blink the text.
	TEXT_BLINK_NEVER TextBlinkMode = C.VTE_TEXT_BLINK_NEVER

	// Allow blinking text only if the terminal is focused.
	TEXT_BLINK_FOCUSED TextBlinkMode = C.VTE_TEXT_BLINK_FOCUSED

	// Allow blinking text only if the terminal is unfocused.
	TEXT_BLINK_UNFOCUSED TextBlinkMode = C.VTE_TEXT_BLINK_UNFOCUSED

	// Allow blinking text. This is the default.
	TEXT_BLINK_ALWAYS TextBlinkMode = C.VTE_TEXT_BLINK_ALWAYS
)
