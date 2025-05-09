package event

import (
	"fmt"
	"os"
	cfg "skiers/config"
)

func WriteEvent(c *cfg.Config, path string, competitors map[uint]*CompetitorState) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for id, state := range competitors {
		var status string
		var displayFirst string
		if len(state.LapEndTimes) == c.Laps {
			status = "Finished"
			totalTime := state.FinishTime.Sub(state.ActualStartTime)
			displayFirst = formatDuration(totalTime)
		} else if state.Disqualified {
			status = "NotFinished"
			displayFirst = status
		} else if state.ActualStartTime.IsZero() {
			status = "NotStarted"
			displayFirst = status
		} else {
			status = "NotFinished"
			displayFirst = status
		}

		lapPart := getLapPart(state, c)
		totalHits, totalShots := getShootingSummary(state)
		totalMisses := totalShots - totalHits
		shootingPart := fmt.Sprintf("%d/%d", totalHits, totalShots)
		penaltyPart := getPenaltyPart(state, c, totalMisses)

		fmt.Fprintf(file, "[%s] %d [%s] %s %s\n", displayFirst, id, lapPart, penaltyPart, shootingPart)
	}

	return nil
}
