package models

type Customer struct {
	Customer_id  int
	Name         string
	Email        string
	Phone_number int
	Address      string
}

type CreateCustomerSchema struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone_number string `json:"phone_number"`
	Address      string `json:"address"`
}

type UpdateCustomerSchema struct {
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone_number string `json:"phone_number,omitempty"`
	Address      string `json:"address,omitempty"`
}
