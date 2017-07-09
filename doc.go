// Package x hosts several packages, modules, and libraries that are common
// across most of my code for easy reuse. This package is very much in the
// spirit of golang.org/x though it does have a slightly longer import path as
// a result of being hosted in my GitHub repository.
//
// One thing that I think is important to note is that most of the subpackages
// in this repository are independent. That is that they are implemented and
// tested separately from other subpackages. Future me and anyone who would
// like to use this package should only go get exactly what they need and
// rely on the documentation on godoc and in the subpackage README.md for more
// information.
//
// Generally speaking, these things are simply ported out of my other
// applications once I discover that they need to be reused. The x repository
// gives me the ability to manage them all in the same version control without
// jumping through all the GitHub hoops. I'm not sure this is what was
// intended by Golang, but managing multiple repositories with just one or two
// files was too much of a pain, hence this system.
package x
