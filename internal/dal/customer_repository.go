package dal

import (
	"database/sql"
	"frappuccino/models"
)

type CustomerRepository interface {
	SaveCustomer(customer models.Customer) error
	IsEmailUnique(email string) (bool, error)
	GetAllCustomers() ([]models.Customer, error)
	GetCustomerByID(id int) (models.Customer, error)
	DeleteCustomer(id int) error
	IsCustomerExist(id int) (bool, error)
	GetCustomerByName(name string) (models.Customer, error)
}

type customerRepo struct {
	DB *sql.DB
}

func NewCustomerRepo(db *sql.DB) *customerRepo {
	return &customerRepo{DB: db}
}

func (r *customerRepo) GetCustomerByName(name string) (models.Customer, error) {
	var cust models.Customer
	err := r.DB.QueryRow(`
        SELECT customer_id, name, email, number 
        FROM customers 
        WHERE name = $1 
        LIMIT 1
    `, name).Scan(&cust.Customer_id, &cust.Name, &cust.Email, &cust.Number)
	if err != nil {
		return models.Customer{}, err
	}
	return cust, nil
}

func (r *customerRepo) SaveCustomer(customer models.Customer) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO customers(name, email, number)
	VALUES($1,$2,$3)`, customer.Name, customer.Email, customer.Number)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *customerRepo) IsEmailUnique(email string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(*)
	FROM customers
	WHERE email=$1
	`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *customerRepo) GetAllCustomers() ([]models.Customer, error) {
	rows, err := r.DB.Query("SELECT customer_id, name, email, number FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Number); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *customerRepo) GetCustomerByID(id int) (models.Customer, error) {
	var customer models.Customer
	err := r.DB.QueryRow("SELECT customer_id, name, email, number FROM customers WHERE customer_id=$1", id).
		Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Number)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (r *customerRepo) DeleteCustomer(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`Delete FROM customers
	WHERE customer_id=$1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *customerRepo) IsCustomerExist(id int) (bool, error) {
	inventoryData, err := r.GetAllCustomers()
	if err != nil {
		return false, err
	}

	for _, inventory := range inventoryData {
		if inventory.Customer_id == id {
			return true, nil
		}
	}
	return false, nil
}
