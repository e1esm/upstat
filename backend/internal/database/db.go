package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (

	//go:embed migrations/sqlite/*.sql
	embedMigrationsSqlite embed.FS

	//go:embed migrations/postgres/*.sql
	embedMigrationsPostgres embed.FS

	postgres = "postgres"
	sqlite   = "sqlite"
)

func DBConnect() (*sql.DB, error) {
	var DB *sql.DB

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = sqlite
	}
	var err error

	switch dbType {
	case postgres:
		DB, err = PostgresConnection()
	default:
		DB, err = SqliteConnection()
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	switch dbType {
	case postgres:
		goose.SetBaseFS(embedMigrationsPostgres)
	default:
		goose.SetBaseFS(embedMigrationsSqlite)
	}

	if err := goose.SetDialect(dbType); err != nil {
		panic(err)
	}

	dbMigrationDir := fmt.Sprintf("migrations/%s", dbType)
	if err := goose.Up(DB, dbMigrationDir); err != nil {
		panic(err)
	}

	return DB, nil
}
