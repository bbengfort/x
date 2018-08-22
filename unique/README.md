# Unique

**Finds unique elements in a slice**.

This is a helper package for common functionality in my Go code. See [Creating Unique Slices in Go](https://kylewbanks.com/blog/creating-unique-slices-in-go) for a blog post about the methodology being used here.

Basic usage:

```go
import "github.com/bbengfort/x/unique"

names := []string{"foo", "bar", "foo", "baz", "zap", "bar"}
deduped := unique.Strings(names)
// []string{"bar", "baz", "foo", "zap"}
```

There is a unique function for all the types that I've needed; sorry if there isn't one for the type you were hoping for.
