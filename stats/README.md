# Stats
**Online computation of descriptive statistics**

This package provides a struct for computing online descriptive statistics without saving all values in an array or to disk. Install the package:

```
$ go get github.com/bbengfort/x/stats
```

Usage as follows:

```go
stats := new(stats.Statistics)

for i := 0; i < 1000; i++ {
    stats.Update(rand.Float64())
}

mu := stats.Mean()
sigma := stats.StdDev()
```

Basically, as samples come in, you can pass them to the `Update` method collecting summary statistics as you go. You can then dump the code out into a JSON dictionary as follows:

```go
data := stats.Serialize()
```

**NOTE:** The `Statistics` object _is thread-safe_ by virtue of a `sync.RWMutex` that locks and unlocks the data structure on every call.

## Bulk Loading

It is possible to bulk-load the statistics object by passing multiple float64 values using variadic arguments:

```go
stats.Update(1.2, 3.1, 4.2, 1.2)
```

Or by passing in an array of float64 values:

```go
var data []float64
stats.Update(data...)
```

This is much faster than loading values individually in a for loop as demonstrated by the following benchmarks:

```
BenchmarkStatistics_Update-8       	20000000	        79.2 ns/op
BenchmarkStatistics_Sequential-8   	      30	  56017960 ns/op
BenchmarkStatistics_BulkLoad-8     	     500	   2514950 ns/op
```

The first benchmark, `BenchmarkStatistics_Update-8`, is the time it takes to
update a single value into the statistics. The second benchmark, `BenchmarkStatistics_Sequential-8` uses a for-loop to Update one value at a time from 1,000,000 values. The third benchmark, `BenchmarkStatistics_BulkLoad-8`, simply passes the entire array of 1M values directly to the function and as a result is 22x faster.

## Blocking vs. Non-Blocking

I received a surprising result when I tried to implement a non-blocking version of the `Statistics` struct by using a buffered-channel. A write-up of that can be found here: [Online Distribution](https://bbengfort.github.io/snippets/2017/08/28/online-distribution.html).

The benchmarks are as follows:

```
BenchmarkBlocking-8      	20000000            81.1 ns/op
BenchmarkNonBlocking-8   	10000000	       140 ns/op
```

As such, the current implementation simply uses thread-safe locks rather than a channel.

## Benchmarks

This package also includes a specialized data structure for computing online statistical distribution of `time.Duration` objects called the `Benchmark`. Similar to the `Statistics` object you can update it, but with `time.Duration` objects, which are then converted into `float64` seconds values using the `time.Duration.Seconds()` method. This reduces the granularity from `int64` nanoseconds, but should still be good at about the microsecond granularity.

The reason for the conversion is because computation with `int64` quickly overflows especially when computing the sum of squares. By converting to a float, the domain of the online distribution is similar to the domain of the `Statistics` object.
