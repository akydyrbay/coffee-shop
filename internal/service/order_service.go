package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"coffee/internal/dal"
	"coffee/models"
)

type OrderService interface {
	GetOrderItemById(id string) (models.Order, error)
	GetOrderItem() ([]models.Order, error)
	PostNewOrder(order models.Order) error
	UpdateOrderStatus(orderId string) error
	UpdateOrder(orderId string, newOrder models.Order) error
	DeleteOrder(orderId string) error
}

type orderService struct {
	orderRepo     dal.OrderRepository
	menuRepo      dal.MenuRepository
	inventoryRepo dal.InventoryRepository
}

func getFormattedTime() string {
	// Get current time in UTC
	currentTime := time.Now().UTC()

	// Format time in the desired format
	return currentTime.Format("2006-01-02T15:04:05Z")
}

func NewOrderService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) *orderService {
	return &orderService{orderRepo: orderRepo, menuRepo: menuRepo, inventoryRepo: inventoryRepo}
}

func (s *orderService) GetOrderItemById(id string) (models.Order, error) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return models.Order{}, err
	}
	for _, orderItem := range orderItems {
		if orderItem.ID == id {
			return orderItem, nil
		}
	}
	return models.Order{}, errors.New("inventory item not found")
}

func (s *orderService) GetOrderItem() ([]models.Order, error) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return []models.Order{}, err
	}

	return orderItems, nil
}

func (s *orderService) UpdateOrderStatus(id string) error {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for i := range orderItems {
		if orderItems[i].ID == id {
			if orderItems[i].Status == "closed" {
				return errors.New("order is already closed")
			}
			orderItems[i].Status = "closed"
			err = s.orderRepo.SaveAll(orderItems)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("order item not found")
}

func (s *orderService) DeleteOrder(orderId string) error {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for i := range orderItems {
		if orderItems[i].ID == orderId {
			if orderItems[i].Status == "closed" {
				return errors.New("order is already closed")
			}
			inventoryRepo, err := s.inventoryRepo.GetAll()
			if err != nil {
				return err
			}
			menus, err := s.menuRepo.GetAll()
			if err != nil {
				return err
			}
			fmt.Println(len(orderItems))
			newInvent, err := UpdateInventoryByOrder(inventoryRepo, orderItems[i], menus, false)
			if err != nil {
				return err
			}
			orderItems = append(orderItems[:i], orderItems[i+1:]...)
			err = s.inventoryRepo.SaveAll(newInvent)
			if err != nil {
				return err
			}
			s.orderRepo.SaveAll(orderItems)
			return nil
		}
	}

	return errors.New("order item not found")
}

func (s *orderService) PostNewOrder(order models.Order) error {
	if !IsOrderValid(order) {
		return errors.New("order is invalid")
	}
	err := IsValidOrder(order, s.menuRepo, s.inventoryRepo)
	if err != nil {
		return err
	}
	order.CreatedAt = getFormattedTime()
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}
	order.ID = strconv.Itoa(rand.Intn(99))
	for {
		pres, err := s.orderRepo.Exists(order.ID)
		if err != nil {
			return err
		}
		if pres {
			order.ID = strconv.Itoa(rand.Intn(99))
			continue
		}
		break
	}

	order.Status = "open"
	orderItems = append(orderItems, order)
	err = s.orderRepo.SaveAll(orderItems)
	if err != nil {
		return err
	}
	inventoryRepo, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}
	menus, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}
	newInvent, err := UpdateInventoryByOrder(inventoryRepo, order, menus, true)
	if err != nil {
		return err
	}
	err = s.inventoryRepo.SaveAll(newInvent)
	if err != nil {
		return err
	}
	return nil
}

func (s *orderService) UpdateOrder(orderId string, newOrder models.Order) error {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}
	if !IsOrderValid(newOrder) {
		return errors.New("invalid order")
	}
	for i := range orders {
		if orders[i].ID == orderId {
			if orders[i].Status == "closed" {
				return errors.New("order is already closed")
			}
			err := IsValidOrder(newOrder, s.menuRepo, s.inventoryRepo)
			if err != nil {
				return err
			}
			orders[i].Items = append(orders[i].Items, newOrder.Items...)
			inventoryRepo, err := s.inventoryRepo.GetAll()
			if err != nil {
				return err
			}
			menus, err := s.menuRepo.GetAll()
			if err != nil {
				return err
			}
			newInvent, err := UpdateInventoryByOrder(inventoryRepo, orders[i], menus, true)
			if err != nil {
				return err
			}
			err = s.orderRepo.SaveAll(orders)
			if err != nil {
				return err
			}
			err = s.inventoryRepo.SaveAll(newInvent)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("order item not found")
}
