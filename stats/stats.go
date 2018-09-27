/*
Package stats implements an online computation of summary statistics.

The primary idea of this package is that samples are coming in at real time,
and online computations of the shape of the distribution: the mean, variance,
and range need to be computed on-demand. Rather than keeping an array of
values, online algorithms update the internal state of the descriptive
statistics at runtime, saving memory.

To track statistics in an online fashion, you need to keep track of the
various aggregates that are used to compute the final descriptives statistics
of the distribution. For simple statistics such as the minimum, maximum,
standard deviation, and mean you need to track the number of samples, the sum
of samples, and the sum of the squares of all samples (along with the minimum
and maximum value seen).

The primary entry point into this function is the Update method, where you
can pass sample values and retrieve data back. All other methods are simply
computations for values.
*/
package stats

import (
	"math"
	"sync"
)

// Statistics keeps track of descriptive statistics in an online fashion at
// runtime without saving each individual sample in an array. It does this by
// updating the internal state of summary aggregates including the number of
// samples seen, the sum of values, and the sum of the value squared. It also
// tracks the minimum and maximum values seen.
//
// The primary entry point to the object is via the Update method, where one
// or more samples can be passed. This object has unexported fields because
// it is thread-safe (via a sync.RWMutex). All properties must be accesesd
// from read-locked access methods.
type Statistics struct {
	sync.RWMutex
	samples uint64  // number of samples seen
	total   float64 // the sum of all samples
	squares float64 // the sum of the squares of each sample
	maximum float64 // the maximum sample observed
	minimum float64 // the minimum sample observed
}

// Update the statistics with a sample or samples (thread-safe). Note that
// this object expects float64 values. While statistical computations for
// integer values are possible, it is simpler to simply transform the values
// into floats ahead of time.
func (s *Statistics) Update(samples ...float64) {
	s.Lock()
	defer s.Unlock()

	for _, sample := range samples {
		s.samples++
		s.total += sample
		s.squares += (sample * sample)

		// If this is our first sample then this value is both our maximum and
		// our minimum value. Otherwise, perform comparisions.
		if s.samples == 1 {
			s.maximum = sample
			s.minimum = sample
		} else {
			if sample > s.maximum {
				s.maximum = sample
			}

			if sample < s.minimum {
				s.minimum = sample
			}
		}
	}
}

// N returns the number of samples observed.
func (s *Statistics) N() uint64 {
	s.RLock()
	defer s.RUnlock()
	return s.samples
}

// Total returns the sum of the samples.
func (s *Statistics) Total() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.total
}

// Mean returns the average for all samples, computed as the sum of values
// divided by the total number of samples seen. If no samples have been added
// then this function returns 0.0. Note that 0.0 is a valid mean and does
// not necessarily mean that no samples have been tracked.
func (s *Statistics) Mean() float64 {
	s.RLock()
	defer s.RUnlock()
	if s.samples > 0 {
		return s.total / float64(s.samples)
	}
	return 0.0
}

// Variance computes the variability of samples and describes the distance of
// the distribution from the mean. If one or none samples have been added to
// the data set then this function returns 0.0 (two or more values are
// required to compute variance).
func (s *Statistics) Variance() float64 {
	s.RLock()
	defer s.RUnlock()

	n := float64(s.samples)
	if s.samples > 1 {
		num := (n*s.squares - s.total*s.total)
		den := (n * (n - 1))
		return num / den
	}

	return 0.0
}

// StdDev returns the standard deviation of samples, the square root of the
// variance. Two or more values are required to comput the standard deviation
// if one or none samples have been added to the data then this function
// returns 0.0.
func (s *Statistics) StdDev() float64 {
	s.RLock()
	defer s.RUnlock()

	if s.samples > 1 {
		return math.Sqrt(s.Variance())
	}
	return 0.0
}

// Maximum returns the maximum value of samples seen. If no samples have been
// added to the dataset, then this function returns 0.0.
func (s *Statistics) Maximum() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.maximum
}

// Minimum returns the minimum value of samples seen. If no samples have been
// added to the dataset, then this function returns 0.0.
func (s *Statistics) Minimum() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.minimum
}

// Range returns the difference between the maximum and minimum of samples.
// If no samples have been added to the dataset, this function returns 0.0.
// This function will also return zero if the maximum value equals the
// minimum value, e.g. in the case only one sample has been added or all of
// the samples are the same value.
func (s *Statistics) Range() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.maximum - s.minimum
}

// Serialize returns a map of summary statistics. This map is useful for
// dumping statistics to disk (using JSON for example) or for reporting the
// statistics elsewhere.
//
// TODO: Create Dump and Load functions to get statistical data to and from
// offline sources.
func (s *Statistics) Serialize() map[string]float64 {
	s.RLock()
	defer s.RUnlock()

	data := make(map[string]float64)
	data["samples"] = float64(s.samples)
	data["total"] = s.total
	data["mean"] = s.Mean()
	data["stddev"] = s.StdDev()
	data["variance"] = s.Variance()
	data["minimum"] = s.Minimum()
	data["maximum"] = s.Maximum()
	data["range"] = s.Range()
	return data
}

// Append another statistics object to the current statistics object,
// incrementing the distribution from the other object.
func (s *Statistics) Append(o *Statistics) {
	// Compute minimum and maximum aggregates by comparing both objects,
	// ensuring that zero valued items are not overriding the comparision.
	// Must come before any other aggregation.
	if o.samples > 0 {
		if o.maximum > s.maximum {
			s.maximum = o.maximum
		}

		if s.samples == 0 || o.minimum < s.minimum {
			s.minimum = o.minimum
		}
	}

	// Update the current statistics object
	s.total += o.total
	s.samples += o.samples
	s.squares += o.squares
}
