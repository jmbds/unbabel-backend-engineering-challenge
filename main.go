package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/events"
	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/statistics"
)

/* The entrypoint of our CLI Application */
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

/* An abstraction of the main function to allow error returns */
func run() error {
	/* Retrieve the arguments passed to the application */
	var (
		inputFilepath  string
		outputFilepath string
		windowSize     int
	)

	flag.StringVar(&inputFilepath, "input_file", "events.json", "path to input file containing events")
	flag.StringVar(&outputFilepath, "output_file", "aggregated_events.out.json", "path to aggregated output file")
	flag.IntVar(&windowSize, "window_size", 10, "size of time window for moving average")
	flag.Parse()

	/* Read and parse the input file */
	transactionDeliveredEvents, err := events.ReadAndUnmarshallEventsFile(inputFilepath)
	if err != nil {
		return err
	}

	/* Format the events to an array of TimeframeData for Moving Average calculation */
	eventsGroupedByMinute, err := events.GroupEventsByUnit(transactionDeliveredEvents, time.Minute)
	if err != nil {
		return err
	}

	/* Calculate the Moving Average */
	movingAverage, err := statistics.CalculateMovingAverage(eventsGroupedByMinute, windowSize)
	if err != nil {
		return err
	}

	/* Generate Moving Averate Output as string */
	output, err := events.GenerateMinuteMovingAverageOutput(transactionDeliveredEvents, movingAverage)
	if err != nil {
		return err
	}

	/* Output Moving Average to file */
	return WriteStringToFile(outputFilepath, output)
}
