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

## Blocking vs. Non-Blocking

I received a surprising result when I tried to implement a non-blocking version of the `Statistics` struct by using a buffered-channel. A write-up of that can be found here: [Online Distribution](https://bbengfort.github.io/snippets/2017/08/28/online-distribution.html).

The benchmarks are as follows:

```
BenchmarkBlocking-8      	20000000            81.1 ns/op
BenchmarkNonBlocking-8   	10000000	       140 ns/op
```

As such, the current implementation simply uses thread-safe locks rather than a channel.
