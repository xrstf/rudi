# Type Handling & Conversions

Rudi has pluggable coalescers, which is the term it uses for the component that is responsible for
handling implicit type conversions (e.g. deciding whether `1` (integer) is the same as `"1"` (string))
and by extension also for equality checks between two values.

* [Background](#background)
* [Coalescers](#coalescers)
  + [Strict](#strict-coalescer)
  + [Pedantic](#pedantic-coalescer)
  + [Humane](#humane-coalescer)
* [Conversion Functions](#conversion-functions)
* [Comparisons](#comparisons)

## Background

Rudi programs are meant to transform data structures, often YAML/JSON files. In some cases, these
are based on strictly typed datastructures, like [Kubernetes](https://kubernetes.io) objects that
are based on an OpenAPI schema. Sometimes however these data structures are much more loosey goosey,
like [Helm](https://helm.sh/) values, where a flag like `enabled: "false"` would ultimately be rendered
using `--enabled={{ .enabled }}` to `--enabled=false`. In these cases, types are less important than
the user's actual intent.

The built-in functions in Rudi already contain helpers like [`to-string`](functions/types-to-string.md)
or [`to-int`](functions/types-to-int.md) to deal with conversions, but it can be tedious to have
type conversions in many places in a Rudi program, just to deal with untyped data.

To help with this, Rudi offers a choice, both to the Rudi program author by having, as well as
when embedding Rudi into other Go programs: Choose your own adventure, or, "coalescer".

When running `rudi`, the coalescing can be changed using `--coalesce`:

```bash
$ rudi --coalesce humane '(+ .foo 2)' data.yaml
```

## Coalescers

A coalescer is a Go interface similar to this:

```go
type Coalescer interface {
   ToBool(val any) (bool, error)
   ToInt64(val any) (int64, error)
   ToString(val any) (string, error)
   // ...
}
```

Its task is to handle type conversions/ensurances across all Rudi functions. When `(add 1 2)` wants
to turn its two arguments into numbers, the coalescer is used. The coalescer decides which values
of which types to convert into the desired target type.

Rudi comes with 3 coalescers to choose from:

* **Strict** is the default. This coalescer does not allow any type conversions, except turning
  `null` into the empty value of any type (i.e. `null` can turn int `0`) and allowing to turn
  floating point numbers into integer numbers if no precision is lost (i.e. `2.0` can turn into `2`).
* **Pedantic** is even more strict than **strict** and does not allow any type conversion whatsoever.
* **Humane** is inspired by, but less forgiving than PHP's type system. This coalescer allows much
  more conversions, like turning `"2.0"` (string) into `2` (int).

You can of course also implement your own coalescer by implementing `pkg/coalescing.Coalescer`.

### Strict Coalescer

A safe default with minimal conversion support. The following list shows which values/types are
allowed for each target type. The left column shows the target type, the other columns list
acceptable source types

|         | null      | bool | int64 | float64 | string | vector | object |
| ------: | :-------: | :--: | :---: | :-----: | :----: | :----: | :----: |
| null    | ✅        | –    | –     | –       | –      | –      | –      |
| bool    | ✅        | ✅   | –     | –       | –      | –      | –      |
| int64   | ✅        | –    | ✅    | ✅(*)   | –      | –      | –      |
| float64 | ✅        | –    | ✅    | ✅      | –      | –      | –      |
| string  | ✅ (`""`) | –    | –     | –       | ✅     | –      | –      |
| vector  | ✅        | –    | –     | –       | –      | ✅     | –      |
| object  | ✅        | –    | –     | –       | –      | –      | ✅     |

\*) only if lossless (e.g. `2.0` can be turned into an int64, `2.1` cannot)

### Pedantic Coalescer

If you really want to be super extra explicit and have strongly typed source data, maybe the pedantic
coalescer is more your style.

|         | null  | bool | int64 | float64 | string | vector | object |
| ------: | :---: | :--: | :---: | :-----: | :----: | :----: | :----: |
| null    | ✅    | –    | –     | –       | –      | –      | –      |
| bool    | –     | ✅   | –     | –       | –      | –      | –      |
| int64   | –     | –    | ✅    | –       | –      | –      | –      |
| float64 | –     | –    | –     | ✅      | –      | –      | –      |
| string  | –     | –    | –     | –       | ✅     | –      | –      |
| vector  | –     | –    | –     | –       | –      | ✅     | –      |
| object  | –     | –    | –     | –       | –      | –      | ✅     |

### Humane Coalescer

For less-than-strongly typed data, sometimes it's easier to just accept humans for what they are and
that `replicas: "2"` really meant `replicas: 2`. For cases like this, use the humane coalescer, which
is inspired by PHP, but a bit less flexible. Also instead of turning `true` into `"1"` like in PHP,
this coalescer returns `"true"` (and `"false"` for `false`).

|         | null  | bool | int64 | float64 | string | vector | object |
| ------: | :---: | :--: | :---: | :-----: | :----: | :----: | :----: |
| null    | ✅    | ✅   | ✅    | ✅      | ✅     | ✅     | ✅     |
| bool    | ✅    | ✅   | ✅    | ✅      | ✅     | ✅     | ✅     |
| int64   | ✅    | ✅   | ✅    | ✅      | ✅     | –      | –      |
| float64 | ✅    | ✅   | ✅    | ✅      | ✅     | –      | –      |
| string  | ✅    | ✅   | ✅    | ✅      | ✅     | –      | –      |
| vector  | ✅    | –    | –     | –       | –      | ✅     | ✅     |
| object  | ✅    | –    | –     | –       | –      | ✅     | ✅     |

Let's look closer at each type's conversion logic:

* **null**: All conversions to `null` are only allowed for the empty value of each source type
  (meaning that `false` and `0` are convertible, but `true` and `"foo"` are not).
* **bool**: Empty values are considered `false`, all others `true`. The string `"0"` and `"false"`
  are also considered `false`.
* **int64**: Empty values are `0`, `true` becomes `1`. If lossless conversion from float64 to int64
  is possible, it's performed, otherwise an error is returned. Strings have their whitespace trimmed;
  if the resulting string is empty, `0` is returned, otherwise the string is parsed as an integer;
  if not successful, parsing as float is attempted and if the resulting float can be losslessly
  converted, it's returned (e.g. `"2.0"` is valid, `"2.1"` is not). Otherwise an error is returned.
* **float64**: Empty values are `0.0`, `true` becomes `1.0`. Integers are converted to floats.
  Strings have their whitespace trimmed; if the resulting string is empty, `0.0` is returned,
  otherwise the string is parsed as a float; if not possible to parse as float, an error is returned.
* **string**: works as expected; floats have their trailing zeros trimmed (`3.12000` becomes
  `"3.12"`). `true` becomes `"true"` and `false` becomes `"false"`. `null` becomes an empty string.
* **vector**: `null` turns into an empty vector and objects can only be converted to an empty vector
  if they are empty, otherwise an error is retuned.
* **object**: `null` turns into an empty object and vectors can only be converted to an empty object
  if they are empty, otherwise an error is retuned.

## Conversion Functions

Rudi offers explicit conversion functions. These always apply the humane coalescing logic.

* [`to-int`](functions/types-to-int.md) converts its argument to an int64 or returns an error if not possible.
* [`to-float`](functions/types-to-float.md) does the same for float64.
* [`to-string`](functions/types-to-string.md) does the same for strings.
* [`to-bool`](functions/types-to-bool.md) does the same for booleans.

## Comparisons

Rudi has 3 functions built-in to check for equality between 2 values:

* [`eq?`](functions/comparisons-eq.md) uses the current coalescer (i.e. by default, strict). If a Rudi program is configured to
  use humane coalescing however, this function will use that coalescing to determine equality.
* [`like?`](functions/comparisons-like.md) always uses humane coalescing.
* [`identical?`](functions/comparisons-identical.md) always uses strict coalescing.

Comparisons work by converting the two values into (hopefully) compatible types that can be compared.
This is done in steps:

1. If either of the arguments is `null`, try to convert to other to `null`.
1. Do the same with `bool`.
1. Do the same with `int64`.
1. Do the same with `float64`.
1. Do the same with `string`.
1. Do the same with `vector`.
1. Do the same with `object`.

**NB:** Equality rules are associative (if `a == b`, then `b == a`), but not transitive, which is
especially apparent with humane coalescing:

* `" " == true` because the string is not empty.
* `" " == 0` because empty strings can turn into `0`.
* `0 == false` because both are the empty values of their types.

If rules were transitive, `0` could not both be `false` and `true` at the same time.
