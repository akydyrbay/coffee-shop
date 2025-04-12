package service

import (
	"errors"
	"fmt"
	"frappuccino/internal/dal"
	"frappuccino/models"
	"time"
)

type OrderService interface {
	Create(orders []models.Order) models.Status
	Read(id int) (models.Order, models.Status)
	ReadAll() ([]models.Order, models.Status)
	UpdateStatus(orderId int) models.Status
	Update(newOrder models.Order) models.Status
	Delete(orderId int) models.Status
}

type orderService struct {
	orderRepo dal.OrderRepository
}

func NewOrderService(orderRepo dal.OrderRepository) *orderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) Read(id int) (models.Order, models.Status) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return models.Order{}, models.Status{ErrorMessage: err, Code: 500}
	}
	for _, orderItem := range orderItems {
		if orderItem.ID == id {
			return orderItem, models.Status{ErrorMessage: nil, Code: 200}
		}
	}
	return models.Order{}, models.NotFound
}

func (s *orderService) ReadAll() ([]models.Order, models.Status) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return []models.Order{}, models.Status{ErrorMessage: fmt.Errorf("failed to get order items: %w", err), Code: 500}
	}

	return orderItems, models.Success
}

func (s *orderService) UpdateStatus(id int) models.Status {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to get order items: %w", err), Code: 500}
	}

	for i := range orderItems {
		if orderItems[i].ID == id {
			if orderItems[i].Status == "closed" {
				return models.Status{ErrorMessage: errors.New("order is already closed"), Code: 400}
			}
			orderItems[i].Status = "closed"
			err = s.orderRepo.Update(orderItems[i])
			if err != nil {
				return models.Status{ErrorMessage: fmt.Errorf("failed to update order item: %w", err), Code: 500}
			}
			return models.Success
		}
	}
	return models.NotFound
}

func (s *orderService) Delete(orderId int) models.Status {
	exists, err := s.orderRepo.Exists(orderId)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
	}
	if !exists {
		return models.Status{ErrorMessage: errors.New("order item not found"), Code: 404}
	}

	if err := s.orderRepo.DeleteItem(orderId); err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to delete order item: %w", err), Code: 500}
	}
	return models.Status{Code: 204} // No content for successful deletion
}

func (s *orderService) Create(orders []models.Order) models.Status {
	for _, order := range orders {
		if !IsOrderValid(order) {
			return models.Status{ErrorMessage: errors.New("invalid order item"), Code: 400}
		}
		// pres, err := s.orderRepo.Exists(order.ID)
		// if err != nil {
		// 	return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
		// }
		// if pres {
		// 	return models.Status{ErrorMessage: errors.New("order item already exists"), Code: 400}
		// }
	}
	for _, order := range orders {
		order.CreatedAt = time.Now()
		if _, err := s.orderRepo.Insert(order); err != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to insert order item: %w", err), Code: 500}
		}
	}

	return models.Success
}

func (s *orderService) Update(newOrder models.Order) models.Status {
	if !IsOrderValid(newOrder) {
		return models.Status{ErrorMessage: errors.New("invalid order item"), Code: 400}
	}
	newOrder.CreatedAt = time.Now()
	if _, err := s.orderRepo.Exists(newOrder.ID); err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
	}
	if exist, _ := s.orderRepo.Exists(newOrder.ID); !exist {
		return models.NotFound
	}
	if err := s.orderRepo.Update(newOrder); err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to update order item: %w", err), Code: 500}
	}
	return models.Success
}
