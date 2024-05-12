package events

import (
	"bufio"
	"encoding/json"
	"errors"
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
	file, err := os.Open(filepath)
	if err != nil {
		return []EventTranslationDelivered{}, err
	}
	defer file.Close()

	events := make([]EventTranslationDelivered, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		event := EventTranslationDelivered{}

		err = json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			return []EventTranslationDelivered{}, errors.New("Content is invalid. Please provide a valid events file.")
		}

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
	startTimestamp, err := time.Parse(InputTimestampFormat, events[0].Timestamp)
	if err != nil {
		return "", errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n")
	}
	startTimestamp = startTimestamp.Truncate(time.Minute)

	textToOutput := ""
	for i, average := range average {
		/*	Calculate timestamp for current moving average	*/
		timestamp := startTimestamp.Add(time.Duration(i) * 1 * time.Minute)
		formattedTimestamp := timestamp.Format(OutputTimestampFormat)

		/* Check if we should remove decimal places of float value */
		if average == float64(int(average)) {
			textToOutput += fmt.Sprintf("{\"date\": \"%s\", \"average_delivery_time\": %d}\n", formattedTimestamp, int(average))
		} else {
			textToOutput += fmt.Sprintf("{\"date\": \"%s\", \"average_delivery_time\": %.1f}\n", formattedTimestamp, average)
		}

	}

	return textToOutput, nil
}
