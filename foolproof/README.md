# DDD Go Template - Fool Proof version

This project was created to illustrate an interesting directory structure
I developed together with [@fabiorodrigues](https://github.com/fabiorodrigues) in the period I was working
for [Dito](https://dito.com.br), they both deserve as much credit as me here.

This very powerful, but yet flat and simple template is organized in 3 directories:

- **cmd/:** Each subdirectory is an entry point for the project,
  e.g. a worker, an API or a CLI. Each of these packages
  is responsible for decoding the configurations, performing the
  dependency injection and setting up any Frameworks if necessary
  (in our case we are using the Fiber as our HTTP framework).

- **domain/:** This package contains the domain language, which is the minimum
  shared language that all packages are allowed to import. Thus, this package is
  meant to be imported by all other packages in order to allow decoupled
  communication between them.

  Each subpackage of the domain pkg is a Service, and this is where we
  should concentrate the domain logic.

- **infra/:** each subdirectory here contains an adapter, i.e. some code
  that adapts an external dependency or logic that is unrelated to your domain
  to an interface declared on `domain/contracts.go`.

  You can also have small helper packages here if necessary for operations that
  are so simple that there is no need to rely on an external dependency.

  These infra packages are meant to contain any logic that is unrelated
  to your domain in order to move as much code as possible away from your services.

  One other thing that we keep here are the repositories, which are often the only
  infra packages that actually use the entities directly, although this is not prohibited
  by DDD.

  The idea here is that the Services contain the most complex and important parts of the project,
  thus, by moving any logic that is not related to the domain away from the Services we can
  keep the Services as simple as they can possibly be, which makes the code a lot easier to maintain.

For Portuguese readers we have a more descriptive explanation of this architecture here:

- [Domain Driven Design Aplicado a um Microserviço Go](https://eng.dito.com.br/domain-driven-design-ddd-aplicado-a-um-microservico-go)

And if you prefer to watch a presentation (also in Portuguese) we have this one:

- [29º Go Talks BH](https://youtu.be/ODft0k1LeHU)
