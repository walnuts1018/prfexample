package migrate

import "embed"

//go:embed *.cql
var Scripts embed.FS
