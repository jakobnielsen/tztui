package tz

import (
	"testing"
	"time"
	_ "time/tzdata"
)

func TestCurrentTime(t *testing.T) {
	tests := []struct {
		iana    string
		wantErr bool
	}{
		{"UTC", false},
		{"Europe/London", false},
		{"America/New_York", false},
		{"Asia/Tokyo", false},
		{"Not/A/Zone", true},
	}
	for _, tc := range tests {
		t.Run(tc.iana, func(t *testing.T) {
			got, err := CurrentTime(tc.iana)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got nil", tc.iana)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.IsZero() {
				t.Fatal("got zero time")
			}
		})
	}
}

func TestListZones_NonEmpty(t *testing.T) {
	zones := ListZones()
	if len(zones) == 0 {
		t.Fatal("ListZones returned empty slice")
	}
}

func TestAllZonesParseable(t *testing.T) {
	for _, z := range Zones {
		t.Run(z.IANA, func(t *testing.T) {
			_, err := time.LoadLocation(z.IANA)
			if err != nil {
				t.Errorf("time.LoadLocation(%q) failed: %v", z.IANA, err)
			}
		})
	}
}

func TestLabelForIANA(t *testing.T) {
	if got := LabelForIANA("Europe/London"); got != "London" {
		t.Errorf("got %q, want %q", got, "London")
	}
	if got := LabelForIANA("Unknown/Zone"); got != "Unknown/Zone" {
		t.Errorf("fallback: got %q, want %q", got, "Unknown/Zone")
	}
}
