package club

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{"Valid Time", "15:04", time.Date(0, 1, 1, 15, 4, 0, 0, time.UTC), false},
		{"Invalid Time", "25:61", time.Time{}, true},
		{"Empty Input", "", time.Time{}, true},
		{"Missing Minute", "15:", time.Time{}, true},
		{"Extra Characters", "15:04pm", time.Time{}, true},
		{"Invalid Format", "3 PM", time.Time{}, true},
		{"Only Hour", "15", time.Time{}, true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("parseTime(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
