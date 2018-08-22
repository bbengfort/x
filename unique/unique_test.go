package unique

import "fmt"

func ExampleStrings() {
	vals := []string{"foo", "bar", "foo", "baz", "bar", "zap", "foo"}
	deduped := Strings(vals)
	fmt.Println(deduped)
	// Output: [foo bar baz zap]
}

func ExampleInts() {
	vals := []int{-3, -2, -1, 1, 2, -3, 1, -2, 3, 1, 2, 3, 4, 5}
	deduped := Ints(vals)
	fmt.Println(deduped)
	// Output: [-3 -2 -1 1 2 3 4 5]
}

func ExampleInt32s() {
	vals := []int32{-3, -2, -1, 1, 2, -3, 1, -2, 3, 1, 2, 3, 4, 5}
	deduped := Int32s(vals)
	fmt.Println(deduped)
	// Output: [-3 -2 -1 1 2 3 4 5]
}

func ExampleInt64s() {
	vals := []int64{-3, -2, -1, 1, 2, -3, 1, -2, 3, 1, 2, 3, 4, 5}
	deduped := Int64s(vals)
	fmt.Println(deduped)
	// Output: [-3 -2 -1 1 2 3 4 5]
}

func ExampleUInts() {
	vals := []uint{1, 2, 3, 1, 2, 3, 1, 2, 3, 4, 5}
	deduped := UInts(vals)
	fmt.Println(deduped)
	// Output: [1 2 3 4 5]
}

func ExampleUInt32s() {
	vals := []uint32{1, 2, 3, 1, 2, 3, 1, 2, 3, 4, 5}
	deduped := UInt32s(vals)
	fmt.Println(deduped)
	// Output: [1 2 3 4 5]
}

func ExampleUInt64s() {
	vals := []uint64{1, 2, 3, 1, 2, 3, 1, 2, 3, 4, 5}
	deduped := UInt64s(vals)
	fmt.Println(deduped)
	// Output: [1 2 3 4 5]
}

func ExampleFloat32s() {
	vals := []float32{1.2, 2.2, 3.3, 1.2, 2.2, 3.1, 1.2, 2.2, 3.1, 3.3, 1.2}
	deduped := Float32s(vals)
	fmt.Println(deduped)
	// Output: [1.2 2.2 3.3 3.1]
}

func ExampleFloat64s() {
	vals := []float64{1.2, 2.2, 3.3, 1.2, 2.2, 3.1, 1.2, 2.2, 3.1, 3.3, 1.2}
	deduped := Float64s(vals)
	fmt.Println(deduped)
	// Output: [1.2 2.2 3.3 3.1]
}
