package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	PgxMigration "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	Context  context.Context
	Database *pgxpool.Pool
}

func (database *Database) Init(dsn string) {
	ctx := context.Background()
	cfg, errConfig := pgxpool.ParseConfig(dsn)
	if errConfig != nil {
		fmt.Println("Bad!")
	}

	dbConn, errConnectConfig := pgxpool.ConnectConfig(ctx, cfg)
	if errConnectConfig != nil {
		log.Fatalf("Failed to connect to database: %v", errConnectConfig)
	}
	database.Database = dbConn
	database.Context = ctx
}

func RunMigrate(dsn string) {
	instance, errOpen := sql.Open("pgx", dsn)
	if errOpen != nil {
		fmt.Println("Failed to open database for migration")
	}
	if errPing := instance.Ping(); errPing != nil {
		fmt.Println("Cannot migrate, failed to connect to target server")
	}

	driver, _ := PgxMigration.WithInstance(instance, &PgxMigration.Config{
		MigrationsTable:       "_migration",
		SchemaName:            "public",
		StatementTimeout:      60 * time.Second,
		MultiStatementEnabled: true,
	})

	migrate, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"pgx", driver,
	)

	if err != nil {
		fmt.Println("Failed to create migration")
	}
	err = migrate.Up()
	if err != nil {
		fmt.Println("Failed to run migrations: ", err)
	}
}