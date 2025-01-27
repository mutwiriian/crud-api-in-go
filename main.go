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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	app := &handlers.Application{
		Logger: logger,
	}

	r := http.NewServeMux()

	r.HandleFunc("GET /customers/get_all", app.GetCustomersHandler(db))
	r.HandleFunc("POST /customers/create", app.CreateCustomerHandler(db))
	r.HandleFunc("GET /customers/get_email", app.GetCustomerByEmailHandler(db))
	r.HandleFunc("PATCH /customers/update", app.UpdateCustomerByEmailHandler(db))
	r.HandleFunc("DELETE /customers/delete_email", app.DeleteCustomerByEmailHandler(db))

	app.Logger.Info("Starting server at", slog.String("port", "8000"))

	err := http.ListenAndServe(":8000", r)

	app.Logger.Error(err.Error())
	os.Exit(1)
}
