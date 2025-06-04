package regexutil

import (
	"fmt"
	"regexp"
)

var (
	// Checks if a string is a valid Discord snowflake id.
	Snowflake = regexp.MustCompile(`^(?<id>\d{17,20})$`)
)

type SnowflakeIDs struct {
	Key string
	ID  string
}

func CheckSnowflakes(ids []SnowflakeIDs) error {
	for _, id := range ids {
		if !Snowflake.MatchString(id.ID) {
			return fmt.Errorf("%s is not a snowflake", id.Key)
		}
	}

	return nil
}
