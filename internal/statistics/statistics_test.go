package statistics_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jmbds/unbabel-backend-engineering-challenge/internal/statistics"
)

func TestCalculateAverage(t *testing.T) {
	testcases := []struct {
		name          string
		timeframeData statistics.DataPoint
		expected      float64
	}{
		{"average", statistics.DataPoint{Total: 40, Count: 10}, 4},
		{"average_for_zero", statistics.DataPoint{Total: 40, Count: 0}, 0},
		{"average_for_negative", statistics.DataPoint{Total: 40, Count: -10}, 0},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timeframeData.CalculateAverage()

			if got != tc.expected {
				t.Errorf("expected %.1f, got %1.f", tc.expected, got)
			}
		})
	}
}

func TestCalculateMovingAverage(t *testing.T) {
	testcases := []struct {
		name          string
		dataset       []statistics.DataPoint
		windowSize    int
		expected      []float64
		expectedError error
	}{
		{
			"valid case",
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
			10,
			[]float64{0, 20, 20, 20, 20, 25.5, 25.5, 25.5, 25.5, 25.5, 25.5, 31, 31, 42.5},
			errors.New(""),
		},
		{
			"invalid case - no dataset",
			[]statistics.DataPoint{},
			10,
			[]float64{},
			errors.New("Dataset was empty, please provide a valid dataset."),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := statistics.CalculateMovingAverage(tc.dataset, tc.windowSize)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
