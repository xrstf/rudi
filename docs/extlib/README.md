# Rudi Extended Library

To keep the dependencies of Rudi itself to a minimum, functionality that requires external libraries
is implemented in standalone Go modules. This makes it possible for integrators to choose exactly
what kind of functions they need and which dependencies they can bear.

All of the following modules are kept in the [rudi-contrib](https://github.com/xrstf/rudi-contrib)
repository on GitHub.

The extended library also serves as a great tutorial on how to wrap existing code in Rudi :smile:

Note that the `rudi` interpreter (the binary) has all of these modules built-in, as the CLI
interpreter is its own Go module and does not contribute to the Rudi language repository.

<!-- BEGIN_EXTLIB_TOC -->
### semver

* [`semver`](../extlib/semver/semver.md) – parses a string as a semantic version

### set

* [`new-key-set`](../extlib/set/new-key-set.md) – create a set filled with the keys of an object
* [`new-set`](../extlib/set/new-set.md) – create a set filled with the given values
* [`set-delete`](../extlib/set/set-delete.md) – returns a copy of the set with the given values removed from it
* [`set-diff`](../extlib/set/set-diff.md) – returns the difference between two sets
* [`set-eq?`](../extlib/set/set-eq.md) – returns true if two sets hold the same values
* [`set-has-any?`](../extlib/set/set-has-any.md) – returns true if the set contains _any_ of the given values
* [`set-has?`](../extlib/set/set-has.md) – returns true if the set contains _all_ of the given values
* [`set-insert`](../extlib/set/set-insert.md) – returns a copy of the set with the newly added values inserted to it
* [`set-intersection`](../extlib/set/set-intersection.md) – returns the insersection of two sets
* [`set-list`](../extlib/set/set-list.md) – returns a sorted vector containing the values of the set
* [`set-size`](../extlib/set/set-size.md) – returns the number of values in the set
* [`set-superset-of?`](../extlib/set/set-superset-of.md) – returns true if the other set is a superset of the base set
* [`set-symdiff`](../extlib/set/set-symdiff.md) – returns the symmetric difference between two sets
* [`set-union`](../extlib/set/set-union.md) – returns the union of two or more sets

### uuid

* [`uuidv4`](../extlib/uuid/uuidv4.md) – returns a new, randomly generated v4 UUID

### yaml

* [`from-yaml`](../extlib/yaml/from-yaml.md) – decodes a YAML string into a Go value
* [`to-yaml`](../extlib/yaml/to-yaml.md) – encodes the given value as YAML
<!-- END_EXTLIB_TOC -->
