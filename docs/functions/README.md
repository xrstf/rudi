# Rudi Standard Library

Rudi ships with a set of built-in functions. When embedding Rudi, this set can
be extended or overwritten as desired to inject custom functions. It is not
possible to define functions in an Rudi program itself.

As a rule of thumb, functions who name ends with a question mark return a boolean,
like `eq?` or `has-prefix?`. Functions with an exclamation point at the end are
not stateless but meant to modify their first argument (see the
[language spec](../language.md) regarding the bang modifier in tuples). The
question mark is part of the function name itself, but the bang modifier can be
applied to all functions (so technically `eq?!` is valid, though weird looking).

<!-- BEGIN_TOC -->
## Core Functions

* [`default`](../functions/core-default.md) – returns the default value if the first argument is empty
* [`delete`](../functions/core-delete.md) – removes a key from an object or an item from a vector
* [`do`](../functions/core-do.md) – eval a sequence of statements where only one expression is valid
* [`empty?`](../functions/core-empty.md) – returns true when the given value is empty-ish (0, false, null, "", ...)
* [`has?`](../functions/core-has.md) – returns true if the given symbol's path expression points to an existing value
* [`if`](../functions/core-if.md) – evaluate one of two expressions based on a condition
* [`set`](../functions/core-set.md) – set a value in a variable/document, only really useful with ! modifier (set!)
* [`try`](../functions/core-try.md) – returns the fallback if the first expression errors out

## Comparisons Functions

* [`eq?`](../functions/comparisons-eq.md) – equality check: return true if both arguments are the same
* [`gt?`](../functions/comparisons-gt.md) – returns a > b
* [`gte?`](../functions/comparisons-gte.md) – returns a >= b
* [`identical?`](../functions/comparisons-identical.md) – like `eq?`, but always uses strict coalecsing
* [`like?`](../functions/comparisons-like.md) – like `eq?`, but always uses humane coalecsing
* [`lt?`](../functions/comparisons-lt.md) – returns a < b
* [`lte?`](../functions/comparisons-lte.md) – return a <= b

## Dates Functions

* [`now`](../functions/dates-now.md) – returns the current date & time (UTC), formatted like a Go date

## Encoding Functions

* [`from-base64`](../functions/encoding-from-base64.md) – decode a base64 encoded string
* [`to-base64`](../functions/encoding-to-base64.md) – apply base64 encoding to the given string

## Hashes Functions

* [`sha1`](../functions/hashes-sha1.md) – return the lowercase hex representation of the SHA-1 hash
* [`sha256`](../functions/hashes-sha256.md) – return the lowercase hex representation of the SHA-256 hash
* [`sha512`](../functions/hashes-sha512.md) – return the lowercase hex representation of the SHA-512 hash

## Lists Functions

* [`append`](../functions/lists-append.md) – appends more strings to a string or arbitrary items into a vector
* [`contains?`](../functions/lists-contains.md) – returns true if a string contains a substring or a vector contains the given element
* [`filter`](../functions/lists-filter.md) – returns a copy of a given vector/object with only those elements remaining that satisfy a condition
* [`len`](../functions/lists-len.md) – returns the length of a string, vector or object
* [`map`](../functions/lists-map.md) – applies an expression to every element in a vector or object
* [`prepend`](../functions/lists-prepend.md) – prepends more strings to a string or arbitrary items into a vector
* [`range`](../functions/lists-range.md) – allows to iterate (loop) over a vector or object
* [`reverse`](../functions/lists-reverse.md) – reverses a string or the elements of a vector

## Logic Functions

* [`and`](../functions/logic-and.md) – returns true if all arguments are true
* [`not`](../functions/logic-not.md) – negates the given argument
* [`or`](../functions/logic-or.md) – returns true if any of the arguments is true

## Math Functions

* [`*`](../functions/math-mult.md) – returns the product of all of its arguments
* [`+`](../functions/math-sum.md) – returns the sum of all of its arguments
* [`-`](../functions/math-sub.md) – returns arg1 - arg2 - .. - argN
* [`/`](../functions/math-div.md) – returns arg1 / arg2 / .. / argN

## Strings Functions

* [`concat`](../functions/strings-concat.md) – concatenate items in a vector using a common glue string
* [`has-prefix?`](../functions/strings-has-prefix.md) – returns true if the given string has the prefix
* [`has-suffix?`](../functions/strings-has-suffix.md) – returns true if the given string has the suffix
* [`split`](../functions/strings-split.md) – split a string into a vector
* [`to-lower`](../functions/strings-to-lower.md) – returns the lowercased version of the given string
* [`to-upper`](../functions/strings-to-upper.md) – returns the uppercased version of the given string
* [`trim`](../functions/strings-trim.md) – returns the given whitespace with leading/trailing whitespace removed
* [`trim-prefix`](../functions/strings-trim-prefix.md) – removes the prefix from the string, if it exists
* [`trim-suffix`](../functions/strings-trim-suffix.md) – removes the suffix from the string, if it exists

## Types Functions

* [`to-bool`](../functions/types-to-bool.md) – try to convert the given argument losslessly to a bool
* [`to-float`](../functions/types-to-float.md) – try to convert the given argument losslessly to a float64
* [`to-int`](../functions/types-to-int.md) – try to convert the given argument losslessly to an int64
* [`to-string`](../functions/types-to-string.md) – try to convert the given argument losslessly to a string
* [`type-of`](../functions/types-type-of.md) – returns the type of a given value (e.g. "string" or "number")
<!-- END_TOC -->
