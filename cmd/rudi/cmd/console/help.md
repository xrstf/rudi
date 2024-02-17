# Welcome to the Rudi interpreter :)

You can enter one of

* A path expression, like `.foo` or `.foo[0].bar` to access the global document.
* An expression like (+ .foo 42) to compute data by functions; see the topics
  below or the Rudi website for a complete list of available functions.
* A scalar JSON value, like `3` or `[1 2 3]`, which will simply return that
  exact value with no further side effects. Not super useful usually.

## Commands

Additionally, the following commands can be used:

* help       – Show this help text.
* help TOPIC – Show help for a specific topic.
* exit       – Exit Rudi immediately.

## Help Topics

The following topics are available and can be accessed using `help TOPIC`:

<!-- BEGIN_HELP_TOPICS_TOC -->
* `language` – A short introduction to the Rudi language
* `coalescing` – How Rudi handles, converts and compares values
<!-- END_HELP_TOPICS_TOC -->

You can also request help for any of these functions using `help FUNCTION`:

<!-- BEGIN_HELP_LIB_TOC -->
* **core**
  * `case` – chooses the first expression for which the test is true
  * `default` – returns the default value if the first argument is empty
  * `delete` – removes a key from an object or an item from a vector
  * `do` – eval a sequence of statements where only one expression is valid
  * `empty?` – returns true when the given value is empty-ish (0, false, null, "", ...)
  * `error` – returns an error
  * `has?` – returns true if the given symbol's path expression points to an existing value
  * `if` – evaluate one of two expressions based on a condition
  * `patch` – applies an expression to all matched values
  * `set` – set a value in a variable/document, most often used with ! modifier (set!)
  * `try` – returns the fallback if the first expression errors out

* **coalesce**
  * `humanely` – evaluates the child expressions using humane coalescing
  * `pedantically` – evaluates the child expressions using pedantic coalescing
  * `strictly` – evaluates the child expressions using strict coalescing

* **compare**
  * `eq?` – equality check: return true if both arguments are the same
  * `gt?` – returns a > b
  * `gte?` – returns a >= b
  * `identical?` – like `eq?`, but always uses strict coalecsing
  * `like?` – like `eq?`, but always uses humane coalecsing
  * `lt?` – returns a < b
  * `lte?` – returns a <= b

* **datetime**
  * `now` – returns the current date & time (UTC), formatted like a Go date

* **encoding**
  * `from-base64` – decode a base64 encoded string
  * `from-json` – decode a JSON string
  * `to-base64` – apply base64 encoding to the given string
  * `to-json` – encode the given value using JSON

* **hashing**
  * `sha1` – return the lowercase hex representation of the SHA-1 hash
  * `sha256` – return the lowercase hex representation of the SHA-256 hash
  * `sha512` – return the lowercase hex representation of the SHA-512 hash

* **lists**
  * `filter` – returns a copy of a given vector/object with only those elements remaining that satisfy a condition
  * `map` – applies an expression to every element in a vector or object
  * `range` – allows to iterate (loop) over a vector or object

* **logic**
  * `and` – returns true if all arguments are true
  * `not` – negates the given argument
  * `or` – returns true if any of the arguments is true

* **math**
  * `add` – returns the sum of all of its arguments
  * `div` – returns arg1 / arg2 / .. / argN (always a floating point division, regardless of arguments)
  * `mult` – returns the product of all of its arguments
  * `sub` – returns arg1 - arg2 - .. - argN

* **strings**
  * `append` – appends more strings to a string or arbitrary items into a vector
  * `concat` – concatenates items in a vector using a common glue string
  * `contains?` – returns true if a string contains a substring or a vector contains the given element
  * `has-prefix?` – returns true if the given string has the prefix
  * `has-suffix?` – returns true if the given string has the suffix
  * `len` – returns the length of a string, vector or object
  * `prepend` – prepends more strings to a string or arbitrary items into a vector
  * `replace` – returns a copy of a string with the a substring replaced by another
  * `reverse` – reverses a string or the elements of a vector
  * `split` – splits a string into a vector
  * `to-lower` – returns the lowercased version of the given string
  * `to-upper` – returns the uppercased version of the given string
  * `trim` – returns the given whitespace with leading/trailing whitespace removed
  * `trim-prefix` – removes the prefix from the string, if it exists
  * `trim-suffix` – removes the suffix from the string, if it exists

* **types**
  * `to-bool` – try to convert the given argument losslessly to a bool
  * `to-float` – try to convert the given argument losslessly to a float64
  * `to-int` – try to convert the given argument losslessly to an int64
  * `to-string` – try to convert the given argument losslessly to a string
  * `type-of` – returns the type of a given value (e.g. "string" or "number")

* **rudifunc**
  * `func` – defines a new function

* **semver**
  * `semver` – parses a string as a semantic version

* **set**
  * `new-key-set` – create a set filled with the keys of an object
  * `new-set` – create a set filled with the given values
  * `set-delete` – returns a copy of the set with the given values removed from it
  * `set-diff` – returns the difference between two sets
  * `set-eq?` – returns true if two sets hold the same values
  * `set-has-any?` – returns true if the set contains _any_ of the given values
  * `set-has?` – returns true if the set contains _all_ of the given values
  * `set-insert` – returns a copy of the set with the newly added values inserted to it
  * `set-intersection` – returns the insersection of two sets
  * `set-list` – returns a sorted vector containing the values of the set
  * `set-size` – returns the number of values in the set
  * `set-superset-of?` – returns true if the other set is a superset of the base set
  * `set-symdiff` – returns the symmetric difference between two sets
  * `set-union` – returns the union of two or more sets

* **uuid**
  * `uuidv4` – returns a new, randomly generated v4 UUID

* **yaml**
  * `from-yaml` – decodes a YAML string into a Go value
  * `to-yaml` – encodes the given value as YAML
<!-- END_HELP_LIB_TOC -->
