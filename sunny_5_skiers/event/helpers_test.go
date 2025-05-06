package event

import (
	"fmt"
	cfg "skiers/config"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "Zero Duration",
			duration: 0,
			want:     "00:00:00.000",
		},
		{
			name:     "One Minute",
			duration: time.Minute,
			want:     "00:01:00.000",
		},
		{
			name:     "Complex Duration",
			duration: 2*time.Hour + 30*time.Minute + 45*time.Second + 123*time.Millisecond,
			want:     "02:30:45.123",
		},
		{
			name:     "Only Milliseconds",
			duration: 456 * time.Millisecond,
			want:     "00:00:00.456",
		},
		{
			name:     "Negative Duration",
			duration: -1*time.Hour - 15*time.Minute,
			want:     "-1:-15:00.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestGetLapPart(t *testing.T) {
	config := &cfg.Config{
		Laps:   3,
		LapLen: 1000, // 1000 meters
	}

	startTime := mustParseTime("15:04:05.000", "12:00:00.000")
	lap1End := startTime.Add(10 * time.Minute) // 10 minutes for lap 1
	lap2End := lap1End.Add(12 * time.Minute)   // 12 minutes for lap 2

	tests := []struct {
		name   string
		state  *CompetitorState
		config *cfg.Config
		want   string
	}{
		{
			name: "All Laps Completed",
			state: &CompetitorState{
				ActualStartTime: startTime,
				LapEndTimes:     []time.Time{lap1End, lap2End, lap2End.Add(15 * time.Minute)},
			},
			config: config,
			want:   "{00:10:00.000, 1.667}, {00:12:00.000, 1.389}, {00:15:00.000, 1.111}",
		},
		{
			name: "Partial Laps",
			state: &CompetitorState{
				ActualStartTime: startTime,
				LapEndTimes:     []time.Time{lap1End},
			},
			config: config,
			want:   "{00:10:00.000, 1.667}, {,}, {,}",
		},
		{
			name: "No Laps",
			state: &CompetitorState{
				ActualStartTime: startTime,
				LapEndTimes:     []time.Time{},
			},
			config: config,
			want:   "{,}, {,}, {,}",
		},
		{
			name: "Zero Start Time",
			state: &CompetitorState{
				ActualStartTime: time.Time{},
				LapEndTimes:     []time.Time{lap1End},
			},
			config: config,
			want:   "{,}, {,}, {,}", // Speed is approximate due to zero start time
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getLapPart(tt.state, tt.config)
			if got != tt.want {
				t.Errorf("getLapPart() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetShootingSummary(t *testing.T) {
	tests := []struct {
		name  string
		state *CompetitorState
		want  struct{ hits, shots int }
	}{
		{
			name: "No Shooting Stages",
			state: &CompetitorState{
				ShootingStages: []ShootingStage{},
			},
			want: struct{ hits, shots int }{0, 0},
		},
		{
			name: "Single Stage",
			state: &CompetitorState{
				ShootingStages: []ShootingStage{
					{Hits: 4, TotalShots: 5},
				},
			},
			want: struct{ hits, shots int }{4, 5},
		},
		{
			name: "Multiple Stages",
			state: &CompetitorState{
				ShootingStages: []ShootingStage{
					{Hits: 4, TotalShots: 5},
					{Hits: 3, TotalShots: 5},
				},
			},
			want: struct{ hits, shots int }{7, 10},
		},
		{
			name: "Zero Hits",
			state: &CompetitorState{
				ShootingStages: []ShootingStage{
					{Hits: 0, TotalShots: 5},
				},
			},
			want: struct{ hits, shots int }{0, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHits, gotShots := getShootingSummary(tt.state)
			if gotHits != tt.want.hits || gotShots != tt.want.shots {
				t.Errorf("getShootingSummary() = (%d, %d), want (%d, %d)", gotHits, gotShots, tt.want.hits, tt.want.shots)
			}
		})
	}
}

func TestGetPenaltyPart(t *testing.T) {
	config := &cfg.Config{
		PenaltyLen: 150,
	}

	startTime := mustParseTime("15:04:05.000", "12:00:00.000")
	penaltyStart := startTime.Add(5 * time.Minute)
	penaltyEnd := penaltyStart.Add(30 * time.Second)

	tests := []struct {
		name        string
		state       *CompetitorState
		config      *cfg.Config
		totalMisses int
		want        string
	}{
		{
			name:        "No Penalties",
			state:       &CompetitorState{PenaltySessions: []PenaltySession{}},
			config:      config,
			totalMisses: 0,
			want:        "{00:00:00.000, 0.000}",
		},
		{
			name: "Single Penalty",
			state: &CompetitorState{
				PenaltySessions: []PenaltySession{
					{StartTime: penaltyStart, EndTime: penaltyEnd},
				},
			},
			config:      config,
			totalMisses: 2,
			want:        "{00:00:30.000, 10.000}", // 300 meters / 30 seconds = 10 m/s
		},
		{
			name: "Multiple Penalties",
			state: &CompetitorState{
				PenaltySessions: []PenaltySession{
					{StartTime: penaltyStart, EndTime: penaltyEnd},
					{StartTime: penaltyEnd, EndTime: penaltyEnd.Add(20 * time.Second)},
				},
			},
			config:      config,
			totalMisses: 3,
			want:        "{00:00:50.000, 9.000}", // 450 meters / 50 seconds = 9 m/s
		},
		{
			name: "Incomplete Penalty Session",
			state: &CompetitorState{
				PenaltySessions: []PenaltySession{
					{StartTime: penaltyStart, EndTime: time.Time{}},
				},
			},
			config:      config,
			totalMisses: 1,
			want:        "{00:00:00.000, 0.000}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPenaltyPart(tt.state, tt.config, tt.totalMisses)
			if got != tt.want {
				t.Errorf("getPenaltyPart() = %q, want %q", got, tt.want)
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
