package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mutwiriian/crud-api-in-go/database"
	"github.com/mutwiriian/crud-api-in-go/models"
)

func internalServerErrorHandler(w http.ResponseWriter, err error, code int) {
	if err != nil {
		http.Error(w, err.Error(), code)
		database.Logger.Error(err.Error())
		return
	}
}

func CreateCustomerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customerPayload models.CreateCustomerSchema

		err := json.NewDecoder(r.Body).Decode(&customerPayload)
		if err != nil {
			internalServerErrorHandler(w, err, http.StatusInternalServerError)
		}

		if customerPayload.Name == "" || customerPayload.Email == "" || customerPayload.Phone_number == "" || customerPayload.Address == "" {
			w.Header().Set("Content-Type", "application/json")
			msg := "All fields must be provided to create customer"
			http.Error(w, msg, http.StatusBadRequest)
			database.Logger.Error(msg)
			return
		}

		insertStmt, err := db.Prepare("insert into customers (name, email, phone_number, address) values ($1, $2, $3, $4)")
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		defer insertStmt.Close()

		_, err = insertStmt.Exec(customerPayload.Name, customerPayload.Email, customerPayload.Phone_number, customerPayload.Address)
		if err != nil {
			internalServerErrorHandler(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		msg := "Customer successfully added!"
		json.NewEncoder(w).Encode(msg)
		database.Logger.Info(msg, "method", "POST", "path", "/customers/create")
	}
}

func GetCustomersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getCustomersStmt, err := db.Prepare("select * from customers")
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		defer getCustomersStmt.Close()

		rows, err := getCustomersStmt.Query()
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		defer rows.Close()

		var customers []models.Customer

		for rows.Next() {
			var customer models.Customer

			err := rows.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			customers = append(customers, customer)
		}

		err = rows.Err()
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"customers": customers,
		}
		json.NewEncoder(w).Encode(response)
		database.Logger.Info("Customers returned", "method", "GET", "path", "/customers/get_all")
	}
}

func GetCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")

		if customerEmail == "" {
			http.Error(w, "Enter valid customer email.", http.StatusBadRequest)
			database.Logger.Error("Customer email not provided")
			return
		}

		searchStmt, err := db.Prepare("select * from customers where email = $1")
		if err != nil {
			internalServerErrorHandler(w, err, http.StatusInternalServerError)
			return
		}

		defer searchStmt.Close()

		row := searchStmt.QueryRow(customerEmail)

		var customer models.Customer

		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Address, &customer.Address)
		if err != nil {
			internalServerErrorHandler(w, err, http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"status": "success",
			"data":   customer,
		}
		json.NewEncoder(w).Encode(response)
		database.Logger.Info("Customer returned", "method", "GET", "path", "/customers/get_email")
	}
}

func UpdateCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")

		var customerPayload models.UpdateCustomerSchema

		err := json.NewDecoder(r.Body).Decode(&customerPayload)
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		searchStmt, err := db.Prepare("select * from customers where email = $1")
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		defer searchStmt.Close()

		var customer models.Customer

		row := searchStmt.QueryRow(customerEmail)
		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		if customerPayload.Name == "" || customerPayload.Phone_number == "" || customerPayload.Address == "" {
			internalServerErrorHandler(w, err, http.StatusInternalServerError)
		}

		tx, err := db.Begin()
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		if customerPayload.Name != "" {
			updateNameStmt, err := tx.Prepare("update customers set name =$1 where email=$2")
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			_, err = updateNameStmt.Exec(customerPayload.Name, customerEmail)
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			defer updateNameStmt.Close()
		}

		if customerPayload.Phone_number != "" {
			updatePhoneNumberStmt, err := tx.Prepare("update customers set phone_number =$1 where email=$2")
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			_, err = updatePhoneNumberStmt.Exec(customerPayload.Phone_number, customerEmail)
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			defer updatePhoneNumberStmt.Close()
		}

		if customerPayload.Address != "" {
			updateAddressStmt, err := tx.Prepare("update customers set address =$1 where email=$2")
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			_, err = updateAddressStmt.Exec(customerPayload.Address, customerEmail)
			internalServerErrorHandler(w, err, http.StatusInternalServerError)

			defer updateAddressStmt.Close()
		}

		if err := tx.Commit(); err != nil {
			internalServerErrorHandler(w, err, http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"status":  "Success",
			"message": "Customer updates successful",
		}

		json.NewEncoder(w).Encode(response)
		database.Logger.Info(response["message"], "method", "POST", "path", "/customers/update")
	}
}

func DeleteCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")

		searchStmt, err := db.Prepare("select * from customers where email = $1")
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		var customer models.Customer

		row := searchStmt.QueryRow(customerEmail)
		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		deleteStmt, err := db.Prepare("delete from customers where email = $1")
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		res, err := deleteStmt.Exec(customerEmail)
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		affected, err := res.RowsAffected()
		internalServerErrorHandler(w, err, http.StatusInternalServerError)

		if affected == 0 {
			msg := "Now rows affected!"
			http.Error(w, msg, http.StatusNotFound)
			database.Logger.Info("Delete operation completed but no rows affected", "method", "DELETE", "path", "/customers/delete_email")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"status":  "success",
			"message": "Customer successfully deleted",
		}

		json.NewEncoder(w).Encode(response)
		database.Logger.Info(response["message"], "method", "DELETE", "path", "/customers/delete_email")
	}
}
