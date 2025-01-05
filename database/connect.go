package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() *sql.DB {
	dsn := "host=172.22.0.1 port=5432 user=postgres password=postgres dbname=customers_crud sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database! %v", err)
	}

	log.Println("Connected to database!")

	return db
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
