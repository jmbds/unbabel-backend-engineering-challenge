package events_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/events"
)

func TestReadAndUnmarshallEventsFile(t *testing.T) {
	testcases := []struct {
		name          string
		filepath      string
		expected      []events.EventTranslationDelivered
		expectedError error
	}{
		{
			"valid case",
			"testcases/events.json",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654", TranslationId: "5aa5b2f39f7254a75aa5", SourceLanguage: "en", TargetLanguage: "fr", ClientName: "airliberty", EventName: "translation_delivered", Duration: 20, NrWords: 30},
				{Timestamp: "2018-12-26 18:15:19.903159", TranslationId: "5aa5b2f39f7254a75aa4", SourceLanguage: "en", TargetLanguage: "fr", ClientName: "airliberty", EventName: "translation_delivered", Duration: 31, NrWords: 30},
				{Timestamp: "2018-12-26 18:23:19.903159", TranslationId: "5aa5b2f39f7254a75bb3", SourceLanguage: "en", TargetLanguage: "fr", ClientName: "taxi-eats", EventName: "translation_delivered", Duration: 54, NrWords: 100},
			},
			errors.New(""),
		},
		{
			"invalid case - file not found",
			"testcases/not_found.json",
			[]events.EventTranslationDelivered{},
			errors.New("open testcases/not_found.json: no such file or directory"),
		},
		{
			"invalid case - invalid format",
			"testcases/invalid_events.json",
			[]events.EventTranslationDelivered{},
			errors.New("Content is invalid. Please provide a valid events file."),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := events.ReadAndUnmarshallEventsFile(tc.filepath)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error: %s", err.Error())
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestGenerateMinuteMovingAverageOutput(t *testing.T) {
	testcases := []struct {
		name          string
		events        []events.EventTranslationDelivered
		average       []float64
		expected      string
		expectedError error
	}{
		{
			"valid case",
			[]events.EventTranslationDelivered{
				{Timestamp: "2018-12-26 18:11:08.509654"},
				{Timestamp: "2018-12-26 18:12:19.903159"},
			},
			[]float64{0, 20, 20},
			"{\"date\": \"2018-12-26 18:11:00\", \"average_delivery_time\": 0.0}\n{\"date\": \"2018-12-26 18:12:00\", \"average_delivery_time\": 20.0}\n{\"date\": \"2018-12-26 18:13:00\", \"average_delivery_time\": 20.0}\n",
			errors.New(""),
		},
		{
			"invalid case - wrong timestamp format",
			[]events.EventTranslationDelivered{
				{Timestamp: "21-07-2016 14:51:08.509654"},
				{Timestamp: "2016-07-22 14:52:19.903159"},
			},
			[]float64{31, 50, 88},
			"",
			errors.New("Invalid date format. Please provide dates in the following format: " + InputTimestampFormat + "\n"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := events.GenerateMinuteMovingAverageOutput(tc.events, tc.average)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
