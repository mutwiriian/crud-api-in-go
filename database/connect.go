package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	host     = flag.String("host", "127.0.0.1", "Database host")
	port     = flag.String("port", "5432", "Port number not in 0-1023")
	database = flag.String("database", "postgres", "Database name")
	username = flag.String("username", "postgres", "Database user name")
	password = flag.String("password", "postgres", "User password")
)

var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

func ConnectDB() *sql.DB {
	flag.Parse()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", *host, *port, *username, *password, *database)
	DB, err := sql.Open("pgx", dsn)
	if err != nil {
		Logger.Error(err.Error())
		os.Exit(1)
	}

	if err := DB.Ping(); err != nil {
		Logger.Error(err.Error())
		os.Exit(1)
	}

	Logger.Info("Connected to database!")

	return DB
}

func CreateCustomersTable(db *sql.DB) {
	stmt := `
  create table if not exists customers (
    id serial primary key,
    name varchar(50) not null,
    email varchar(50) unique not null,
    phone_number varchar(10) not null,
    address text not null
  );
  `
	_, err := db.Exec(stmt)
	if err != nil {
		Logger.Error(err.Error())
		os.Exit(1)
	}

	Logger.Info("Database migration successful!")
}
