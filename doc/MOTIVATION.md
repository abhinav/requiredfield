# Motivation

A common pattern in Go is to use a struct
to pass several parameters to a function.
This is often referred to as a "parameter object" or a "parameter struct".
If you're unfamiliar with the concept, you can read more about it in
[Designing Go Libraries > Parameter objects](https://abhinavg.net/2022/12/06/designing-go-libraries/#parameter-objects).

In short, the pattern provides some advantages:

1. readability:
   names of fields are visible at call sites,
   allowing them to act as a form of documentation
   similar to named parameters in other languages
2. flexibility:
   new fields can be added without updating all existing call sites

These are both desirable properties for libraries:
users of the library get a readable API
and maintainers of the library can add new **optional** fields
without a major version bump.

For applications, however,
the flexibility afforded by the pattern can turn into a problem.
Application-internal packages rarely cares about API backwards compatibility
and are prone to adding new *required* parameters to functions.
If they use parameter objects,
they lose the ability to safely add these required parameters:
they can no longer have the compiler tell them that they missed a spot.

So application developers are left to choose between:

- parameter objects: get readability, lose safety
- functions with tens of parameters: lose readability, get safety

requiredfield aims to fill this gap with parameter objects
so that applications can still get the readability benefits of using them
without sacrificing safety.
