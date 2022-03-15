# DDD Go Template

This project contains the same logical architecture organized in two different ways:

1. The first one is the `foolproof` version, which is very simple to explain and understand.
   The drawback of this approach is that the names of the interfaces and DTOs end up a little bit verbose, e.g.:
   `domain.LogProvider` instead of `log.Provider`, and `domain.LogBody` instead of `log.Body`
2. The second one is the `advanced` one, it is harder to explain, but follows the exact same logical structure
   of the first example, and have better names and overall it just feels more "right".
   It is also more similar to how the standard library is organized, which is a plus.

That said, I recommend reading the `foolproof/README.md` first
then understanding the project organization in its simplest form
before moving to the `advanced/README.md` example.

Reading in this order should help you (1) understand all the design decisions
and (2) understand the differences and which one would fit your team better.

