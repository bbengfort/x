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

Basically, as samples come in, you can pass them to the `Update` method collecting summary statistics as you go.
