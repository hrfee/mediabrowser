package mediabrowser

import (
	"time"
)

// TimeFormat is the go time.Time layout for formatting/parsing times sent to/from Jellyfin & Emby.
const TimeFormat = "2006-01-02T15:04:05"

type magicParse struct {
	Parsed time.Time `json:"parseme"`
}

// Time embeds time.Time with a custom JSON Unmarshal method to work with Jellyfin & Emby's time formatting.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	// Trim quotes from beginning and end, and any number of Zs (indicates UTC).
	for b[0] == '"' {
		b = b[1:]
	}
	for b[len(b)-1] == '"' || b[len(b)-1] == 'Z' {
		b = b[:len(b)-1]
	}
	// Trim nanoseconds and anything after, we don't care
	i := len(b) - 1
	for b[i] != '.' && i > 0 {
		i--
	}
	if i != 0 {
		b = b[:i]
	}
	t.Time, err = time.Parse(TimeFormat, string(b))
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Format(TimeFormat) + "\""), nil
}

func toMediabrowserTime(t ConvertableToMediabrowserTime) string {
	return t.Format(TimeFormat)
}

type ConvertableToMediabrowserTime interface {
	Format(layout string) string
}
