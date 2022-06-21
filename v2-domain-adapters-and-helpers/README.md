# DDD Go Template - Domain Adapters and Helpers

If you haven't read yet, I recommend reading the `v1-very-simple/README.md` first.

This template as the other two is based on the same ideas, so it is useful to
learn the simplest version first.

This example template has the same logical structure as the `v1-very-simple` version,
but organizes the interfaces and DTOs in a way that allows better names, such as:

- `rest.Provider` instead of `domain.RestProvider`
- `cache.Provider` instead of `domain.CacheProvider`

Another difference is that on the `v1-very-simple` version we put both
helpers and adapters inside the `infra/` directory, but in this version
we actually split the infra directory into 2 new directories:

- `adapters/`
- `helpers/`

This is an important distinction because helpers are meant to be very simple
pieces of code that simplify a few common use-cases, and we usually don't mind
depending directly on them, exactly because they are very simple and because
they exist inside on the same version control of our code, so we can safely
update them if necessary.

But the same is not valid for adapters, so having this separation expressed
on the directory structure is kind of important.

## Reorganizing the infra package

For that we reorganized all adapters in the `infra/` package by nesting
them inside a package that contains the interface they implement, e.g.:

- The `infra/http` package was moved to `adapters/rest/http`
- The `infra/memorycache` package was moved to `adapters/cache/memorycache`
- The `infra/redis` package was moved to `adapters/cache/redis`

And so on.

For each of these new packages a new file `contracts.go` was created containing
only the relevant interfaces for that dependency, so now we have 3 new files with that name:

- `adapters/rest/contracts.go`
- `adapters/cache/contracts.go`
- `adapters/cache/contracts.go`

And the old `domain/contracts.go` was deleted.

We also moved the two helpers that we used to have to the helpers directory:

- The `infra/env` pkg was moved to `helpers/env`
- The `infra/maps` pkg was moved to `helpers/maps`

