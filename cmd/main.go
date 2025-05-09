package main

import (
	"flag"
	"log"
	"os"
	cfg "skiers/config"
	event "skiers/event"
)

var configPath string
var eventsPath string
var logPath string
var outputPath string

func main() {

	flag.StringVar(&configPath, "config", "config.json", "path to config file")
	flag.StringVar(&eventsPath, "events", "events", "path to events file")
	flag.StringVar(&logPath, "log", "output.log", "path to log file")
	flag.StringVar(&outputPath, "output", "results.txt", "path to output file")
	flag.Parse()

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	config, err := cfg.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	events, err := event.ReadEvent(eventsPath)
	if err != nil {
		log.Fatal(err)
	}

	originalStdout := os.Stdout
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("Failed to create output file:", err)
	}
	defer outputFile.Close()
	os.Stdout = outputFile

	states, err := event.EventLog(&config, logPath, events)
	if err != nil {
		log.Fatal(err)
	}

	err = event.WriteEvent(&config, outputPath, states)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout = originalStdout

}
