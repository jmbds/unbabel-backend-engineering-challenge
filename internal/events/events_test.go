package events_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/events"
	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/statistics"
)

const InputTimestampFormat string = "2006-01-02 15:04:05.000000"

func TestGroupEventsByUnit(t *testing.T) {

	testcases := []struct {
		name          string
		events        []events.EventTranslationDelivered
		unit          time.Duration
		expected      []statistics.DataPoint
		expectedError error
	}{
		{
			"valid case",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", Duration: 20},
				{Timestamp: "2018-12-26 18:15:19.903159", Duration: 31},
				{Timestamp: "2018-12-26 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			[]statistics.DataPoint{
				{Total: 0, Count: 0},
				{Total: 20, Count: 1},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 31, Count: 1},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 0, Count: 0},
				{Total: 54, Count: 1},
			},
			errors.New(""),
		},
		{
			"invalid case - no dataset",
			[]events.EventTranslationDelivered{},
			time.Minute,
			[]statistics.DataPoint{},
			errors.New("No events found. Please provide a valid list of events."),
		},
		{
			"invalid case - wrong date formats in first event",
			[]events.EventTranslationDelivered{
				{Timestamp: "26-12-2018 18:11:08.509654", Duration: 20},
				{Timestamp: "26-12-2018 18:15:19.903159", Duration: 31},
				{Timestamp: "26-12-2018 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			[]statistics.DataPoint{},
			errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n"),
		},
		{
			"invalid case - wrong date formats in middle event",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", Duration: 20},
				{Timestamp: "26-12-2018 18:15:19.903159", Duration: 31},
				{Timestamp: "2018-12-26 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			[]statistics.DataPoint{},
			errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n"),
		},
		{
			"invalid case -  wrong date format in last event",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", Duration: 20},
				{Timestamp: "2018-12-26 18:15:19.903159", Duration: 31},
				{Timestamp: "26-12-2018 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			[]statistics.DataPoint{},
			errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := events.GroupEventsByUnit(tc.events, tc.unit)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestGetEventWindowByUnit(t *testing.T) {
	testcases := []struct {
		name          string
		events        []events.EventTranslationDelivered
		unit          time.Duration
		expectedStart time.Time
		expectedEnd   time.Time
		expectedError error
	}{
		{
			"valid case",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", Duration: 20},
				{Timestamp: "2018-12-26 18:15:19.903159", Duration: 31},
				{Timestamp: "2018-12-26 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			time.Date(2018, 12, 26, 18, 11, 0, 0, time.UTC),
			time.Date(2018, 12, 26, 18, 24, 0, 0, time.UTC),
			errors.New(""),
		},
		{
			"invalid case - no dataset",
			[]events.EventTranslationDelivered{},
			time.Minute,
			time.Time{},
			time.Time{},
			errors.New("No events found. Please provide a valid list of events."),
		},
		{
			"invalid case - wrong date format in first event",
			[]events.EventTranslationDelivered{
				{Timestamp: "26-12-2018 18:11:08.509654", Duration: 20},
				{Timestamp: "2018-12-26 18:15:19.903159", Duration: 31},
				{Timestamp: "2018-12-26 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			time.Time{},
			time.Time{},
			errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n"),
		},
		{
			"valid case -  wrong date format in middle event",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", Duration: 20},
				{Timestamp: "26-12-2018 18:15:19.903159", Duration: 31},
				{Timestamp: "2018-12-26 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			time.Date(2018, 12, 26, 18, 11, 0, 0, time.UTC),
			time.Date(2018, 12, 26, 18, 24, 0, 0, time.UTC),
			errors.New(""),
		},
		{
			"invalid case -  wrong date format in last event",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", Duration: 20},
				{Timestamp: "2018-12-26 18:15:19.903159", Duration: 31},
				{Timestamp: "26-12-2018 18:23:19.903159", Duration: 54},
			},
			time.Minute,
			time.Time{},
			time.Time{},
			errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotStart, gotEnd, err := events.GetEventWindowByUnit(tc.events, tc.unit)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if gotStart.Compare(tc.expectedStart) != 0 {
				t.Errorf("expected %v, got %v", tc.expectedStart, gotStart)
			}

			if gotEnd.Compare(tc.expectedEnd) != 0 {
				t.Errorf("expected %v, got %v", tc.expectedEnd, gotEnd)
			}
		})
	}
}
