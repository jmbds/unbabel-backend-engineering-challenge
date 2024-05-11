package events

import (
	"errors"
	"time"

	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/statistics"
)

const (
	InputTimestampFormat  = "2006-01-02 15:04:05.000000"
	OutputTimestampFormat = "2006-01-02 15:04:05"
)

type EventTranslationDelivered struct {
	Timestamp      string `json:"timestamp"`
	TranslationId  string `json:"translation_id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	ClientName     string `json:"client_name"`
	EventName      string `json:"event_name"`
	Duration       int    `json:"duration"`
	NrWords        int    `json:"nr_words"`
}

func GroupEventsByUnit(events []EventTranslationDelivered, unit time.Duration) ([]statistics.DataPoint, error) {
	/* Get the start and finish timestamps aka window. */
	windowStart, windowEnd, err := GetEventWindowByUnit(events, unit)
	if err != nil {
		return []statistics.DataPoint{}, err
	}

	/*	Calculate number of minutes in event input. */
	datasetLength := calculateUnitDifference(windowStart, windowEnd, unit) + 1

	/*	Initialize dataset with 0 values.	*/
	dataset := make([]statistics.DataPoint, datasetLength, datasetLength)

	/* Iterate through events and group them by minute difference to window start.	*/
	for _, event := range events {
		timestamp, err := time.Parse(InputTimestampFormat, event.Timestamp)
		if err != nil {
			return []statistics.DataPoint{}, err
		}

		/*	Include an extra unit (second, minute, ...) in the calculation, as events are logged based on the unit immediately following their occurrence.	*/
		timestamp = timestamp.Add(1 * unit).Truncate(unit)

		/* Calculate the corresponding index for the event, based on minute difference to the window start. */
		index := calculateUnitDifference(windowStart, timestamp, unit)

		/*	Increment the number of Events and total Duration for specific index	*/
		dataset[index] = statistics.DataPoint{
			Total: dataset[index].Total + float64(event.Duration),
			Count: dataset[index].Count + 1,
		}
	}
	return dataset, nil
}

func GetEventWindowByUnit(events []EventTranslationDelivered, unit time.Duration) (time.Time, time.Time, error) {
	nrEvents := len(events)
	if nrEvents == 0 {
		return time.Time{}, time.Time{}, errors.New("No events found. Please provide a valid list of events.")
	}

	initialEvent, finalEvent := events[0], events[nrEvents-1]

	initialEventTimestamp, err := time.Parse(InputTimestampFormat, initialEvent.Timestamp)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	initialEventTimestamp = initialEventTimestamp.Truncate(1 * unit)

	finalEventTimestamp, err := time.Parse(InputTimestampFormat, finalEvent.Timestamp)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	finalEventTimestamp = finalEventTimestamp.Add(1 * unit).Truncate(1 * unit)

	return initialEventTimestamp, finalEventTimestamp, nil
}

func calculateUnitDifference(start time.Time, end time.Time, unit time.Duration) int {
	return int(end.Sub(start) / unit)
}
