package jobj

import (
	"fmt"
	"time"
)

const RFC3339Basic = "2006-01-02T15:04:05Z"

// parsePublishedTime attempts to parse a string representing a published time using various time formats. Typically
// used to parse feeds.
func parsePublishedTime(published string) (time.Time, error) {
	formats := []string{
		RFC3339Basic,
		"2006-01-02T15:04:05-07:00", // ISO 8601 format with timezone
		"2006-01-02 15:04:05",       // Default format without timezone
		time.DateOnly,
		time.RFC3339,
		time.RubyDate,
		time.RFC822,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC850,
		time.RFC3339,
		time.RFC3339Nano,
		time.Layout,
		time.ANSIC,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	var err error
	var publishedTime time.Time

	for _, format := range formats {
		publishedTime, err = time.Parse(format, published)
		if err == nil {
			return publishedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse published time: %v", err)
}
