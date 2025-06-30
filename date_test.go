package jobj

import (
	"testing"
	"time"
)

func TestJsonDateEquality(t *testing.T) {
	date1 := JsonDateTime{}
	err := date1.UnmarshalJSON([]byte(`"2023-01-01"`))
	if err != nil {
		t.Fatalf("Failed to unmarshal date1: %v", err)
	}

	date2 := JsonDateTime{}
	err = date2.UnmarshalJSON([]byte(`"2023-01-01"`))
	if err != nil {
		t.Fatalf("Failed to unmarshal date2: %v", err)
	}

	if !date1.Equal(date2.Time) {
		t.Errorf("Expected date1 and date2 to be equal, but they are not")
	}

	if date1.Format("2006-01-02") != date2.Format("2006-01-02") {
		t.Errorf("Expected formatted dates to be equal, but they are not")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		want      string
		wantErr   bool
	}{
		{
			name:      "Simple date format",
			jsonInput: `"2023-01-01"`,
			want:      "2023-01-01",
			wantErr:   false,
		},
		{
			name:      "RFC3339 format",
			jsonInput: `"2023-01-01T12:34:56Z"`,
			want:      "2023-01-01",
			wantErr:   false,
		},
		{
			name:      "RFC3339 with timezone",
			jsonInput: `"2023-01-01T12:34:56-07:00"`,
			want:      "2023-01-01",
			wantErr:   false,
		},
		{
			name:      "Invalid date format",
			jsonInput: `"not-a-date"`,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "Invalid JSON",
			jsonInput: `{"date": "2023-01-01"}`,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "Null input",
			jsonInput: `null`,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "European date format (unsupported)",
			jsonInput: `"01/01/2023"`,
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var date JsonDateTime
			err := date.UnmarshalJSON([]byte(tt.jsonInput))

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := date.Format("2006-01-02")
				if got != tt.want {
					t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestParsePublishedTime(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "RFC3339",
			timeStr: "2023-01-02T15:04:05Z",
			want:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "ISO8601 with timezone",
			timeStr: "2023-01-02T15:04:05-07:00",
			want:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
			wantErr: false,
		},
		{
			name:    "RFC3339Nano",
			timeStr: "2023-01-02T15:04:05.999999999Z",
			want:    time.Date(2023, 1, 2, 15, 4, 5, 999999999, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date only",
			timeStr: "2023-01-02",
			want:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "RFC1123",
			timeStr: "Mon, 02 Jan 2023 15:04:05 MST",
			want:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", 0)),
			wantErr: false,
		},
		{
			name:    "RFC1123Z",
			timeStr: "Mon, 02 Jan 2023 15:04:05 -0700",
			want:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
			wantErr: false,
		},
		{
			name:    "RFC822",
			timeStr: "02 Jan 23 15:04 MST",
			want:    time.Date(2023, 1, 2, 15, 4, 0, 0, time.FixedZone("MST", 0)),
			wantErr: false,
		},
		{
			name:    "RFC822Z",
			timeStr: "02 Jan 23 15:04 -0700",
			want:    time.Date(2023, 1, 2, 15, 4, 0, 0, time.FixedZone("", -7*60*60)),
			wantErr: false,
		},
		{
			name:    "Date and time without timezone",
			timeStr: "2023-01-02 15:04:05",
			want:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Invalid format",
			timeStr: "not-a-time",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Empty string",
			timeStr: "",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePublishedTime(tt.timeStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("parsePublishedTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !compareTimeComponents(got, tt.want) {
					t.Errorf("parsePublishedTime() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// Helper function to compare time components without worrying about location or nanoseconds
func compareTimeComponents(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute() &&
		t1.Second() == t2.Second()
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		wantErr   bool
	}{
		{
			name:      "Empty JSON string",
			jsonInput: `""`,
			wantErr:   true,
		},
		{
			name:      "Whitespace only",
			jsonInput: `"   "`,
			wantErr:   true,
		},
		{
			name:      "Number instead of string",
			jsonInput: `20230101`,
			wantErr:   true,
		},
		{
			name:      "Boolean instead of string",
			jsonInput: `true`,
			wantErr:   true,
		},
		{
			name:      "Array instead of string",
			jsonInput: `["2023-01-01"]`,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var date JsonDateTime
			err := date.UnmarshalJSON([]byte(tt.jsonInput))

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
