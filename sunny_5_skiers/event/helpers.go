package event

import (
	"fmt"
	cfg "skiers/config"
	"strings"
	"time"
)

func formatDuration(d time.Duration) string {
	hours := d / time.Hour
	remaining := d % time.Hour
	minutes := remaining / time.Minute
	remaining = remaining % time.Minute
	seconds := remaining / time.Second
	remaining = remaining % time.Second
	milliseconds := remaining / time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func getLapPart(state *CompetitorState, c *cfg.Config) string {
	lapStr := make([]string, c.Laps)
	for i := range lapStr {
		if i < len(state.LapEndTimes) {
			var startTime time.Time
			if i == 0 {
				startTime = state.ActualStartTime
			} else {
				startTime = state.LapEndTimes[i-1]
			}
			endTime := state.LapEndTimes[i]

			if startTime.IsZero() {
				lapStr[i] = "{,}"
				continue
			}

			lapTime := endTime.Sub(startTime)
			lapSpeed := float64(c.LapLen) / lapTime.Seconds()
			lapStr[i] = fmt.Sprintf("{%s, %.3f}", formatDuration(lapTime), lapSpeed)
		} else {
			lapStr[i] = "{,}"
		}
	}
	return strings.Join(lapStr, ", ")
}

func getShootingSummary(state *CompetitorState) (int, int) {
	totalHits := 0
	totalShots := 0
	for _, stage := range state.ShootingStages {
		totalHits += stage.Hits
		totalShots += stage.TotalShots
	}
	return totalHits, totalShots
}

func getPenaltyPart(state *CompetitorState, c *cfg.Config, totalMisses int) string {
	totalPenaltyDistance := float64(totalMisses) * float64(c.PenaltyLen)
	totalPenaltyTime := time.Duration(0)
	for _, ps := range state.PenaltySessions {
		if !ps.EndTime.IsZero() {
			totalPenaltyTime += ps.EndTime.Sub(ps.StartTime)
		}
	}
	if totalPenaltyTime > 0 {
		averagePenaltySpeed := totalPenaltyDistance / totalPenaltyTime.Seconds()
		return fmt.Sprintf("{%s, %.3f}", formatDuration(totalPenaltyTime), averagePenaltySpeed)
	}
	return "{00:00:00.000, 0.000}"
}
