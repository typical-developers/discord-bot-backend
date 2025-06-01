package regexutil

import (
	"regexp"
)

var (
	// Checks if a string is a valid Discord snowflake id.
	Snowflake = regexp.MustCompile(`^(?<id>\d{17,20})$`)
)
