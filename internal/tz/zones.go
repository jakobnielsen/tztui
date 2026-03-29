// Package tz provides timezone listing, lookup, and system timezone helpers.
package tz

import (
	"errors"
	"fmt"
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

// Zones is the curated list of zones shown in the TUI.
var Zones = []Zone{
	{"UTC", "UTC"},
	{"Etc/GMT+12", "UTC-12"},
	{"Etc/GMT+11", "UTC-11"},
	{"Etc/GMT+10", "UTC-10"},
	{"Etc/GMT+9", "UTC-9"},
	{"Etc/GMT+8", "UTC-8"},
	{"Etc/GMT+7", "UTC-7"},
	{"Etc/GMT+6", "UTC-6"},
	{"Etc/GMT+5", "UTC-5"},
	{"Etc/GMT+4", "UTC-4"},
	{"Etc/GMT+3", "UTC-3"},
	{"Etc/GMT+2", "UTC-2"},
	{"Etc/GMT+1", "UTC-1"},
	{"Etc/GMT-1", "UTC+1"},
	{"Etc/GMT-2", "UTC+2"},
	{"Etc/GMT-3", "UTC+3"},
	{"Etc/GMT-4", "UTC+4"},
	{"Etc/GMT-5", "UTC+5"},
	{"Etc/GMT-6", "UTC+6"},
	{"Etc/GMT-7", "UTC+7"},
	{"Etc/GMT-8", "UTC+8"},
	{"Etc/GMT-9", "UTC+9"},
	{"Etc/GMT-10", "UTC+10"},
	{"Etc/GMT-11", "UTC+11"},
	{"Etc/GMT-12", "UTC+12"},
	{"Africa/Abidjan", "Abidjan"},
	{"Africa/Cairo", "Cairo"},
	{"Africa/Casablanca", "Casablanca"},
	{"Africa/Johannesburg", "Johannesburg"},
	{"Africa/Lagos", "Lagos"},
	{"Africa/Nairobi", "Nairobi"},
	{"Africa/Tunis", "Tunis"},
	{"America/Anchorage", "Anchorage"},
	{"America/Argentina/Buenos_Aires", "Buenos Aires"},
	{"America/Bogota", "Bogota"},
	{"America/Caracas", "Caracas"},
	{"America/Chicago", "Chicago"},
	{"America/Denver", "Denver"},
	{"America/Halifax", "Halifax"},
	{"America/Lima", "Lima"},
	{"America/Los_Angeles", "Los Angeles"},
	{"America/Mexico_City", "Mexico City"},
	{"America/New_York", "New York"},
	{"America/Phoenix", "Phoenix"},
	{"America/Santiago", "Santiago"},
	{"America/Sao_Paulo", "Sao Paulo"},
	{"America/St_Johns", "St. Johns"},
	{"America/Toronto", "Toronto"},
	{"America/Vancouver", "Vancouver"},
	{"Asia/Bangkok", "Bangkok"},
	{"Asia/Colombo", "Colombo"},
	{"Asia/Dubai", "Dubai"},
	{"Asia/Hong_Kong", "Hong Kong"},
	{"Asia/Jakarta", "Jakarta"},
	{"Asia/Kabul", "Kabul"},
	{"Asia/Karachi", "Karachi"},
	{"Asia/Kathmandu", "Kathmandu"},
	{"Asia/Kolkata", "Kolkata"},
	{"Asia/Kuala_Lumpur", "Kuala Lumpur"},
	{"Asia/Kuwait", "Kuwait"},
	{"Asia/Manila", "Manila"},
	{"Asia/Riyadh", "Riyadh"},
	{"Asia/Seoul", "Seoul"},
	{"Asia/Shanghai", "Shanghai"},
	{"Asia/Singapore", "Singapore"},
	{"Asia/Tashkent", "Tashkent"},
	{"Asia/Tehran", "Tehran"},
	{"Asia/Tokyo", "Tokyo"},
	{"Asia/Vladivostok", "Vladivostok"},
	{"Asia/Yangon", "Yangon"},
	{"Asia/Yekaterinburg", "Yekaterinburg"},
	{"Atlantic/Azores", "Azores"},
	{"Atlantic/Reykjavik", "Reykjavik"},
	{"Pacific/Auckland", "Auckland"},
	{"Pacific/Fiji", "Fiji"},
	{"Pacific/Guam", "Guam"},
	{"Pacific/Honolulu", "Honolulu"},
	{"Pacific/Midway", "Midway"},
	{"Pacific/Port_Moresby", "Port Moresby"},
	{"Australia/Adelaide", "Adelaide"},
	{"Australia/Brisbane", "Brisbane"},
	{"Australia/Darwin", "Darwin"},
	{"Australia/Hobart", "Hobart"},
	{"Australia/Perth", "Perth"},
	{"Australia/Sydney", "Sydney"},
	{"Europe/Amsterdam", "Amsterdam"},
	{"Europe/Athens", "Athens"},
	{"Europe/Berlin", "Berlin"},
	{"Europe/Brussels", "Brussels"},
	{"Europe/Bucharest", "Bucharest"},
	{"Europe/Budapest", "Budapest"},
	{"Europe/Copenhagen", "Copenhagen"},
	{"Europe/Dublin", "Dublin"},
	{"Europe/Helsinki", "Helsinki"},
	{"Europe/Istanbul", "Istanbul"},
	{"Europe/Kiev", "Kyiv"},
	{"Europe/Lisbon", "Lisbon"},
	{"Europe/London", "London"},
	{"Europe/Madrid", "Madrid"},
	{"Europe/Moscow", "Moscow"},
	{"Europe/Oslo", "Oslo"},
	{"Europe/Paris", "Paris"},
	{"Europe/Prague", "Prague"},
	{"Europe/Rome", "Rome"},
	{"Europe/Sofia", "Sofia"},
	{"Europe/Stockholm", "Stockholm"},
	{"Europe/Vienna", "Vienna"},
	{"Europe/Warsaw", "Warsaw"},
	{"Europe/Zurich", "Zurich"},
}

// ListZones returns all zones in the curated list.
func ListZones() []Zone {
	return Zones
}

// CurrentTime returns the current time in the given IANA zone.
func CurrentTime(iana string) (time.Time, error) {
	loc, err := time.LoadLocation(iana)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown timezone %q: %w", iana, err)
	}
	return time.Now().In(loc), nil
}

// LabelForIANA returns the display label for an IANA identifier.
func LabelForIANA(iana string) string {
	for _, z := range Zones {
		if z.IANA == iana {
			return z.Label
		}
	}
	return iana
}
