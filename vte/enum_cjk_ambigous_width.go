package vte

// CJKAmbiguousWidth is an enumeration type that represents width of
// ambiguous-width characters.
type CJKAmbiguousWidth int

const (
	// Narrow characters.
	CJK_AMBIGUOUS_WIDTH_NARROW CJKAmbiguousWidth = 1

	// Wide characters.
	CJK_AMBIGUOUS_WIDTH_WIDE CJKAmbiguousWidth = 2
)
