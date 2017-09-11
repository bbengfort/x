/*
Package bandit implements multi-armed bandit strategies for random choice.
*/
package bandit

import (
	"math"
	"math/rand"
)

// Strategy specifies the methods required by an algorithm to compute
// multi-armed bandit probabilities for reinforcement learning. The basic
// mechanism allows you to initialize a strategy with n arms (or n choices).
// The Select() method will return a selected index based on the internal
// strategy, and the Update() method allows external callers to update the
// reward function for the selected arm.
type Strategy interface {
	Init(nArms int)         // Initialize the bandit with n choices
	Select() int            // Selects an arm and returns the index of the choice
	Update(arm, reward int) // Update the given arm with a reward
	Counts() []uint64       // The frequency of each arm being selected
	Values() []float64      // The reward distributions for each arm
	Serialize() interface{} // Return a JSON representation of the strategy
}

//===========================================================================
// Epsilon Greedy Multi-Armed Bandit
//===========================================================================

// EpsilonGreedy implements a reinforcement learning strategy such that the
// maximizing value is selected with probability epsilon and a uniform random
// selection is made with probability 1-epsilon.
type EpsilonGreedy struct {
	Epsilon float64   // Probability of selecting maximizing value
	counts  []uint64  // Number of times each index was selected
	values  []float64 // Reward values condition by frequency
}

// Init the bandit with nArms number of possible choices, which are referred
// to by index in both the Counts and Values arrays.
func (b *EpsilonGreedy) Init(nArms int) {
	b.counts = make([]uint64, nArms, nArms)
	b.values = make([]float64, nArms, nArms)
}

// Select the arm with the maximizing value with probability epsilon,
// otherwise uniform random selection of all arms with probability 1-epsilon.
func (b *EpsilonGreedy) Select() int {
	if rand.Float64() > b.Epsilon {
		// Select the maximal value from values.
		max := -1.0
		idx := -1

		// Find the index of the maximal value.
		for i, val := range b.values {
			if val > max {
				max = val
				idx = i
			}
		}

		return idx
	}

	// Otherwise return any of the values
	return rand.Intn(len(b.values))
}

// Update the selected arm with the reward so that the strategy can learn the
// maximizing value (conditioned by the frequency of selection).
func (b *EpsilonGreedy) Update(arm, reward int) {
	// Update the frequency
	b.counts[arm]++
	n := float64(b.counts[arm])

	value := b.values[arm]
	b.values[arm] = ((n-1)/n)*value + (1/n)*float64(reward)
}

// Counts returns the frequency each arm was selected
func (b *EpsilonGreedy) Counts() []uint64 {
	return b.counts
}

// Values returns the reward distribution of each arm
func (b *EpsilonGreedy) Values() []float64 {
	return b.values
}

// Serialize the bandit strategy to dump to JSON.
func (b *EpsilonGreedy) Serialize() interface{} {
	data := make(map[string]interface{})
	data["strategy"] = "epsilon greedy"
	data["epsilon"] = b.Epsilon
	data["counts"] = b.counts
	data["values"] = b.values
	return data
}

//===========================================================================
// Annealing Epsilon Greedy Multi-Armed Bandit
//===========================================================================

// AnnealingEpsilonGreedy implements a reinforcement learning strategy such
// that value of epsilon starts small then grows increasingly bigger, leading
// to an exploring learning strategy at start and prefering exploitation as
// more selections are made.
type AnnealingEpsilonGreedy struct {
	counts []uint64  // Number of times each index was selected
	values []float64 // Reward values condition by frequency
}

// Init the bandit with nArms number of possible choices, which are referred
// to by index in both the Counts and Values arrays.
func (b *AnnealingEpsilonGreedy) Init(nArms int) {
	b.counts = make([]uint64, nArms, nArms)
	b.values = make([]float64, nArms, nArms)
}

// Epsilon is computed by the current number of trials such that the more
// trials have occured, the smaller epsilon is (on a log scale).
func (b *AnnealingEpsilonGreedy) Epsilon() float64 {
	// Compute epsilon based on the total number of trials
	t := uint64(1)
	for _, i := range b.counts {
		t += i
	}

	// The more trials the smaller that epsilon is
	return 1 / math.Log(float64(t)+0.0000001)
}

// Select the arm with the maximizing value with probability epsilon,
// otherwise uniform random selection of all arms with probability 1-epsilon.
func (b *AnnealingEpsilonGreedy) Select() int {
	if rand.Float64() > b.Epsilon() {
		// Select the maximal value from values.
		max := -1.0
		idx := -1

		// Find the index of the maximal value.
		for i, val := range b.values {
			if val > max {
				max = val
				idx = i
			}
		}

		return idx
	}

	// Otherwise return any of the values
	return rand.Intn(len(b.values))
}

// Update the selected arm with the reward so that the strategy can learn the
// maximizing value (conditioned by the frequency of selection).
func (b *AnnealingEpsilonGreedy) Update(arm, reward int) {
	// Update the frequency
	b.counts[arm]++
	n := float64(b.counts[arm])

	value := b.values[arm]
	b.values[arm] = ((n-1)/n)*value + (1/n)*float64(reward)
}

// Counts returns the frequency each arm was selected
func (b *AnnealingEpsilonGreedy) Counts() []uint64 {
	return b.counts
}

// Values returns the reward distribution of each arm
func (b *AnnealingEpsilonGreedy) Values() []float64 {
	return b.values
}

// Serialize the bandit strategy to dump to JSON.
func (b *AnnealingEpsilonGreedy) Serialize() interface{} {
	data := make(map[string]interface{})
	data["strategy"] = "annealing epsilon greedy"
	data["epsilon"] = b.Epsilon()
	data["counts"] = b.counts
	data["values"] = b.values
	return data
}

//===========================================================================
// Uniform strategy
//===========================================================================

// Uniform selects all values with an equal likelihood on every selection.
// While it tracks the frequency of selection and the reward costs, this
// information does not affect the way it selects values.
type Uniform struct {
	counts []uint64  // Number of times each index was selected
	values []float64 // Reward values condition by frequency
}

// Init the bandit with nArms number of possible choices, which are referred
// to by index in both the Counts and Values arrays.
func (b *Uniform) Init(nArms int) {
	b.counts = make([]uint64, nArms, nArms)
	b.values = make([]float64, nArms, nArms)
}

// Select the arm with equal probability for each choice.
func (b *Uniform) Select() int {
	return rand.Intn(len(b.values))
}

// Update the selected arm with the reward so that the strategy can learn the
// maximizing value (conditioned by the frequency of selection).
func (b *Uniform) Update(arm, reward int) {
	// Update the frequency
	b.counts[arm]++
	n := float64(b.counts[arm])

	value := b.values[arm]
	b.values[arm] = ((n-1)/n)*value + (1/n)*float64(reward)
}

// Counts returns the frequency each arm was selected
func (b *Uniform) Counts() []uint64 {
	return b.counts
}

// Values returns the reward distribution of each arm
func (b *Uniform) Values() []float64 {
	return b.values
}

// Serialize the bandit strategy to dump to JSON.
func (b *Uniform) Serialize() interface{} {
	data := make(map[string]interface{})
	data["strategy"] = "uniform selection"
	data["counts"] = b.counts
	data["values"] = b.values
	return data
}
