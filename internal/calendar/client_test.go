package calendar

import (
	"regexp"
	"testing"
)

func TestFormatEventID(t *testing.T) {
	// Google Calendar event ID requirement: 5-1024 characters from [a-v0-9]
	validRegex := regexp.MustCompile(`^[a-v0-9]+$`)

	testIDs := []int64{1, 42, 2000, 999999, 123456789}

	for _, id := range testIDs {
		formatted := FormatEventID(id)
		if len(formatted) < 5 || len(formatted) > 1024 {
			t.Errorf("FormatEventID(%d) length %d out of bounds [5, 1024]", id, len(formatted))
		}
		if !validRegex.MatchString(formatted) {
			t.Errorf("FormatEventID(%d) = %q does not match Google Calendar ID format requirement", id, formatted)
		}
	}
}
