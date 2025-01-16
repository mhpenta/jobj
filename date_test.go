package jobj

import (
	"testing"
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
