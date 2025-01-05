package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() *sql.DB {
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	db := os.Getenv("DBNAME")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")

	if host == "" || port == "" || db == "" || user == "" || password == "" {
		log.Fatal("All environment variables(DBHOST, DBPORT, DB, DBUSER, DBPASS) must be provided")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db)
	DB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database! %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Unable to connect to database! %v", err)
	}

	log.Println("Connected to database!")

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
		log.Fatalf("Database migration failed: %v", err)
		return
	}

	log.Println("Database migration successful!")
}
