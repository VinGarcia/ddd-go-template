# DDD Go Template

This project contains the same logical architecture organized in two different ways:

1. The first one is the `v1-very-simple` version, which is very simple to explain and understand.
   The drawback of this approach is that the names of the interfaces and DTOs end up a little bit verbose, e.g.:
   `domain.LogProvider` instead of `log.Provider`, and `domain.LogBody` instead of `log.Body`
2. The second one is the `v2-domain-adapters-and-helpers`,
   this one slightly less simple because the interfaces are now spread on
   different places, but this allows for better names for the interfaces and
   it also helps to organize your code if you have a few big interfaces and you
   don't want all of them in the same file.

That said, I recommend reading the `v1-very-simple/README.md` first
then understanding the project organization in its simplest form
before moving to the `v2-domain-adapters-and-helpers/README.md` example.

Reading in this order should help you (1) understand all the design decisions
and (2) understand the differences and which one would fit your team better.

