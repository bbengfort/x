# Console

**Package console implements simple hierarchical logging functionality.**

This is a pretty standard hierarchical console logging module that allows you to use tiered functions to manage what is printed to the command line. For example if the level is set to `Info`, then `Trace` and `Debug` messages will be implemented as no-ops, reducing the amount of information printed.

Setting up the console usually happens in the `init()` method of a package:

```go
func init() {
	// Initialize our debug logging with our prefix
	console.Init("[myapp] ", log.Lmicroseconds)
    console.SetLevel(console.LevelInfo)
}
```

Now the logging functions can automatically be used:

```go
console.Trace("routine %s happening", thing)
console.Debug("sending message #%d from %s to %s", msg, send, recv)
console.Info("listening on %s", addr)
console.Status("completed %d out of %d tasks", completed, nTasks)
console.Warn("limit of %d queries reached", nQueries)
console.Warne("something bad happened: %s", err)
```

The purpose of these functions were to have simple pout and perr methods inside of applications. Another way to use this library is simply to copy and paste this code and lowercase the function names into your app. 
