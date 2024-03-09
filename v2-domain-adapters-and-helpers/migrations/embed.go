package migrations

import "embed"

//go:embed *.sql
var Dir embed.FS
