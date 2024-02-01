# case

`case` implements the common switch/case pattern: The function expects an even
number of expressions and will sequentially process them in pairs. If the first
expression yields `true`, `case` will return the second expression. Otherwise
the next pair is checked. If none of the pairs match, `case` returns `null`.

## Examples

* `(case true 1 false)` ➜ invalid
* `(case true 1 false 2)` ➜ `1`
* `(case false 1 true 2)` ➜ `2`
* `(case false 1 (eq? 1 2) 2)` ➜ `null`

## Forms

### `(case exprs:expression…)` ➜ `any`

* `exprs` is an even number of expressions.
