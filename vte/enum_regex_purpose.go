package vte

// RegexPurpose is an enumeration type that represents purpose of the regular
// expression.
type RegexPurpose int

const (
	// [Regex] for match. Can be used in [Terminal.MatchSetRegex].
	REGEX_PURPOSE_MATCH RegexPurpose = iota

	// [Regex] for search. Can be used in [Terminal.SearchSetRegex].
	REGEX_PURPOSE_SEARCH
)
