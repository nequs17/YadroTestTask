package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type configJSON struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type Config struct {
	Laps        int           `json:"laps"`
	LapLen      int           `json:"lapLen"`
	PenaltyLen  int           `json:"penaltyLen"`
	FiringLines int           `json:"firingLines"`
	Start       time.Time     `json:"start"`
	StartDelta  time.Duration `json:"startDelta"`
}

func ReadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cj configJSON

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cj)
	if err != nil {
		return Config{}, err
	}

	start, err := time.Parse("15:04:05.000", cj.Start)
	if err != nil {
		return Config{}, err
	}

	var h, m, s int
	_, err = fmt.Sscanf(cj.StartDelta, "%02d:%02d:%02d", &h, &m, &s)
	if err != nil {
		return Config{}, fmt.Errorf("parse error startDelta: %w", err)
	}
	startDelta := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second

	config := Config{
		Laps:        cj.Laps,
		LapLen:      cj.LapLen,
		PenaltyLen:  cj.PenaltyLen,
		FiringLines: cj.FiringLines,
		Start:       start,
		StartDelta:  startDelta,
	}

	return config, nil
}
