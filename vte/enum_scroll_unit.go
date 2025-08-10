package vte

// ScrollUnit is an enumeration type that scroll units used by the [Terminal].
type ScrollUnit int

const (
	// Measure scroll amount in lines.
	SCROLL_UNIT_LINES ScrollUnit = 0

	// Measure scroll amount in pixels.
	SCROLL_UNIT_PIXELS ScrollUnit = 1
)
