# Conflict-Free Replicated Versions (CFRV)

This package implements conflict-free replicated versions of multiple types. A conflict-free version number implemented on a distributed system means that two versions can be concurrently generated but still be ordered. Unlike auto-incrementing primary keys or sequences, this means that versions are structs that contain additional information.

The simplest type of CFRV (and the primary implementation of this package) I typically refer to as a Lamport Scalar. Each replica assigns scalar version numbers with a monotonically increasing counter and additionally include a unique process id. When a version is replicated locally, the counter is updated to the latest scalar. If two versions have the same scalar, then the process id is used to break the tie (e.g. the version with the lower process id is ordered before the later one).

Other types of CFRVs include vector and matrix clocks and more complex data structures (even things like TrueTime!). This package may implement these in the future, but primarily relies on vector component representations.

Rules for vectors:

1. All versions must be monotonically increasing
2. All versions must be unique and able to be ordered
3. Versions must point to their parent to show version sequences
