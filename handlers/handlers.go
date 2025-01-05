package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mutwiriian/crud-api-in-go/models"
)

func CreateCustomerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customerPayload models.CreateCustomerSchema

		if err := json.NewDecoder(r.Body).Decode(&customerPayload); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)

			response := map[string]any{
				"status":  "fail",
				"message": err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		if customerPayload.Name == "" || customerPayload.Email == "" || customerPayload.Phone_number == "" || customerPayload.Address == "" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "All fields must be provided to create customer", http.StatusBadRequest)
			log.Println("All field must be provided to create customer")
			return
		}

		insertStmt, err := db.Prepare("insert into customers (name, email, phone_number, address) values ($1, $2, $3, $4)")
		if err != nil {
			log.Printf("Failed to create SQL statement: %v", err)
			http.Error(w, "Failed to create SQL statement", http.StatusInternalServerError)
			return
		}

		defer insertStmt.Close()

		_, err = insertStmt.Exec(customerPayload.Name, customerPayload.Email, customerPayload.Phone_number, customerPayload.Address)
		if err != nil {
			log.Printf("Customer creation failed: %v", err)
			http.Error(w, "Customer creation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)

		response := map[string]any{
			"status":  "success",
			"message": "Customer successfully added!",
		}
		json.NewEncoder(w).Encode(response)
		log.Println("Customer successfully added!")
	}
}

func GetCustomersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getCustomersStmt, err := db.Prepare("select * from customers")
		if err != nil {
			log.Printf("Failed to create SQL statement: %v", err)
			http.Error(w, "Failed to create SQL statement", http.StatusConflict)
		}

		defer getCustomersStmt.Close()

		rows, err := getCustomersStmt.Query()
		if err != nil {
			log.Printf("Failed to execute statement: %v", err)
			http.Error(w, "Failed to execute statement", http.StatusConflict)
			return
		}

		defer rows.Close()

		var customers []models.Customer

		for rows.Next() {
			var customer models.Customer

			err := rows.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
			if err != nil {
				log.Printf("Failed to load customer: %v", err)
				http.Error(w, "Failed to load customer", http.StatusConflict)
				return
			}
			customers = append(customers, customer)
		}

		err = rows.Err()
		if err != nil {
			log.Printf("An error occurred while fetching customers: %v", err)
			http.Error(w, "An error occurred while fetching customers", http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"status": "success",
			"data":   customers,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func GetCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")
		if customerEmail == "" {
			log.Println("Enter valid customer email")
			http.Error(w, "Enter valid customer email.", http.StatusBadRequest)
			return
		}

		searchStmt, err := db.Prepare("select * from customers where email = $1")
		if err != nil {
			log.Printf("Failed to define SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		defer searchStmt.Close()

		row := searchStmt.QueryRow(customerEmail)
		var customer models.Customer

		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Address, &customer.Address)
		if err == sql.ErrNoRows {
			log.Printf("No customer with given email exists: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		} else if err != nil {
			log.Printf("An error occurred while finding customer!: %v", err)
			http.Error(w, "An error occurred while finding customer", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"status": "success",
			"data":   customer,
		}
		json.NewEncoder(w).Encode(response)
		log.Println("Customer returned")
	}
}

func UpdateCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")

		searchStmt, err := db.Prepare("select * from customers where email = $1")
		if err != nil {
			log.Printf("Failed to prepare SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var customer models.Customer

		row := searchStmt.QueryRow(customerEmail)
		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
		if err == sql.ErrNoRows {
			log.Printf("Customer with given email not found: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if err != nil {
			log.Printf("An error occured while fetching customer: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var customerPayload models.UpdateCustomerSchema

		err = json.NewDecoder(r.Body).Decode(&customerPayload)
		if err != nil {
			log.Printf("Failed to decode update body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updateNameStmt, err := db.Prepare("update customers set name =$1 where email=$2")
		if err != nil {
			log.Printf("Failed to prepare SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updatePhoneNumberStmt, err := db.Prepare("update customers set phone_number =$1 where email=$2")
		if err != nil {
			log.Printf("Failed to prepare SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updateAddressStmt, err := db.Prepare("update customers set address =$1 where email=$2")
		if err != nil {
			log.Printf("Failed to prepare SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if customerPayload.Name != "" {
			_, err := updateNameStmt.Exec(customerPayload.Name, customerEmail)
			if err != nil {
				log.Printf("Failed to update %s: %v", customerPayload.Name, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if customerPayload.Phone_number != "" {
			_, err := updatePhoneNumberStmt.Exec(customerPayload.Phone_number, customerEmail)
			if err != nil {
				log.Printf("Failed to update %s: %v", customerPayload.Phone_number, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if customerPayload.Address != "" {
			_, err := updateAddressStmt.Exec(customerPayload.Address, customerEmail)
			if err != nil {
				log.Printf("Failed to update %s: %v", customerPayload.Address, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]any{
			"status":  "success",
			"message": "Customer updates successful",
		}

		json.NewEncoder(w).Encode(response)
		log.Println("Customer updates successful")
	}
}

func DeleteCustomerByEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerEmail := r.URL.Query().Get("email")
		searchStmt, err := db.Prepare("select * from customers where email = $1")
		if err != nil {
			log.Printf("Failed to define SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		var customer models.Customer

		row := searchStmt.QueryRow(customerEmail)
		err = row.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Phone_number, &customer.Address)
		if err == sql.ErrNoRows {
			log.Println("No customer with given email found!")
			http.Error(w, "No customer with given email found!", http.StatusInternalServerError)
			return
		} else if err != nil {
			log.Printf("An error occurred while fetching customer: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deleteStmt, err := db.Prepare("delete from customers where email = $1")
		if err != nil {
			log.Printf("Failed to define SQL statement: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		res, err := deleteStmt.Exec(customerEmail)
		if err != nil {
			log.Printf("An error occured: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			log.Printf("An error occured: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if affected == 0 {
			log.Printf("No rows affected by delete operation: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		response := map[string]any{
			"status":  "success",
			"message": "User successfully deleted",
		}

		json.NewEncoder(w).Encode(response)
		log.Println("User successfully deleted")
	}
}
