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

### uuid

* [`uuidv4`](../extlib/uuid/uuidv4.md) – returns a new, randomly generated v4 UUID

### yaml

* [`from-yaml`](../extlib/yaml/from-yaml.md) – decodes a YAML string into a Go value
* [`to-yaml`](../extlib/yaml/to-yaml.md) – encodes the given value as YAML
<!-- END_EXTLIB_TOC -->
