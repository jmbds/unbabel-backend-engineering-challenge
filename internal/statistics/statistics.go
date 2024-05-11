package statistics

import (
	"errors"
)

/*
A struct that hold data from a given point in time.

Total is the total value of data we are holding.
Count is the total number of occurrences.
*/
type DataPoint struct {
	Total float64
	Count int
}

/*
A function to calculate the average of a Datapoint.

Returns 0 if the datapoint Count is 0 or lower.
Returns the average if the datapoint Count is bigger than 0.
*/
func (dp *DataPoint) CalculateAverage() float64 {
	/* When count is zero, return 0 */
	if dp.Count <= 0 {
		return 0
	}

	/* Otherwise return the average */
	return dp.Total / float64(dp.Count)
}

/*
A function to calculate the Moving Average given an array of datapoints and a window size.

A datapoint is a struct that contains a value and the number of occurrences.
The window size is the length of previous datapoints used to calculate the moving average.

Returns an array of float values, each one being the moving average in the corresponding window, and an error.
*/
func CalculateMovingAverage(dataPoints []DataPoint, windowSize int) ([]float64, error) {
	/*	When no datapoints are provided, return an error	*/
	if len(dataPoints) == 0 {
		return []float64{}, errors.New("Dataset was empty, please provide a valid dataset.")
	}

	/*	Slice for moving average return values	*/
	movingAverage := make([]float64, 0)

	/*	Struct with Total Duration and NrEvents in Window	*/
	window := DataPoint{Total: 0, Count: 0}

	/*	Queue with last K elements from dataset. The queue size is at most the size of window.	*/
	queue := make([]DataPoint, 0, windowSize)

	/* Struct to hold tail element from queue */
	tail := DataPoint{}

	/*	Iterate through every datapoint. */
	for i := 0; i < len(dataPoints); i++ {
		/*	Append datapoint to queue */
		queue = append(queue, dataPoints[i])

		/*
			Verify if the queue has reached its maximum capacity (windowSize).
			If the queue is full, remove the first element from the queue.
			If the queue isn't full, it implies that there's no data point to remove, so we initialize it with a value of 0.
		*/
		if len(queue) > windowSize {
			tail, queue = queue[0], queue[1:]
		} else {
			tail = DataPoint{}
		}

		/*
			Update the Total Duration and NrEvents in Window.
			We remove the tail element that left the queue and add the datapoint that was appended.
		*/
		window.Total = window.Total - tail.Total + dataPoints[i].Total
		window.Count = window.Count - tail.Count + dataPoints[i].Count

		/*	Append Window average to return slice	*/
		movingAverage = append(movingAverage, window.CalculateAverage())
	}

	return movingAverage, nil
}
