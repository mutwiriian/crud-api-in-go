package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mutwiriian/crud-api-in-go/database"
	"github.com/mutwiriian/crud-api-in-go/models"
)

func CreateCustomerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customerPayload models.CreateCustomerSchema

		if err := json.NewDecoder(r.Body).Decode(&customerPayload); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)

			response := map[string]any{
				"status":  "Fail",
				"message": err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		if customerPayload.Name == "" || customerPayload.Email == "" || customerPayload.Phone_number == "" || customerPayload.Address == "" {
			w.Header().Set("Content-Type", "application/json")
			msg := "All fields must be provided to create customer"
			http.Error(w, msg, http.StatusBadRequest)
			database.Logger.Error(msg)
			return
		}

		insertStmt, err := db.Prepare("insert into customers (name, email, phone_number, address) values ($1, $2, $3, $4)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		defer insertStmt.Close()

		_, err = insertStmt.Exec(customerPayload.Name, customerPayload.Email, customerPayload.Phone_number, customerPayload.Address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)

		response := map[string]any{
			"status":  "Success",
			"message": "Customer successfully added!",
		}
		json.NewEncoder(w).Encode(response)
		database.Logger.Info("msg", response["message"], "method", "POST", "path", "/customers/create")
	}
}

func GetCustomersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getCustomersStmt, err := db.Prepare("select * from customers")
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			database.Logger.Error(err.Error())
		}

		defer getCustomersStmt.Close()

		rows, err := getCustomersStmt.Query()
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			database.Logger.Error(err.Error())
			return
		}

		defer rows.Close()

		var customers []models.Customer

		for rows.Next() {
			var customer models.Customer

			err := rows.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
			if err != nil {
				http.Error(w, err.Error(), http.StatusConflict)
				database.Logger.Error(err.Error())
				return
			}
			customers = append(customers, customer)
		}

		err = rows.Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"status": "Success",
			"data":   customers,
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		defer searchStmt.Close()

		row := searchStmt.QueryRow(customerEmail)
		var customer models.Customer

		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Address, &customer.Address)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusBadGateway)
			database.Logger.Error(sql.ErrNoRows.Error())
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())

			return
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

		searchStmt, err := db.Prepare("select * from customers where email = $1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}
		var customer models.Customer

		row := searchStmt.QueryRow(customerEmail)
		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			database.Logger.Error(err.Error())
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		var customerPayload models.UpdateCustomerSchema

		err = json.NewDecoder(r.Body).Decode(&customerPayload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		updateNameStmt, err := db.Prepare("update customers set name =$1 where email=$2")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		updatePhoneNumberStmt, err := db.Prepare("update customers set phone_number =$1 where email=$2")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		updateAddressStmt, err := db.Prepare("update customers set address =$1 where email=$2")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		if customerPayload.Name != "" {
			_, err := updateNameStmt.Exec(customerPayload.Name, customerEmail)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				database.Logger.Error(err.Error())
				return
			}
		}

		if customerPayload.Phone_number != "" {
			_, err := updatePhoneNumberStmt.Exec(customerPayload.Phone_number, customerEmail)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				database.Logger.Error(err.Error())
				return
			}
		}

		if customerPayload.Address != "" {
			_, err := updateAddressStmt.Exec(customerPayload.Address, customerEmail)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				database.Logger.Error(err.Error())
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"status":  "Success",
			"message": "Customer updates successful",
		}

		json.NewEncoder(w).Encode(response)
		database.Logger.Info("Customer updates successful", "method", "POST", "path", "/customers/update")
	}
}

func DeleteCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")
		searchStmt, err := db.Prepare("select * from customers where email = $1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			database.Logger.Error(err.Error())
			return
		}

		var customer models.Customer

		row := searchStmt.QueryRow(customerEmail)
		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		deleteStmt, err := db.Prepare("delete from customers where email = $1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		res, err := deleteStmt.Exec(customerEmail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		affected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			database.Logger.Error(err.Error())
			return
		}

		if affected == 0 {
			msg := "Now rows affected!"
			http.Error(w, msg, http.StatusNotFound)
			database.Logger.Info("Delete operation completed but no rows affected", "method", "DELETE", "path", "/customers/delete_email")

			w.Header().Set("Content-Type", "application/json")
			response := map[string]any{
				"status":  "success",
				"message": "Customer successfully deleted",
			}

			json.NewEncoder(w).Encode(response)
			database.Logger.Info("Customer successfully deleted", "method", "DELETE", "path", "/customers/delete_email")
		}
	}
}
