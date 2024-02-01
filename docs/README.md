# Documentation

Welcome to the Rudi documentation :smile:

## General

<!-- BEGIN_TOPICS_TOC -->
* [The Rudi Language](language.md) – A short introduction to the Rudi language
* [Type Handling & Conversions](coalescing.md) – How Rudi handles, converts and compares values
<!-- END_TOPICS_TOC -->

## Standard Library

These are all the built-in functions, available in the `rudi` interpreter and when embedding Rudi
and not specifying another, custom function set.

<!-- BEGIN_STDLIB_TOC -->
### core

* [`case`](stdlib/core/case.md) – chooses the first expression for which the test is true
* [`default`](stdlib/core/default.md) – returns the default value if the first argument is empty
* [`delete`](stdlib/core/delete.md) – removes a key from an object or an item from a vector
* [`do`](stdlib/core/do.md) – eval a sequence of statements where only one expression is valid
* [`empty?`](stdlib/core/empty.md) – returns true when the given value is empty-ish (0, false, null, "", ...)
* [`error`](stdlib/core/error.md) – returns an error
* [`has?`](stdlib/core/has.md) – returns true if the given symbol's path expression points to an existing value
* [`if`](stdlib/core/if.md) – evaluate one of two expressions based on a condition
* [`set`](stdlib/core/set.md) – set a value in a variable/document, only really useful with ! modifier (set!)
* [`try`](stdlib/core/try.md) – returns the fallback if the first expression errors out

### coalesce

* [`humanely`](stdlib/coalesce/humanely.md) – evaluates the child expressions using humane coalescing
* [`pedantically`](stdlib/coalesce/pedantically.md) – evaluates the child expressions using pedantic coalescing
* [`strictly`](stdlib/coalesce/strictly.md) – evaluates the child expressions using strict coalescing

### compare

* [`eq?`](stdlib/compare/eq.md) – equality check: return true if both arguments are the same
* [`gt?`](stdlib/compare/gt.md) – returns a > b
* [`gte?`](stdlib/compare/gte.md) – returns a >= b
* [`identical?`](stdlib/compare/identical.md) – like `eq?`, but always uses strict coalecsing
* [`like?`](stdlib/compare/like.md) – like `eq?`, but always uses humane coalecsing
* [`lt?`](stdlib/compare/lt.md) – returns a < b
* [`lte?`](stdlib/compare/lte.md) – returns a <= b

### datetime

* [`now`](stdlib/datetime/now.md) – returns the current date & time (UTC), formatted like a Go date

### encoding

* [`from-base64`](stdlib/encoding/from-base64.md) – decode a base64 encoded string
* [`from-json`](stdlib/encoding/from-json.md) – decode a JSON string
* [`to-base64`](stdlib/encoding/to-base64.md) – apply base64 encoding to the given string
* [`to-json`](stdlib/encoding/to-json.md) – encode the given value using JSON

### hashing

* [`sha1`](stdlib/hashing/sha1.md) – return the lowercase hex representation of the SHA-1 hash
* [`sha256`](stdlib/hashing/sha256.md) – return the lowercase hex representation of the SHA-256 hash
* [`sha512`](stdlib/hashing/sha512.md) – return the lowercase hex representation of the SHA-512 hash

### lists

* [`filter`](stdlib/lists/filter.md) – returns a copy of a given vector/object with only those elements remaining that satisfy a condition
* [`map`](stdlib/lists/map.md) – applies an expression to every element in a vector or object
* [`range`](stdlib/lists/range.md) – allows to iterate (loop) over a vector or object

### logic

* [`and`](stdlib/logic/and.md) – returns true if all arguments are true
* [`not`](stdlib/logic/not.md) – negates the given argument
* [`or`](stdlib/logic/or.md) – returns true if any of the arguments is true

### math

* [`add`](stdlib/math/add.md) – returns the sum of all of its arguments
* [`div`](stdlib/math/div.md) – returns arg1 / arg2 / .. / argN (always a floating point division, regardless of arguments)
* [`mult`](stdlib/math/mult.md) – returns the product of all of its arguments
* [`sub`](stdlib/math/sub.md) – returns arg1 - arg2 - .. - argN

### strings

* [`append`](stdlib/strings/append.md) – appends more strings to a string or arbitrary items into a vector
* [`concat`](stdlib/strings/concat.md) – concatenates items in a vector using a common glue string
* [`contains?`](stdlib/strings/contains.md) – returns true if a string contains a substring or a vector contains the given element
* [`has-prefix?`](stdlib/strings/has-prefix.md) – returns true if the given string has the prefix
* [`has-suffix?`](stdlib/strings/has-suffix.md) – returns true if the given string has the suffix
* [`len`](stdlib/strings/len.md) – returns the length of a string, vector or object
* [`prepend`](stdlib/strings/prepend.md) – prepends more strings to a string or arbitrary items into a vector
* [`replace`](stdlib/strings/replace.md) – returns a copy of a string with the a substring replaced by another
* [`reverse`](stdlib/strings/reverse.md) – reverses a string or the elements of a vector
* [`split`](stdlib/strings/split.md) – splits a string into a vector
* [`to-lower`](stdlib/strings/to-lower.md) – returns the lowercased version of the given string
* [`to-upper`](stdlib/strings/to-upper.md) – returns the uppercased version of the given string
* [`trim`](stdlib/strings/trim.md) – returns the given whitespace with leading/trailing whitespace removed
* [`trim-prefix`](stdlib/strings/trim-prefix.md) – removes the prefix from the string, if it exists
* [`trim-suffix`](stdlib/strings/trim-suffix.md) – removes the suffix from the string, if it exists

### types

* [`to-bool`](stdlib/types/to-bool.md) – try to convert the given argument losslessly to a bool
* [`to-float`](stdlib/types/to-float.md) – try to convert the given argument losslessly to a float64
* [`to-int`](stdlib/types/to-int.md) – try to convert the given argument losslessly to an int64
* [`to-string`](stdlib/types/to-string.md) – try to convert the given argument losslessly to a string
* [`type-of`](stdlib/types/type-of.md) – returns the type of a given value (e.g. "string" or "number")

### rudifunc

* [`func`](stdlib/rudifunc/func.md) – defines a new function
<!-- END_STDLIB_TOC -->

## Extended Library

These modules are only available when explicitly importing their Go modules and adding them to the
Rudi function set. They are however all available by default in the `rudi` interpreter.

<!-- BEGIN_EXTLIB_TOC -->
### semver

* [`semver`](extlib/semver/semver.md) – parses a string as a semantic version

### set

* [`new-key-set`](extlib/set/new-key-set.md) – create a set filled with the keys of an object
* [`new-set`](extlib/set/new-set.md) – create a set filled with the given values
* [`set-delete`](extlib/set/set-delete.md) – returns a copy of the set with the given values removed from it
* [`set-diff`](extlib/set/set-diff.md) – returns the difference between two sets
* [`set-eq?`](extlib/set/set-eq.md) – returns true if two sets hold the same values
* [`set-has-any?`](extlib/set/set-has-any.md) – returns true if the set contains _any_ of the given values
* [`set-has?`](extlib/set/set-has.md) – returns true if the set contains _all_ of the given values
* [`set-insert`](extlib/set/set-insert.md) – returns a copy of the set with the newly added values inserted to it
* [`set-intersection`](extlib/set/set-intersection.md) – returns the insersection of two sets
* [`set-list`](extlib/set/set-list.md) – returns a sorted vector containing the values of the set
* [`set-size`](extlib/set/set-size.md) – returns the number of values in the set
* [`set-superset-of?`](extlib/set/set-superset-of.md) – returns true if the other set is a superset of the base set
* [`set-symdiff`](extlib/set/set-symdiff.md) – returns the symmetric difference between two sets
* [`set-union`](extlib/set/set-union.md) – returns the union of two or more sets

### uuid

* [`uuidv4`](extlib/uuid/uuidv4.md) – returns a new, randomly generated v4 UUID

### yaml

* [`from-yaml`](extlib/yaml/from-yaml.md) – decodes a YAML string into a Go value
* [`to-yaml`](extlib/yaml/to-yaml.md) – encodes the given value as YAML
<!-- END_EXTLIB_TOC -->
