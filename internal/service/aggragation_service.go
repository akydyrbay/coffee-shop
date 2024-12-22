package service

import (
	"errors"
	"sort"

	"coffee/internal/dal"
	"coffee/models"
)

type AggragationService interface {
	GetTotalSales() (float64, error)
	GetPopularMenuItems() ([]models.OrderItem, error)
}

type aggragationService struct {
	orderRepo dal.OrderRepository
	menuRepo  dal.MenuRepository
}

func NewAggragationService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository) *aggragationService {
	return &aggragationService{orderRepo: orderRepo, menuRepo: menuRepo}
}

func (s *aggragationService) GetTotalSales() (float64, error) {
	allMenuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return 0, err
	}
	menuMap := make(map[string]models.MenuItem)
	for _, menuItem := range allMenuItems {
		menuMap[menuItem.ID] = menuItem
	}
	allOrderItems, err := s.orderRepo.GetAll()

	var totalSales float64
	for _, orderItem := range allOrderItems {
		if orderItem.Status == "open" {
			continue
		}
		for _, ingr := range orderItem.Items {
			itemMenu := menuMap[ingr.ProductID]
			if itemMenu.ID == "" {
				return 0, errors.New("menu item not found: " + orderItem.ID)
			} else {
				items := orderItem.Items
				for _, item := range items {
					totalSales += float64(item.Quantity) * itemMenu.Price
				}
			}
		}
	}
	return totalSales, nil
}

func (s *aggragationService) GetPopularMenuItems() ([]models.OrderItem, error) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}

	itemCount := make(map[string]int)

	for _, order := range orderItems {
		if order.Status == "closed" {
			for _, item := range order.Items {
				itemCount[item.ProductID] += item.Quantity
			}
		}
	}

	var popularItems []models.OrderItem
	for itemID, quantity := range itemCount {
		popularItems = append(popularItems, models.OrderItem{
			ProductID: itemID,
			Quantity:  quantity,
		})
	}

	sort.Slice(popularItems, func(i, j int) bool {
		return popularItems[i].Quantity > popularItems[j].Quantity
	})

	return popularItems, nil
}
