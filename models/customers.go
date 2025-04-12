package models

type Customer struct {
	Customer_id int    `json:"customer_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Number      string `json:"number"`
}
