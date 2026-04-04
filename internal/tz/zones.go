// Package tz provides timezone listing, lookup, and system timezone helpers.
package tz

import (
	"errors"
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"
)

// ErrPermission is returned by SetSystemTZ when the process lacks the
// privileges needed to change the system timezone. The caller should
// retry via SetSystemTZWithSudo after obtaining a password.
var ErrPermission = errors.New("permission denied: root privileges required")

// Zone represents an IANA timezone entry with a human-friendly label.
type Zone struct {
	IANA  string
	Label string
}

// ListZones returns all IANA timezones from the embedded tzdata database,
// sorted alphabetically.
func ListZones() []Zone {
	zones := make([]Zone, 0, len(ianaZones))
	for _, iana := range ianaZones {
		zones = append(zones, Zone{IANA: iana, Label: labelFromIANA(iana)})
	}
	return zones
}

// labelFromIANA derives a human-friendly label from an IANA timezone name,
// e.g. "America/New_York" → "New York", "Europe/London" → "London".
func labelFromIANA(iana string) string {
	parts := strings.Split(iana, "/")
	return strings.ReplaceAll(parts[len(parts)-1], "_", " ")
}

// LabelForIANA returns the display label for an IANA identifier.
func LabelForIANA(iana string) string {
	return labelFromIANA(iana)
}

// CurrentTime returns the current time in the given IANA zone.
func CurrentTime(iana string) (time.Time, error) {
	loc, err := time.LoadLocation(iana)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown timezone %q: %w", iana, err)
	}
	return time.Now().In(loc), nil
}
