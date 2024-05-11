package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

/*
A function that reads a file and unmarshalls its contents, line by line, into EventTranslationDelivered structs.

Receives a path to a file containing the list of events.
Returns a list of events and an error.
*/
func ReadAndUnmarshallEventsFile(filepath string) ([]EventTranslationDelivered, error) {
	/*	Open the file for reading	*/
	file, err := os.Open(filepath)
	if err != nil {
		return []EventTranslationDelivered{}, err
	}
	/*	Close the file when exiting	*/
	defer file.Close()

	/*	Initialize a slice for the events	*/
	events := make([]EventTranslationDelivered, 0)

	/*	Initialize the file scanner	*/
	scanner := bufio.NewScanner(file)

	/*	Read file, line by line	*/
	for scanner.Scan() {
		/* Initialize an EventTransationDelivered struct */
		event := EventTranslationDelivered{}

		/*	Unmarshall the scanned bytes into the EventTransationDelivered struct	*/
		err = json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			return []EventTranslationDelivered{}, err
		}

		/*	Append to the event slice	*/
		events = append(events, event)
	}

	return events, nil
}

/*
A function that generates the string output with the desired format.

Receives a list of events and the calculated moving average values.
Returns the string output with date and average_delivery_time and an error.
*/
func GenerateMinuteMovingAverageOutput(events []EventTranslationDelivered, average []float64) (string, error) {
	/*	Parse the timestamp of first event to get initial Date */
	startTimestamp, err := time.Parse(InputTimestampFormat, events[0].Timestamp)
	if err != nil {
		return "", err
	}
	/* Truncate to minute */
	startTimestamp = startTimestamp.Truncate(time.Minute)

	textToOutput := ""
	for i, average := range average {
		/*	Calculate timestamp for current moving average	*/
		timestamp := startTimestamp.Add(time.Duration(i) * 1 * time.Minute)

		/*	Format timestamp into output format	*/
		formattedTimestamp := timestamp.Format(OutputTimestampFormat)

		/*	Generate string with output values */
		textToOutput += fmt.Sprintf("{\"date\": \"%s\", \"average_delivery_time\": %.1f}\n", formattedTimestamp, average)
	}

	return textToOutput, nil
}
