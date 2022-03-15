package ksqlrepo

import "github.com/vingarcia/ksql"

// For each new repositóry just add it here as an unnamed attribute
// so that we have a single struct implementing all the repositories
// that are implemented using ksql
//
// In the domain we'll only use the interfaces, so no one will know
// it's the same struct implementing multiple interfaces.
type Repo struct {
	UsersRepo
}

func New(db ksql.Provider) Repo {
	return Repo{
		UsersRepo: newUsersRepo(db),
	}
}
