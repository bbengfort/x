package stats

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func loadBenchData() ([]time.Duration, error) {
	data := make([]time.Duration, 0, 1000000)

	f, err := os.Open("testdata/latencies.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	for buf.Scan() {
		var val time.Duration
		if val, err = time.ParseDuration(buf.Text()); err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	if buf.Err() != nil {
		return nil, err
	}

	return data, nil
}

func ExampleBenchmark() {
	stats := new(Benchmark)
	samples, _ := loadBenchData()

	for _, sample := range samples {
		stats.Update(sample)
	}

	data, _ := json.MarshalIndent(stats.Serialize(), "", "  ")
	fmt.Println(string(data))
	// Output:
	// {
	//   "fastest": "41.219436ms",
	//   "mean": "120.993689ms",
	//   "range": "167.175236ms",
	//   "samples": 1000000,
	//   "slowest": "208.394672ms",
	//   "stddev": "17.283562ms",
	//   "throughput": 8.264893850648656,
	//   "total": "33h36m33.689461785s",
	//   "variance": "298.721µs"
	// }
}

func TestBenchmark(t *testing.T) {
	RegisterTestingT(t)

	data, err := loadBenchData()
	Ω(err).ShouldNot(HaveOccurred())

	stats := new(Benchmark)

	for _, v := range data {
		stats.Update(v)
	}

	Ω(stats.Mean()).Should(Equal(time.Duration(120993689)))
	Ω(stats.StdDev()).Should(Equal(time.Duration(17283562)))
	Ω(stats.Variance()).Should(Equal(time.Duration(298721)))
	Ω(stats.Slowest()).Should(Equal(time.Duration(208394672)))
	Ω(stats.Fastest()).Should(Equal(time.Duration(41219436)))
	Ω(stats.Range()).Should(Equal(time.Duration(167175236)))
}

func TestBenchmarkBulk(t *testing.T) {
	RegisterTestingT(t)

	data, err := loadBenchData()
	Ω(err).ShouldNot(HaveOccurred())

	stats := new(Benchmark)
	stats.Update(data...)

	Ω(stats.Mean()).Should(Equal(time.Duration(120993689)))
	Ω(stats.StdDev()).Should(Equal(time.Duration(17283562)))
	Ω(stats.Variance()).Should(Equal(time.Duration(298721)))
	Ω(stats.Slowest()).Should(Equal(time.Duration(208394672)))
	Ω(stats.Fastest()).Should(Equal(time.Duration(41219436)))
	Ω(stats.Range()).Should(Equal(time.Duration(167175236)))
}

func BenchmarkBenchmark_Update(b *testing.B) {
	rand.Seed(42)
	stats := new(Benchmark)

	for i := 0; i < b.N; i++ {
		val := time.Duration(rand.Int31n(1000)) * time.Millisecond
		stats.Update(val)
	}
}

func BenchmarkBenchmark_Sequential(b *testing.B) {
	data, _ := loadBenchData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats := new(Benchmark)
		for _, val := range data {
			stats.Update(val)
		}
	}
}

func BenchmarkBenchmark_BulkLoad(b *testing.B) {
	data, _ := loadBenchData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stats := new(Benchmark)
		stats.Update(data...)
	}
}
