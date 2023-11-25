# Documentation

Welcome to the Rudi documentation :smile:

<!-- BEGIN_TOC -->
## General

* [language](language.md) – A short introduction to the Rudi language

## Core Functions

* [`default`](functions/core-default.md) – returns the default value if the first argument is empty
* [`delete`](functions/core-delete.md) – removes a key from an object or an item from a vector
* [`do`](functions/core-do.md) – eval a sequence of statements where only one expression is valid
* [`empty?`](functions/core-empty.md) – returns true when the given value is empty-ish (0, false, null, "", ...)
* [`has?`](functions/core-has.md) – returns true if the given symbol's path expression points to an existing value
* [`if`](functions/core-if.md) – evaluate one of two expressions based on a condition
* [`set`](functions/core-set.md) – set a value in a variable/document, only really useful with ! modifier (set!)
* [`try`](functions/core-try.md) – returns the fallback if the first expression errors out

## Math Functions

* [`*`](functions/math-mult.md) – returns the product of all of its arguments
* [`+`](functions/math-sum.md) – returns the sum of all of its arguments
* [`-`](functions/math-sub.md) – returns arg1 - arg2 - .. - argN
* [`/`](functions/math-div.md) – returns arg1 / arg2 / .. / argN

## Strings Functions

* [`concat`](functions/strings-concat.md) – concatenate items in a vector using a common glue string
* [`split`](functions/strings-split.md) – split a string into a vector
<!-- END_TOC -->
