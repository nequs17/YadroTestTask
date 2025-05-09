package event

import (
	"fmt"
	"os"
	cfg "skiers/config"

	"time"
)

type CompetitorState struct {
	ScheduledStartTime time.Time
	ActualStartTime    time.Time
	FinishTime         time.Time
	LapEndTimes        []time.Time
	ShootingStages     []ShootingStage
	PenaltySessions    []PenaltySession
	Disqualified       bool
	Comment            string
}

type ShootingStage struct {
	Hits       int
	TotalShots int
}

type PenaltySession struct {
	StartTime time.Time
	EndTime   time.Time
}

// nolint:gocyclo // all's good, it's test project
func EventLog(c *cfg.Config, path string, events []Event) (map[uint]*CompetitorState, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	competitors := make(map[uint]*CompetitorState)

	for _, e := range events {
		if _, exists := competitors[e.CompetitorID]; !exists {
			competitors[e.CompetitorID] = &CompetitorState{}
		}
		state := competitors[e.CompetitorID]

		switch e.EventID {
		case 1:
			fmt.Fprintf(file, "[%s] The competitor(%d) registered\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 2:
			scheduledStart, err := time.Parse("15:04:05.000", e.Comment)
			if err == nil {
				state.ScheduledStartTime = scheduledStart
			}
			fmt.Fprintf(file, "[%s] The start time for the competitor(%d) was set by a draw to %s\n", e.Time.Format("15:04:05.000"), e.CompetitorID, e.Comment)

		case 3:
			fmt.Fprintf(file, "[%s] The competitor(%d) is on the start line\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 4:
			state.ActualStartTime = e.Time
			fmt.Fprintf(file, "[%s] The competitor(%d) has started\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 5:
			newStage := ShootingStage{Hits: 0, TotalShots: 5}
			state.ShootingStages = append(state.ShootingStages, newStage)
			fmt.Fprintf(file, "[%s] The competitor(%d) is on the firing range(%s)\n", e.Time.Format("15:04:05.000"), e.CompetitorID, e.Comment)

		case 6:
			if len(state.ShootingStages) > 0 {
				lastStage := &state.ShootingStages[len(state.ShootingStages)-1]
				lastStage.Hits++
			}
			fmt.Fprintf(file, "[%s] The target(%s) has been hit by competitor(%d)\n", e.Time.Format("15:04:05.000"), e.Comment, e.CompetitorID)

		case 7:
			fmt.Fprintf(file, "[%s] The competitor(%d) left the firing range\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 8:
			newPenalty := PenaltySession{StartTime: e.Time}
			state.PenaltySessions = append(state.PenaltySessions, newPenalty)
			fmt.Fprintf(file, "[%s] The competitor(%d) entered the penalty laps\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 9:
			if len(state.PenaltySessions) > 0 {
				lastPenalty := &state.PenaltySessions[len(state.PenaltySessions)-1]
				lastPenalty.EndTime = e.Time
			}
			fmt.Fprintf(file, "[%s] The competitor(%d) left the penalty laps\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 10:
			state.LapEndTimes = append(state.LapEndTimes, e.Time)
			if len(state.LapEndTimes) == c.Laps {
				state.FinishTime = e.Time
			}
			fmt.Fprintf(file, "[%s] The competitor(%d) ended the main lap\n", e.Time.Format("15:04:05.000"), e.CompetitorID)

		case 11:
			state.Disqualified = true
			state.Comment = e.Comment
			fmt.Fprintf(file, "[%s] The competitor(%d) can't continue: %s\n", e.Time.Format("15:04:05.000"), e.CompetitorID, e.Comment)
		}
	}

	return competitors, nil
}
