package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	upCmd := flag.Bool("up", false, "Migrate the DB to the most recent version")
	downCmd := flag.Bool("down", false, "Rollback the DB to the previous version")
	version := flag.Int("version", 0, "Migrate to a specific version")
	flag.Parse()

	if !*upCmd && !*downCmd && *version == 0 {
		flag.Usage()
		os.Exit(1)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if *upCmd {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		log.Println("Migration completed successfully")
	} else if *downCmd {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Rollback completed successfully")
	} else if *version > 0 {
		if err := m.Migrate(uint(*version)); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate to version %d: %v", *version, err)
		}
		log.Printf("Migration to version %d completed successfully", *version)
	}
}
