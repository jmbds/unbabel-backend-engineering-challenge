package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

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
			return []EventTranslationDelivered{}, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GenerateMinuteMovingAverageOutput(events []EventTranslationDelivered, average []float64) (string, error) {
	startTimestamp, err := time.Parse(InputTimestampFormat, events[0].Timestamp)
	if err != nil {
		return "", err
	}
	startTimestamp = startTimestamp.Truncate(time.Minute)

	textToOutput := ""
	for i, average := range average {
		/*	Calculate timestamp for current average	*/
		timestamp := startTimestamp.Add(time.Duration(i) * 1 * time.Minute)

		/*	Format timestamp into output format	*/
		formattedTimestamp := timestamp.Format(OutputTimestampFormat)

		/*	Generate string with output values */
		textToOutput += fmt.Sprintf("{\"date\": \"%s\", \"average_delivery_time\": %.1f}\n", formattedTimestamp, average)
	}

	return textToOutput, nil
}
