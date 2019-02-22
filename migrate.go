package main

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func startMigration() error {
	pgPort := 5432
	// if os.Getenv("PGPORT") != "" {
	// 	pgPort = strconv.Atoi(os.Getenv("PGPORT"))
	// }
	pgURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&binary_parameters=yes", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), pgPort, os.Getenv("PGDATABASE"))

	// driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	// m, err := migrate.NewWithDatabaseInstance("file://./migrations/", "postgres", driver)

	m, err := migrate.New("file://./migrations", pgURI)

	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Up()

	switch {
	case err == migrate.ErrNoChange:
		return nil
	case err != nil:
		return err
	}

	return nil
}
