package stats

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
)

func loadTestData() ([]float64, error) {
	data := make([]float64, 0, 1000000)

	f, err := os.Open("testdata/standardized.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	for buf.Scan() {
		var val float64
		if val, err = strconv.ParseFloat(buf.Text(), 64); err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	if buf.Err() != nil {
		return nil, err
	}

	return data, nil
}

func ExampleStatistics() {
	stats := new(Statistics)
	samples, _ := loadTestData()

	for _, sample := range samples {
		stats.Update(sample)
	}

	data, _ := json.MarshalIndent(stats.Serialize(), "", "  ")
	fmt.Println(string(data))
	// Output:
	// {
	//   "maximum": -4.72206033824,
	//   "mean": 0.00041124313405184064,
	//   "minimum": 5.30507026071,
	//   "range": 10.02713059895,
	//   "samples": 1000000,
	//   "stddev": 0.9988808397330513,
	//   "total": 411.2431340518406,
	//   "variance": 0.9977629319858057
	// }
}

func TestStatistics(t *testing.T) {
	RegisterTestingT(t)

	data, err := loadTestData()
	Ω(err).ShouldNot(HaveOccurred())

	stats := new(Statistics)

	Ω(stats.values).Should(BeNil())

	for _, v := range data {
		stats.Update(v)
	}

	Ω(stats.Mean()).Should(Equal(0.00041124313405184064))
	Ω(stats.StdDev()).Should(Equal(0.9988808397330513))
	Ω(stats.Variance()).Should(Equal(0.9977629319858057))
	Ω(stats.Maximum()).Should(Equal(5.30507026071))
	Ω(stats.Minimum()).Should(Equal(-4.7220603382400004))
	Ω(stats.Range()).Should(Equal(10.02713059895))
}

func BenchmarkStatistics_Update(b *testing.B) {
	rand.Seed(42)
	stats := new(Statistics)

	for i := 0; i < b.N; i++ {
		val := rand.Float64()
		stats.Update(val)
	}
}
