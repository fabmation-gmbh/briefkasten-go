package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

// Migrations holds all registered SQL migrations.
var Migrations = migrate.NewMigrations()

//go:embed *.sql
var migFS embed.FS

func init() {
	if err := Migrations.Discover(migFS); err != nil {
		panic(err)
	}
	// if err := Migrations.DiscoverCaller(); err != nil {
	//     panic(err)
	// }
}
