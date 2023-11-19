# Welcome to the Otto interpreter :)

You can enter one of

* A path expression, like `.foo` or `.foo[0].bar` to access the global document.
* An expression like (+ .foo 42) to compute data by functions; see the topics
  below or the Otto website for a complete list of available functions.
* A scalar JSON value, like `3` or `[1 2 3]`, which will simply return that
  exact value with no further side effects. Not super useful usually.

## Commands

Additionally, the following commands can be used:

* help       – Show this help text.
* help TOPIC – Show help for a specific topic.
* exit       – Exit Otti immediately.

## Help Topics

The following topics are available and can be accessed using `help TOPIC`:
