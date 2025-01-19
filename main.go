package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/mutwiriian/crud-api-in-go/database"
	"github.com/mutwiriian/crud-api-in-go/handlers"
)

var db *sql.DB

func init() {
	db = database.ConnectDB()
	database.CreateCustomersTable(db)
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /customers/get_all", handlers.GetCustomersHandler(db))
	r.HandleFunc("POST /customers/create", handlers.CreateCustomerHandler(db))
	r.HandleFunc("GET /customers/get_email", handlers.GetCustomerByEmailHandler(db))
	r.HandleFunc("PATCH /customers/update", handlers.UpdateCustomerByEmailHandler(db))
	r.HandleFunc("DELETE /customers/delete_email", handlers.DeleteCustomerByEmailHandler(db))

	database.Logger.Info("Starting server at", slog.String("port", "8000"))

	err := http.ListenAndServe(":8000", r)

	database.Logger.Error(err.Error())
	os.Exit(1)
}
