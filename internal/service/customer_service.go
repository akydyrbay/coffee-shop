package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/models"
	"net/http"
)

type CustomerService interface {
	CreateCustomer(customer models.Customer) models.Status
	GetAllCustomers(w http.ResponseWriter) ([]models.Customer, models.Status)
	GetCustomer(w http.ResponseWriter, id int) (models.Customer, models.Status)
	DeleteCustomer(id int) models.Status
}

type customerService struct {
	repo dal.CustomerRepository
}

func NewServiceCustomer(r dal.CustomerRepository) *customerService {
	return &customerService{repo: r}
}

func (s *customerService) CreateCustomer(customer models.Customer) models.Status {
	unique, err := s.repo.IsEmailUnique(customer.Email)
	if err != nil {
		return models.Status{ErrorMessage: errors.New("customer cannot check if unique"), Code: 400}
	}
	b := IsCustomerValid(customer)
	if !b {
		return models.Status{ErrorMessage: errors.New("incorrect customer"), Code: 400}
	}
	if !unique {
		return models.Status{ErrorMessage: errors.New("customer not unique"), Code: 400}
	}
	err = s.repo.SaveCustomer(customer)
	if err != nil {
		return models.Status{ErrorMessage: errors.New("customer not saved"), Code: 400}
	}
	return models.Status{ErrorMessage: nil, Code: 200}
}

func (s *customerService) GetAllCustomers(w http.ResponseWriter) ([]models.Customer, models.Status) {
	customers, err := s.repo.GetAllCustomers()
	if err != nil {
		return []models.Customer{}, models.Status{ErrorMessage: errors.New("customer not got"), Code: 400}
	}
	return customers, models.Status{ErrorMessage: nil, Code: 200}
}

func (s *customerService) GetCustomer(w http.ResponseWriter, id int) (models.Customer, models.Status) {
	if _, err := s.repo.IsCustomerExist(id); err != nil {
		return models.Customer{}, models.Status{ErrorMessage: errors.New("customer not exists"), Code: 400}
	}
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return models.Customer{}, models.Status{ErrorMessage: errors.New("customer not found"), Code: 400}
	}
	return customer, models.Status{ErrorMessage: nil, Code: 200}
}

func (s *customerService) DeleteCustomer(id int) models.Status {
	if _, err := s.repo.IsCustomerExist(id); err != nil {
		return models.Status{ErrorMessage: errors.New("customer not found"), Code: 400}
	}
	err := s.repo.DeleteCustomer(id)
	if err != nil {
		return models.Status{ErrorMessage: errors.New("customer not deleted"), Code: 404004}
	}
	return models.Status{ErrorMessage: nil, Code: 200}
}
