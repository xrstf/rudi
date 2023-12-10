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

Additional functionality is available in the [extended library](../extlib/).

<!-- BEGIN_STDLIB_TOC -->
### core

* [`default`](../stdlib/core/default.md) – returns the default value if the first argument is empty
* [`delete`](../stdlib/core/delete.md) – removes a key from an object or an item from a vector
* [`do`](../stdlib/core/do.md) – eval a sequence of statements where only one expression is valid
* [`empty?`](../stdlib/core/empty.md) – returns true when the given value is empty-ish (0, false, null, "", ...)
* [`error`](../stdlib/core/error.md) – returns an error
* [`has?`](../stdlib/core/has.md) – returns true if the given symbol's path expression points to an existing value
* [`if`](../stdlib/core/if.md) – evaluate one of two expressions based on a condition
* [`set`](../stdlib/core/set.md) – set a value in a variable/document, only really useful with ! modifier (set!)
* [`try`](../stdlib/core/try.md) – returns the fallback if the first expression errors out

### coalesce

* [`humanely`](../stdlib/coalesce/humanely.md) – evaluates the child expressions using humane coalescing
* [`pedantically`](../stdlib/coalesce/pedantically.md) – evaluates the child expressions using pedantic coalescing
* [`strictly`](../stdlib/coalesce/strictly.md) – evaluates the child expressions using strict coalescing

### compare

* [`eq?`](../stdlib/compare/eq.md) – equality check: return true if both arguments are the same
* [`gt?`](../stdlib/compare/gt.md) – returns a > b
* [`gte?`](../stdlib/compare/gte.md) – returns a >= b
* [`identical?`](../stdlib/compare/identical.md) – like `eq?`, but always uses strict coalecsing
* [`like?`](../stdlib/compare/like.md) – like `eq?`, but always uses humane coalecsing
* [`lt?`](../stdlib/compare/lt.md) – returns a < b
* [`lte?`](../stdlib/compare/lte.md) – returns a <= b

### datetime

* [`now`](../stdlib/datetime/now.md) – returns the current date & time (UTC), formatted like a Go date

### encoding

* [`from-base64`](../stdlib/encoding/from-base64.md) – decode a base64 encoded string
* [`from-json`](../stdlib/encoding/from-json.md) – decode a JSON string
* [`to-base64`](../stdlib/encoding/to-base64.md) – apply base64 encoding to the given string
* [`to-json`](../stdlib/encoding/to-json.md) – encode the given value using JSON

### hashing

* [`sha1`](../stdlib/hashing/sha1.md) – return the lowercase hex representation of the SHA-1 hash
* [`sha256`](../stdlib/hashing/sha256.md) – return the lowercase hex representation of the SHA-256 hash
* [`sha512`](../stdlib/hashing/sha512.md) – return the lowercase hex representation of the SHA-512 hash

### lists

* [`filter`](../stdlib/lists/filter.md) – returns a copy of a given vector/object with only those elements remaining that satisfy a condition
* [`map`](../stdlib/lists/map.md) – applies an expression to every element in a vector or object
* [`range`](../stdlib/lists/range.md) – allows to iterate (loop) over a vector or object

### logic

* [`and`](../stdlib/logic/and.md) – returns true if all arguments are true
* [`not`](../stdlib/logic/not.md) – negates the given argument
* [`or`](../stdlib/logic/or.md) – returns true if any of the arguments is true

### math

* [`add`](../stdlib/math/add.md) – returns the sum of all of its arguments
* [`div`](../stdlib/math/div.md) – returns arg1 / arg2 / .. / argN (always a floating point division, regardless of arguments)
* [`mult`](../stdlib/math/mult.md) – returns the product of all of its arguments
* [`sub`](../stdlib/math/sub.md) – returns arg1 - arg2 - .. - argN

### strings

* [`append`](../stdlib/strings/append.md) – appends more strings to a string or arbitrary items into a vector
* [`concat`](../stdlib/strings/concat.md) – concatenates items in a vector using a common glue string
* [`contains?`](../stdlib/strings/contains.md) – returns true if a string contains a substring or a vector contains the given element
* [`has-prefix?`](../stdlib/strings/has-prefix.md) – returns true if the given string has the prefix
* [`has-suffix?`](../stdlib/strings/has-suffix.md) – returns true if the given string has the suffix
* [`len`](../stdlib/strings/len.md) – returns the length of a string, vector or object
* [`prepend`](../stdlib/strings/prepend.md) – prepends more strings to a string or arbitrary items into a vector
* [`replace`](../stdlib/strings/replace.md) – returns a copy of a string with the a substring replaced by another
* [`reverse`](../stdlib/strings/reverse.md) – reverses a string or the elements of a vector
* [`split`](../stdlib/strings/split.md) – splits a string into a vector
* [`to-lower`](../stdlib/strings/to-lower.md) – returns the lowercased version of the given string
* [`to-upper`](../stdlib/strings/to-upper.md) – returns the uppercased version of the given string
* [`trim`](../stdlib/strings/trim.md) – returns the given whitespace with leading/trailing whitespace removed
* [`trim-prefix`](../stdlib/strings/trim-prefix.md) – removes the prefix from the string, if it exists
* [`trim-suffix`](../stdlib/strings/trim-suffix.md) – removes the suffix from the string, if it exists

### types

* [`to-bool`](../stdlib/types/to-bool.md) – try to convert the given argument losslessly to a bool
* [`to-float`](../stdlib/types/to-float.md) – try to convert the given argument losslessly to a float64
* [`to-int`](../stdlib/types/to-int.md) – try to convert the given argument losslessly to an int64
* [`to-string`](../stdlib/types/to-string.md) – try to convert the given argument losslessly to a string
* [`type-of`](../stdlib/types/type-of.md) – returns the type of a given value (e.g. "string" or "number")
<!-- END_STDLIB_TOC -->
