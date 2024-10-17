package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/xYurii/Bell/src/database"
)

//go:embed migrations/*.sql
var sqlMigrations embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a migration action\nmigrate: run the migration\nrollback: rollback the last migration\ninit: create migrations .sql files")
		return
	}

	ctx := context.Background()
	db, _ := database.InitDatabase(database.GetEnvDatabaseConfig())

	switch os.Args[1] {
	case "migrate":
		err := runMigrations(ctx, db)
		if err != nil {
			log.Fatal(err)
		}
	case "rollback":
		err := RollbackMigrations(ctx, db)
		if err != nil {
			log.Fatal(err)
		}
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a migration name")
			return
		}
		name := os.Args[2]
		timestamp := time.Now().Format("20060102150405")

		upFileName := fmt.Sprintf("src/database/migrator/migrations/%s_%s.up.sql", timestamp, name)
		downFileName := fmt.Sprintf("src/database/migrator/migrations/%s_%s.down.sql", timestamp, name)

		createFile(upFileName)
		createFile(downFileName)
	default:
		fmt.Println("Invalid migration command")
	}
}

func runMigrations(ctx context.Context, db *bun.DB) error {
	migrations := migrate.NewMigrations()

	if err := migrations.Discover(sqlMigrations); err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	err := migrator.Init(ctx)
	if err != nil {
		return err
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		log.Println("No new migrations to run")
	} else {
		log.Printf("Migrated to %s\n", group)
	}

	return nil
}

func RollbackMigrations(ctx context.Context, db *bun.DB) error {
	migrations := migrate.NewMigrations()

	if err := migrations.Discover(sqlMigrations); err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)

	group, err := migrator.Rollback(ctx)
	if err != nil {
		return err
	}

	log.Printf("Rolled back migration group %s\n", group)
	return nil
}

func createFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating migration file:", err)
		return
	}
	defer file.Close()

	fmt.Printf("Created migration file: %s\n", fileName)
}
