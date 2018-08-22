/*
Package unique implements type-specific functions to deduplicate a slice and
return a slice with only the unique elements in it. It uses an intermediate
map to store values already seen in the slice, then returns the new slice.
Note that no guarantees are made about ordering.
*/
package unique

// Strings returns a slice of strings without any duplicates in it.
func Strings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// Ints returns a slice of ints without any duplicates in it.
func Ints(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// Int32s returns a slice of int32 without any duplicates in it.
func Int32s(input []int32) []int32 {
	u := make([]int32, 0, len(input))
	m := make(map[int32]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// Int64s returns a slice of int64 without any duplicates in it.
func Int64s(input []int64) []int64 {
	u := make([]int64, 0, len(input))
	m := make(map[int64]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// UInts returns a slice of uints without any duplicates in it.
func UInts(input []uint) []uint {
	u := make([]uint, 0, len(input))
	m := make(map[uint]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// UInt32s returns a slice of uint32 without any duplicates in it.
func UInt32s(input []uint32) []uint32 {
	u := make([]uint32, 0, len(input))
	m := make(map[uint32]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// UInt64s returns a slice of uint64 without any duplicates in it.
func UInt64s(input []uint64) []uint64 {
	u := make([]uint64, 0, len(input))
	m := make(map[uint64]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// Float32s returns a slice of float32 without any duplicates in it.
func Float32s(input []float32) []float32 {
	u := make([]float32, 0, len(input))
	m := make(map[float32]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// Float64s returns a slice of float64 without any duplicates in it.
func Float64s(input []float64) []float64 {
	u := make([]float64, 0, len(input))
	m := make(map[float64]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
