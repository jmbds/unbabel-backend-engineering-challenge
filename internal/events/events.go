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

/*
A function that aggregates the duration of events by time unit.
Receives a list of events and a unit of time.

Returns a list of data points, with the total duration of events and number of occurrences aggregated by time unit.
*/
func GroupEventsByUnit(events []EventTranslationDelivered, unit time.Duration) ([]statistics.DataPoint, error) {
	/* Get the start and finish timestamps aka window. */
	windowStart, windowEnd, err := GetEventWindowByUnit(events, unit)
	if err != nil {
		return []statistics.DataPoint{}, err
	}

	/*	Calculate number of time units in event input. */
	datasetLength := calculateUnitDifference(windowStart, windowEnd, unit) + 1

	/*	Initialize dataset with 0 values.	*/
	dataset := make([]statistics.DataPoint, datasetLength, datasetLength)

	/* Iterate through events and group them by minute difference to window start.	*/
	for _, event := range events {
		timestamp, err := time.Parse(InputTimestampFormat, event.Timestamp)
		if err != nil {
			return []statistics.DataPoint{}, errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n")
		}

		/*	Include an extra unit (second, minute, ...) in the calculation, as events are logged based on the unit immediately following their occurrence.	*/
		timestamp = timestamp.Add(1 * unit).Truncate(unit)

		/* Calculate the corresponding index for the event, based on time unit difference to the window start. */
		index := calculateUnitDifference(windowStart, timestamp, unit)

		/*	Increment the number of Events and total Duration for specific index	*/
		dataset[index] = statistics.DataPoint{
			Total: dataset[index].Total + float64(event.Duration),
			Count: dataset[index].Count + 1,
		}
	}
	return dataset, nil
}

/*
A function that calculates the event window by unit of time.
Receives a list of events and a unit of time.

Returns a start time and an end time according to the unit of time given.
*/
func GetEventWindowByUnit(events []EventTranslationDelivered, unit time.Duration) (time.Time, time.Time, error) {
	/* Check if there are events to calculate the time window */
	nrEvents := len(events)
	if nrEvents == 0 {
		return time.Time{}, time.Time{}, errors.New("No events found. Please provide a valid list of events.")
	}

	/* Since events are ordered, grab the first and last event of the list */
	initialEvent, finalEvent := events[0], events[nrEvents-1]

	/* Parse the timestamp of first event */
	initialEventTimestamp, err := time.Parse(InputTimestampFormat, initialEvent.Timestamp)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n")
	}
	/* Truncate to time unit */
	initialEventTimestamp = initialEventTimestamp.Truncate(1 * unit)

	/* Parse the timestamp of last event */
	finalEventTimestamp, err := time.Parse(InputTimestampFormat, finalEvent.Timestamp)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n")
	}
	/* Add 1 unit, since the last event will only be counted in next time unit, and truncate to time unit */
	finalEventTimestamp = finalEventTimestamp.Add(1 * unit).Truncate(1 * unit)

	/* Return initial and start dates [aka event window] */
	return initialEventTimestamp, finalEventTimestamp, nil
}

/*
Function to calculate the time difference in provided unit of time
*/
func calculateUnitDifference(start time.Time, end time.Time, unit time.Duration) int {
	return int(end.Sub(start) / unit)
}
