# github.com/bbengfort/x

**Go packages that are common to many of my projects -- in the spirit of golang.org/x**

## Usage

To get these packages on your system, it's as easy as:

```
$ go get github.com/bbengfort/x/[pkg]
```

Where `[pkg]` is the name of the package you want to install. Import the packages required in your code, vendoring as necessary in order to use them in other projects.

## Tests and Benchmarks

Packages in this repository use different tests and benchmarking tools. Beyond the standard `go test` environment, most packages will probably use [ginkgo](http://onsi.github.io/ginkgo/) for BDD style testing and [gomega](http://onsi.github.io/gomega/) for matching and assertions. To install these packages:

```
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega
```

This will install the libraries as well as the runner executables. 
