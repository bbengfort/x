// Package cfrv implements conflict-free replicated versions of multiple types.
// A conflict-free version number implemented on a distributed system means
// that two versions can be concurrently generated but still be ordered.
// Unlike auto-incrementing primary keys or sequences, this means that
// versions are structs that contain additional information.
//
// The simplest type of CFRV (and the primary implementation of this package)
// is typically referred to as a Lamport Scalar. Each replica assigns scalar
// version numbers with a monotonically increasing counter and additionally
// include a unique process id. When a version is replicated locally, the
// counter is updated to the latest scalar. If two versions have the same
// scalar, then the process id is used to break the tie (e.g. the version with
// the lower process id is ordered before the later one).
//
// Other types of CFRVs include vector and matrix clocks and more complex data
// structures (even things like TrueTime!). This package may implement these
// in the future, but primarily relies on vector component representations.
package cfrv

//===========================================================================
// CFRV and Factory Interfaces
//===========================================================================

// CFRV describes the behavior of a conflict-free replicated version data
// structure, particularly the methods for comparing two versions to each
// other to order the versions. All CFRVs must implement this interface to be
// used interchangeably in systems with varying requirements.
type CFRV interface {
	IsZero() bool                 // Returns true if the version is the zero value
	Equals(other CFRV) bool       // Returns true if version == other
	Greater(other CFRV) bool      // Returns true if version > other
	GreaterEqual(other CFRV) bool // Returns true if version >= other
	Lesser(other CFRV) bool       // Returns true if version < other
	LesserEqual(other CFRV) bool  // Returns true if version <= other
	String() string               // Returns a parseable string representation of the version
	Parse(s string) error         // Parses the string into the version struct
}

// Factory describes the behavior of a datastructure that generates conflict-
// free replicated versions. The factory maintains the global state to issue
// new versions, allowing version objects themselves to be stateless. Factories
// maintain versions for all objects in the system (even if there is just one)
// therefore they must identify the object by a unique name, the key. Factories
// that don't support keys can simply accept an empty string.
type Factory interface {
	Next() CFRV                   // return the next version for the given key
	Update(vers CFRV)             // update the state of the factory with the given version
	Parse(s string) (CFRV, error) // parse the version from a string and return
}
