# DDD Go Template

This project was created to illustrate a great architectural structure
I developed together with [@fabiorodrigues](https://github.com/fabiorodrigues) in the period I was working
for [Dito](https://dito.com.br), they both deserve as much credit as me here.

This very powerful, but yet flat and simple, template is organized in 3 directories:

- **cmd/:** Each subdirectory is an entrypoint for the project,
  e.g. a worker, an API or a CLI interface, each of these packages
  are responsible for decoding the configurations, performing the
  dependency injection and setting any Frameworks that might need to be setup
  (in our case we are using the fasthttp as an HTTP framework).

- **domain/:** This package contains the domain language and is meant to
  be imported by all other packages in order to allow a decoupled comunication
  between them.

  Each subpackage of the domain is a Service, and its where we should concentrate
  the domain logic.

- **infra/:** each subdirectory contains either an adapter pattern making
  an external functionality available to the domain in the form of an interface
  (check the memorycache package for an example on this) or simple packages that
  extract logic that is unrelated to the to domain in order to move as much code
  away from the services as possible.

  One other thing that we keep here are the repositories, which are the often the only
  infra packages that actually use the entities directly, althought this is not prohibited
  by DDD.

  The idea here is that the Services contain the most complex and important parts of the project,
  thus, by moving any logic that is not related to the domain away from the Services we can
  keep the Services as simple as they can possibly be, which makes the code a lot easier to maintain.

For portuguese readers we have a more descriptive explanation of this architecture here:

- [Domain Driven Design Aplicado a um Microserviço Go](https://eng.dito.com.br/domain-driven-design-ddd-aplicado-a-um-microservico-go)

And if you prefer to watch a presentation (also in portuguese) we have this one:

- [29º Go Talks BH](https://youtu.be/ODft0k1LeHU)
