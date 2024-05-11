package events_test

import (
	"errors"
	"testing"

	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/events"
)

func TestReadAndUnmarshallEventsFile(t *testing.T) {

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
			"valid case",
			[]events.EventTranslationDelivered{
				{Timestamp: "2016-07-21 14:51:08.509654"},
				{Timestamp: "2016-07-22 14:52:19.903159"},
			},
			[]float64{31, 50, 88},
			"{\"date\": \"2016-07-21 14:51:00\", \"average_delivery_time\": 31.0}\n{\"date\": \"2016-07-21 14:52:00\", \"average_delivery_time\": 50.0}\n{\"date\": \"2016-07-21 14:53:00\", \"average_delivery_time\": 88.0}\n",
			errors.New(""),
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
