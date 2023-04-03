# FAQ

## Why a comment instead of a struct tag?

The reasons for this choice are both, philosophical and cosmetic.

First, the philosophical reason:
requiredfield is a linter that runs at compile-time,
and therefore wants its footprint limited to compile-time only.
Struct tags get compiled into your binary
and are available at runtime via reflection.
It would become possible for someone to
change how the program behaves based on the value of those struct tags.
requiredfield considers that a violation of the linter's boundaries,
and aims to prevent that by using comments instead.

The cosmetic reason is much easier to explain:
Struct tags are uglier than line comments.

```go
Author ID `required:"true"`

// versus

Author ID // required
```
