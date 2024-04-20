package pkg

import (
	"sort"
	"time"
)

type TimeSeries[T any] struct {
	Datapoints []T
	Timestamps []time.Time
}

func NewTimeSeriesFromSortedDataPoints[T any](datapoints []T, timestamps []time.Time) TimeSeries[T] {
	return TimeSeries[T]{
		Datapoints: datapoints,
		Timestamps: timestamps,
	}
}

func (t *TimeSeries[T]) AddUnsortedDataPoint(datapoint T, timestamp time.Time) {
	// Find the index of the first timestamp that is greater than or equal to the given timestamp
	index := sort.Search(len(t.Timestamps), func(i int) bool {
		return t.Timestamps[i].After(timestamp)
	})

	// Insert the new data point and timestamp into the slice
	t.Datapoints = append(t.Datapoints[:index+1], t.Datapoints[index:]...)
	t.Datapoints[index] = datapoint

	t.Timestamps = append(t.Timestamps[:index+1], t.Timestamps[index:]...)
	t.Timestamps[index] = timestamp
}

func (t *TimeSeries[T]) GetDataPointsWithin(start time.Time, end time.Time) TimeSeries[T] {
	// Find the index of the first timestamp that is greater than or equal to the given start time
	startIndex := sort.Search(len(t.Timestamps), func(i int) bool {
		return t.Timestamps[i].After(start)
	})

	// Find the index of the first timestamp that is greater than or equal to the given end time
	endIndex := sort.Search(len(t.Timestamps), func(i int) bool {
		return t.Timestamps[i].After(end)
	})

	return NewTimeSeriesFromSortedDataPoints(
		t.Datapoints[startIndex:endIndex],
		t.Timestamps[startIndex:endIndex],
	)
}
