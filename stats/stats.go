package stats

import (
	"math"
	"sync"
)

// Statistics is an object that can keep track of descriptive statistics in
// an online fashion at runtime without saving each individual sample.
type Statistics struct {
	sync.RWMutex
	samples uint64       // number of samples seen
	total   float64      // the sum of all samples
	squares float64      // the sum of the squares of each sample
	maximum float64      // the maximum sample observed
	minimum float64      // the minimum sample observed
	async   bool         // whether or not to run in async mode
	done    chan bool    // the updater signals when it is done
	values  chan float64 // the channel to serialize values to
}

// Init the statistics; if async will create a buffered channel so that
// external callers can just dump values into it.
func (s *Statistics) Init(async bool) {
	s.async = async
	if s.async {
		s.done = make(chan bool)
		s.values = make(chan float64, 5000)
		go s.updater()
	}
}

// Close the updater to ensure that the values are all finalized.
func (s *Statistics) Close() {
	if s.async {
		close(s.values)
		<-s.done
	}
}

// Update the statistics with a single value (thread-safe)
func (s *Statistics) Update(sample float64) {
	if s.async {
		s.values <- sample
	} else {
		s.update(sample)
	}
}

// internal update method
func (s *Statistics) update(sample float64) {
	s.Lock()
	defer s.Unlock()

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

// updater just loops on the channel until it is closed, updating as it goes.
func (s *Statistics) updater() {
	for sample := range s.values {
		s.update(sample)
	}
	s.done <- true
}

// Mean returns the average for all samples.
func (s *Statistics) Mean() float64 {
	s.RLock()
	defer s.RUnlock()
	if s.samples > 0 {
		return s.total / float64(s.samples)
	}
	return 0.0
}

// Variance computes the variability of samples.
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

// StdDev returns the standard deviation of samples
func (s *Statistics) StdDev() float64 {
	s.RLock()
	defer s.RUnlock()

	if s.samples > 1 {
		return math.Sqrt(s.Variance())
	}
	return 0.0
}

// Maximum returns the maximum value of samples
func (s *Statistics) Maximum() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.maximum
}

// Minimum returns the minimum value of samples
func (s *Statistics) Minimum() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.minimum
}

// Range returns the difference between the maximum and minimum of samples
func (s *Statistics) Range() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.maximum - s.minimum
}

// Serialize returns a map of the stats that can be saved to disk
func (s *Statistics) Serialize() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	data := make(map[string]interface{})
	data["samples"] = s.samples
	data["total"] = s.total
	data["mean"] = s.Mean()
	data["stddev"] = s.StdDev()
	data["variance"] = s.Variance()
	data["minimum"] = s.minimum
	data["maximum"] = s.maximum
	return data
}
