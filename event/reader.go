package event

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	EventID      uint
	CompetitorID uint
	Comment      string
}

func ReadEvent(path string) ([]Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "[") {
			continue
		}

		parts := strings.SplitN(line, "]", 2)
		if len(parts) != 2 {
			continue
		}

		timeStr := strings.TrimPrefix(parts[0], "[")
		fields := strings.Fields(parts[1])

		if len(fields) < 2 {
			continue
		}

		parsedTime, err := time.Parse("15:04:05.000", timeStr)
		if err != nil {
			continue
		}

		EventID, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			continue
		}

		competitorID, err := strconv.ParseUint(fields[1], 10, 32)
		if err != nil {
			continue
		}

		comment := ""
		if len(fields) > 2 {
			comment = strings.Join(fields[2:], " ")
		}

		events = append(events, Event{
			Time:         parsedTime,
			EventID:      uint(EventID),
			CompetitorID: uint(competitorID),
			Comment:      comment,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	return events, nil
}
