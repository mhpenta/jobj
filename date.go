package jobj

import (
	"encoding/json"
	"fmt"
	"time"
)

// JsonDateTime is a custom date type that can be unmarshalled from JSON for YYYY-MM-DD format.
// It embeds the time.Time type and provides custom unmarshaling behavior.
type JsonDateTime struct {
	time.Time
}

// UnmarshalJSON is the custom unmarshaling method for JsonDateTime.
// It expects the JSON date string to be in the format "YYYY-MM-DD".
// It parses the date string and assigns the parsed time to the embedded Time field.
//
// Parameters:
//   - data: The JSON data to be unmarshalled, represented as a byte slice.
//
// Returns:
//   - error: An error if the JSON data is invalid or the date string is not in the expected format.
//     Returns nil if the unmarshaling is successful.
func (d *JsonDateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		t, err = parsePublishedTime(s)
		if err != nil {
			return fmt.Errorf("invalid date format: %v", err)
		}
	}
	d.Time = t
	return nil
}

const (
	RFC3339Basic       = "2006-01-02T15:04:05Z"
	RFC3339WithOffset  = "2006-01-02T15:04:05-07:00"
	DateTimeWithSpace  = "2006-01-02 15:04:05"
)

var timeFormats = []string{
	RFC3339WithOffset,
	time.RFC3339,
	RFC3339Basic,
	time.RFC3339Nano,
	DateTimeWithSpace,
	time.DateOnly,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RubyDate,
	time.RFC850,
	time.Layout,
	time.ANSIC,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

// parsePublishedTime attempts to parse a string representing a published time using various time formats. Typically
// used to parse feeds.
func parsePublishedTime(published string) (time.Time, error) {
	var err error
	var publishedTime time.Time
	for _, format := range timeFormats {
		publishedTime, err = time.Parse(format, published)
		if err == nil {
			return publishedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse published time: %v", err)
}
