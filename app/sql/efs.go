package sql

import "embed"

//go:embed migrations
var SqlFiles embed.FS
