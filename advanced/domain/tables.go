package domain

import "github.com/vingarcia/ksql"

var UsersTable = ksql.NewTable("users", "id")
