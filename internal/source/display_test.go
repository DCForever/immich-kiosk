package source

import (
	"testing"
	"time"
)

func TestIsoWeekRange(t *testing.T) {
	utc := time.UTC
	tests := []struct {
		name     string
		input    time.Time
		wantMon  time.Time
		wantSun  time.Time
	}{
		{
			name:    "Tuesday March 13 2007",
			input:   time.Date(2007, 3, 13, 12, 0, 0, 0, utc),
			wantMon: time.Date(2007, 3, 12, 0, 0, 0, 0, utc),
			wantSun: time.Date(2007, 3, 18, 0, 0, 0, 0, utc),
		},
		{
			name:    "Monday March 12 2007",
			input:   time.Date(2007, 3, 12, 0, 0, 0, 0, utc),
			wantMon: time.Date(2007, 3, 12, 0, 0, 0, 0, utc),
			wantSun: time.Date(2007, 3, 18, 0, 0, 0, 0, utc),
		},
		{
			name:    "Sunday March 18 2007",
			input:   time.Date(2007, 3, 18, 23, 59, 59, 0, utc),
			wantMon: time.Date(2007, 3, 12, 0, 0, 0, 0, utc),
			wantSun: time.Date(2007, 3, 18, 0, 0, 0, 0, utc),
		},
		{
			name:    "Monday Jan 1 2007",
			input:   time.Date(2007, 1, 1, 0, 0, 0, 0, utc),
			wantMon: time.Date(2007, 1, 1, 0, 0, 0, 0, utc),
			wantSun: time.Date(2007, 1, 7, 0, 0, 0, 0, utc),
		},
		{
			name:    "Sunday Jan 1 2023 (ISO week in previous year)",
			input:   time.Date(2023, 1, 1, 0, 0, 0, 0, utc),
			wantMon: time.Date(2022, 12, 26, 0, 0, 0, 0, utc),
			wantSun: time.Date(2023, 1, 1, 0, 0, 0, 0, utc),
		},
		{
			name:    "Wednesday mid-year",
			input:   time.Date(2023, 7, 19, 14, 30, 0, 0, utc),
			wantMon: time.Date(2023, 7, 17, 0, 0, 0, 0, utc),
			wantSun: time.Date(2023, 7, 23, 0, 0, 0, 0, utc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMon, gotSun := IsoWeekRange(tt.input)
			if !gotMon.Equal(tt.wantMon) {
				t.Errorf("IsoWeekRange() start = %v, want %v", gotMon, tt.wantMon)
			}
			if !gotSun.Equal(tt.wantSun) {
				t.Errorf("IsoWeekRange() end = %v, want %v", gotSun, tt.wantSun)
			}
			if gotMon.Weekday() != time.Monday {
				t.Errorf("IsoWeekRange() start weekday = %v, want Monday", gotMon.Weekday())
			}
			if gotSun.Weekday() != time.Sunday {
				t.Errorf("IsoWeekRange() end weekday = %v, want Sunday", gotSun.Weekday())
			}
		})
	}
}
