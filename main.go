package main

import (
	"database/sql"
	"net/http"

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

	http.ListenAndServe(":8000", r)
}
