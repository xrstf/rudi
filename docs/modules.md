# Rudi Extended Library

To keep the dependencies of Rudi itself to a minimum, functionality that requires external libraries
is implemented in standalone Go modules. This makes it possible for integrators to choose exactly
what kind of functions they need and which dependencies they can bear.

All of the following modules are kept in the [rudi-contrib](https://github.com/xrstf/rudi-contrib)
repository on GitHub.

The extended library also serves as a great tutorial on how to wrap existing code in Rudi :smile:

Note that the `rudi` interpreter (the binary) has all of these modules built-in, as the CLI
interpreter is its own Go module and does not contribute to the Rudi language repository.

## semver

The [semver module](https://github.com/xrstf/rudi-contrib/tree/main/semver) integrates the
[blang/semver](https://github.com/blang/semver) library and allows to parse and compare semantic
versions.

* `(semver "9.2.0")`

## uuid

The [uuid module](https://github.com/xrstf/rudi-contrib/tree/main/uuid) contains functionality to
create UUIDs.

* `(uuid)`

## yaml

The [yaml module](https://github.com/xrstf/rudi-contrib/tree/main/yaml) contains integrates
[gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) and allows encoding/decoding YAML.

* `(from-yaml "{foo: 3}")`
* `(to-yaml {foo 3})`

