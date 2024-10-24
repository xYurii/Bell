package database

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/xYurii/Bell/src/database/adapter"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var Database *bun.DB
var User adapter.UserAdapter
var Guild adapter.GuildAdapter
var Raffle adapter.RaffleAdapter

func GetEnvDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}
}

func InitDatabase(config *DatabaseConfig) (*bun.DB, error) {
	db := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", config.Host, config.Port)),
		pgdriver.WithUser(config.User),
		pgdriver.WithPassword(config.Password),
		pgdriver.WithDatabase(config.Database),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(db)
	err := sqldb.Ping()
	if err != nil {
		return nil, err
	}

	Database = bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	Database.SetMaxOpenConns(maxOpenConns)
	Database.SetMaxIdleConns(maxOpenConns)

	User = adapter.NewUserAdapter(Database)
	Guild = adapter.NewGuildAdapter(Database)
	Raffle = adapter.NewRaffleAdapter(Database)

	return Database, nil
}
