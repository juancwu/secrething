package database

import (
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/juancwu/konbini/env"
)

func Migrate() {
	if db == nil {
		log.Fatal("Can't perform migration without a database connection. Call \"Connect()\" before calling \"Migrate()\"")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error initiating postgres driver: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", env.Values().DB_NAME, driver)
	if err != nil {
		log.Fatalf("Error initiating migrate instace: %v\n", err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatalf("Error migrating: %v\n", err)
	}
}
