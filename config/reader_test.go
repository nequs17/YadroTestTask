package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantErr    bool
		wantConfig Config
	}{
		{
			name:    "Valid Config",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "12:00:00.000", "startDelta": "00:01:00"}`,
			wantErr: false,
			wantConfig: Config{
				Laps:        3,
				LapLen:      1000,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       mustParseTime("15:04:05.000", "12:00:00.000"),
				StartDelta:  time.Minute,
			},
		},
		{
			name:    "File Not Found",
			content: "",
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "12:00:00.000", "startDelta": "00:01:00"`,
			wantErr: true,
		},
		{
			name:    "Missing Start",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "startDelta": "00:01:00"}`,
			wantErr: true,
		},
		{
			name:    "Invalid Start Format",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "invalid", "startDelta": "00:01:00"}`,
			wantErr: true,
		},
		{
			name:    "Missing StartDelta",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "12:00:00.000"}`,
			wantErr: true,
		},
		{
			name:    "Invalid StartDelta Format",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "12:00:00.000", "startDelta": "invalid"}`,
			wantErr: true,
		},
		{
			name:    "Type Mismatch",
			content: `{"laps": "three", "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "12:00:00.000", "startDelta": "00:01:00"}`,
			wantErr: true,
		},
		{
			name:    "Extra Fields",
			content: `{"laps": 3, "lapLen": 1000, "penaltyLen": 150, "firingLines": 2, "start": "12:00:00.000", "startDelta": "00:01:00", "extra": "field"}`,
			wantErr: false,
			wantConfig: Config{
				Laps:        3,
				LapLen:      1000.0,
				PenaltyLen:  150.0,
				FiringLines: 2,
				Start:       mustParseTime("15:04:05.000", "12:00:00.000"),
				StartDelta:  time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			var path string
			if tt.content != "" {
				filePath := filepath.Join(tempDir, "config.json")
				err := os.WriteFile(filePath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatal(err)
				}
				path = filePath
			} else {
				path = filepath.Join(tempDir, "nonexistent.json")
			}

			gotConfig, err := ReadConfig(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
					t.Errorf("ReadConfig() = %v, want %v", gotConfig, tt.wantConfig)
				}
			}
		})
	}
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(fmt.Sprintf("failed to parse time: %v", err))
	}
	return t
}
