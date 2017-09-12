package stats

import (
	"sync"
	"time"
)

// Benchmark keeps track of a distrubtion of durations, e.g. to benchmark the
// performance or timing of an operation. It returns descriptive statistics
// as durations so that they can be read as timings. Benchmark works in an
// online fashion similar to the Statistics object, but works on
// time.Duration samples instead of floats. Instead of minimum and maximum
// values it returns the fastest and slowest times.
//
// The primary entry point to the object is via the Update method, where one
// or more time.Durations can be passed. This object has unexported fields
// because it is thread-safe (via a sync.RWMutex). All properties must be
// accesesd from read-locked access methods.
type Benchmark struct {
	sync.RWMutex
	Statistics
	timeouts uint64 // the number of 0 durations (null durations) or timeouts
}

// Update the benchmark with a duration or durations (thread-safe). If a
// duration of 0 is passed, then it is interpreted as a timeout -- e.g. a
// maximal duration bound had been reached. Timeouts are recorded in a
// separate counter and can be used to express failure measures.
func (s *Benchmark) Update(durations ...time.Duration) {
	s.Lock()
	defer s.Unlock()

	for _, duration := range durations {
		// Record any timeouts in the benchmark
		if duration == 0 {
			s.timeouts++
			continue
		}

		s.Statistics.Update(duration.Seconds())
	}
}

// Throughput returns the number of samples per second, measured as the
// inverse mean: number of samples divided by the total duration in seconds.
// This metric does not express a duration, so a float64 value is returned
// instead. If the duration or number of accesses is zero, 0.0 is returned.
func (s *Benchmark) Throughput() float64 {
	s.RLock()
	defer s.RUnlock()
	if s.samples > 0 && s.total > 0 {
		return float64(s.Statistics.samples) / s.Statistics.total
	}
	return 0.0
}

// Total returns the total duration recorded across all samples.
func (s *Benchmark) Total() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.total)
}

// Mean returns the average for all durations expressed as float64 seconds
// and returns a time.Duration which is expressed in int64 nanoseconds. This
// can mean some loss in precision of the mean value, but also allows the
// caller to compute the mean in varying timescales. Since microseconds is a
// pretty fine granularity for timings, truncating the floating point of the
// nanosecond seems acceptable.
//
// If no durations have been recorded, a zero valued duration is returned.
func (s *Benchmark) Mean() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.Mean())
}

// Variance computes the variability of samples and describes the distance of
// the distribution from the mean. This function returns a time.Duration,
// which can mean a loss in precision lower than the microsecond level. This
// is usually acceptable for most applications.
//
// If no more than 1 durations were recorded, returns a zero valued duration.
func (s *Benchmark) Variance() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.Variance())
}

// StdDev returns the standard deviation of samples, the square root of the
// variance. This function returns a time.Duration which represents a loss in
// precision from int64 nanoseconds to float64 seconds.
//
// If no more than 1 durations were recorded, returns a zero valued duration.
func (s *Benchmark) StdDev() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.StdDev())
}

// Slowest returns the maximum value of durations seen. If no durations have
// been added to the dataset, then this function returns a zero duration.
func (s *Benchmark) Slowest() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.Maximum())
}

// Fastest returns the minimum value of durations seen. If no durations have
// been added to the dataset, then this function returns a zero duration.
func (s *Benchmark) Fastest() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.Minimum())
}

// Range returns the difference between the slowest and fastest durations.
// If no samples have been added to the dataset, this function returns a zero
// duration. It will also return zero if the fastest and slowest durations
// are equal. E.g. in the case only one duration has been recorded or such
// that all durations have the same value.
func (s *Benchmark) Range() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.castSeconds(s.Statistics.Range())
}

// Serialize returns a map of summary statistics. This map is useful for
// dumping statistics to disk (using JSON for example) or for reporting the
// statistics elsewhere. The values in the maps are string representations of
// the time.Duration objects, which are reported in a human readable form.
// They can be converted back to durations with time.ParseDuration.
//
// TODO: Create Dump and Load functions to get statistical data to and from
// offline sources.
func (s *Benchmark) Serialize() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	data := make(map[string]interface{})
	data["samples"] = s.samples
	data["total"] = s.Total().String()
	data["mean"] = s.Mean().String()
	data["stddev"] = s.StdDev().String()
	data["variance"] = s.Variance().String()
	data["fastest"] = s.Fastest().String()
	data["slowest"] = s.Slowest().String()
	data["range"] = s.Range().String()
	data["throughput"] = s.Throughput()
	return data
}

// Internal Helper Method to cast float64 seconds into a duration
func (s *Benchmark) castSeconds(seconds float64) time.Duration {
	return time.Duration(float64(time.Second) * seconds)
}
