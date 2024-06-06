package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/utils"
)

func Migrate() {
	if db == nil {
		utils.Logger().Fatal("Can't perform migration without a database connection. Call \"Connect()\" before calling \"Migrate()\"")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		utils.Logger().Fatalf("Error initiating postgres driver: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", env.Values().DB_NAME, driver)
	if err != nil {
		utils.Logger().Fatalf("Error initiating migrate instace: %v\n", err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		utils.Logger().Fatalf("Error migrating: %v\n", err)
	}
}
