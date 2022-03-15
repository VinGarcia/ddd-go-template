# DDD Go Template - Advanced version

If you haven't read yet, I recommend reading the `foolproof/README.md` first.

This example template has the same logical structure as the `foolproof` version,
but organizes the interfaces and DTOs in a way that allows better names, such as:

- `rest.Provider` instead of `domain.RestProvider`
- `cache.Provider` instead of `domain.CacheProvider`

## Reorganizing the infra package

For that we reorganized all adapters in the `infra/` package by nesting
them inside a package that contains the interface they implement, e.g.:

- The `infra/http` package was moved to `infra/rest/http`
- The `infra/memorycache` package was moved to `infra/cache/memorycache`
- The `infra/redis` package was moved to `infra/cache/redis`

And so on.

For each of these new packages a new file `contracts.go` was created containing
only the relevant interfaces for that dependency, so now we have 3 new files with that name:

- `infra/rest/contracts.go`
- `infra/cache/contracts.go`
- `infra/cache/contracts.go`

And the old `domain/contracts.go` was deleted.
