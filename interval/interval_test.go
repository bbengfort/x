package interval_test

import (
	"math/rand"
	"time"

	"github.com/bbengfort/x/events"
	. "github.com/bbengfort/x/interval"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interval", func() {

	var calls int64
	var echan chan<- error
	var counter events.Callback
	var started time.Time
	var since time.Duration

	delay := 5 * time.Millisecond
	wait := 24 * time.Millisecond

	BeforeEach(func() {
		calls = 0
		since = 0
		echan = make(chan error, 10)
		counter = func(e events.Event) error {
			calls++
			since += time.Since(started)
			started = time.Now()
			return nil
		}
	})

	Describe("Fixed Interval", func() {

		It("should not start the interval on init", func() {
			ticker := new(FixedInterval)
			ticker.Init(delay, events.TimeoutEvent, echan)
			ticker.Register(counter)
			started = time.Now()

			time.Sleep(wait)
			Ω(calls).Should(BeZero())
		})

		It("should not start an uninitialized interval", func() {
			ticker := new(FixedInterval)
			Ω(ticker.Start()).Should(BeFalse())
		})

		It("should emit an event at a fixed interval", func() {
			ticker := new(FixedInterval)
			ticker.Init(delay, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Start the ticker
			started = time.Now()
			Ω(ticker.Start()).Should(BeTrue())
			time.Sleep(wait)
			ticker.Stop()

			Ω(calls).Should(BeNumerically("==", 4))
			Ω(since).Should(BeNumerically(">=", 4*time.Millisecond))
			Ω(since).Should(BeNumerically("<", wait+delay))
		})

		It("should be able to determine if the interval is running", func() {
			ticker := new(FixedInterval)
			ticker.Init(delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Should not be running to start
			Ω(ticker.Running()).Should(BeFalse())

			// Start and stop a bunch of times
			for i := 0; i < 100; i++ {
				Ω(ticker.Start()).Should(BeTrue())
				Ω(ticker.Running()).Should(BeTrue())
				Ω(ticker.Stop()).Should(BeTrue())
				Ω(ticker.Running()).Should(BeFalse())
			}

			// No calls should have executed
			Ω(calls).Should(BeZero())
		})

		It("should be able to interrupt an interval", func() {
			ticker := new(FixedInterval)
			ticker.Init(delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Start the ticker
			started = time.Now()
			Ω(ticker.Start()).Should(BeTrue())
			time.Sleep(5 * time.Millisecond)

			// Interrupt the ticker
			Ω(ticker.Interrupt()).Should(BeTrue())

			time.Sleep(7 * time.Millisecond)
			Ω(calls).Should(Equal(int64(0)))

			// Wait until ticker is timed out.
			time.Sleep(18 * time.Millisecond)
			Ω(ticker.Stop()).Should(BeTrue())
			Ω(calls).Should(Equal(int64(2)))
		})

		It("should be able to stop an interval", func() {
			ticker := new(FixedInterval)
			ticker.Init(delay, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Start the ticker
			started = time.Now()
			Ω(ticker.Start()).Should(BeTrue())
			time.Sleep(7 * time.Millisecond)
			ticker.Stop()

			time.Sleep(22 * time.Millisecond)
			Ω(calls).Should(BeNumerically("==", 1))
		})

	})

	Describe("Random Interval", func() {

		BeforeEach(func() {
			// Set the random seed to produce expected behavior
			rand.Seed(42)
		})

		It("should not start the interval on init", func() {
			ticker := new(RandomInterval)
			ticker.Init(delay, delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)
			started = time.Now()

			time.Sleep(wait)
			Ω(calls).Should(BeZero())
		})

		It("should not start an uninitialized interval", func() {
			ticker := new(RandomInterval)
			Ω(ticker.Start()).Should(BeFalse())
		})

		It("should return a different delay on GetDelay", func() {
			ticker := new(RandomInterval)
			ticker.Init(delay, delay*2, events.TimeoutEvent, echan)

			delays := make([]time.Duration, 0, 100)
			for i := 0; i < 1000; i++ {
				delays = append(delays, ticker.GetDelay())
				Ω(delays[i]).Should(BeNumerically("<", delay*2))
				Ω(delays[i]).Should(BeNumerically(">", delay))

				if i > 0 {
					Ω(delays[i]).ShouldNot(Equal(delays[i-1]))
				}
			}
		})

		It("should emit an event at a random interval", func() {
			ticker := new(RandomInterval)
			ticker.Init(delay, delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Start the ticker
			started = time.Now()
			Ω(ticker.Start()).Should(BeTrue())
			time.Sleep(wait)
			ticker.Stop()

			Ω(calls).Should(BeNumerically(">=", 2))
			Ω(calls).Should(BeNumerically("<=", 4))
			Ω(since).Should(BeNumerically(">=", 4*time.Millisecond))
			Ω(since).Should(BeNumerically("<", wait+delay))
		})

		It("should be able to determine if the interval is running", func() {
			ticker := new(RandomInterval)
			ticker.Init(delay, delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Should not be running to start
			Ω(ticker.Running()).Should(BeFalse())

			// Start and stop a bunch of times
			for i := 0; i < 100; i++ {
				Ω(ticker.Start()).Should(BeTrue())
				Ω(ticker.Running()).Should(BeTrue())
				Ω(ticker.Stop()).Should(BeTrue())
				Ω(ticker.Running()).Should(BeFalse())
			}

			// No calls should have executed
			Ω(calls).Should(BeZero())
		})

		It("should be able to interrupt an interval", func() {
			ticker := new(RandomInterval)
			ticker.Init(delay, delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Start the ticker
			started = time.Now()
			Ω(ticker.Start()).Should(BeTrue())
			time.Sleep(delay - (1 * time.Millisecond))

			// Interrupt the ticker
			Ω(ticker.Interrupt()).Should(BeTrue())

			time.Sleep(delay - (1 * time.Millisecond))
			Ω(calls).Should(Equal(int64(0)))

			// Wait until ticker is timed out.
			time.Sleep(wait)
			Ω(ticker.Stop()).Should(BeTrue())
			Ω(calls).Should(BeNumerically("<=", 4))
		})

		It("should be able to stop an interval", func() {
			ticker := new(RandomInterval)
			ticker.Init(delay, delay*2, events.TimeoutEvent, echan)
			ticker.Register(counter)

			// Start the ticker
			started = time.Now()
			Ω(ticker.Start()).Should(BeTrue())
			time.Sleep(delay * 2)
			ticker.Stop()

			time.Sleep(22 * time.Millisecond)
			Ω(calls).Should(BeNumerically("==", 1))
		})
	})

})
