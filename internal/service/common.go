package service

import (
	"errors"
	"fmt"
	"frappuccino/internal/dal"
	"frappuccino/models"
	"strconv"
	"time"
)

type Service interface {
	CloseOrdersByID(id int) models.Status
	ProcessBatchOrders(request models.BatchProcessRequest) (models.BatchProcessResponse, error)
}

type service struct {
	orderRepo     dal.OrderRepository
	inventoryRepo dal.InventoryRepository
	menuRepo      dal.MenuRepository
	customerRepo  dal.CustomerRepository
}

func NewService(orderRepo dal.OrderRepository, inventoryRepo dal.InventoryRepository, menuRepo dal.MenuRepository, custRepo dal.CustomerRepository) *service {
	return &service{orderRepo: orderRepo, inventoryRepo: inventoryRepo, menuRepo: menuRepo, customerRepo: custRepo}
}

func (s *service) CloseOrdersByID(orderId int) models.Status {
	// Check if order exists - do this only once
	exist, err := s.orderRepo.Exists(orderId)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check if order exists: %w", err), Code: 500}
	}
	if !exist {
		return models.Status{ErrorMessage: errors.New("order does not exist"), Code: 404}
	}

	orderService := &orderService{orderRepo: s.orderRepo}
	order, status := orderService.Read(orderId)
	if status.ErrorMessage != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to read order: %w", status.ErrorMessage), Code: 500}
	}

	if order.Status == "closed" {
		return models.Status{ErrorMessage: errors.New("order is already closed"), Code: 400}
	}

	// Create services once, outside the loops
	ms := &menuService{menuRepo: s.menuRepo}
	ins := &inventoryService{inventoryRepo: s.inventoryRepo}

	// Collect all inventory updates so we can apply them at once
	type inventoryUpdate struct {
		ID       int
		Quantity float64
	}

	inventoryUpdates := make(map[int]float64)

	// First pass: check if all inventory operations would succeed
	for _, orderItem := range order.Items {
		menuItem, status := ms.Read(orderItem.MenuItemID)
		if status.ErrorMessage != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to read menu item: %w", status.ErrorMessage), Code: 500}
		}

		for _, ingredient := range menuItem.Ingredients {
			inventoryItem, status := ins.Read(ingredient.IngredientID)
			if status.ErrorMessage != nil {
				return models.Status{ErrorMessage: fmt.Errorf("failed to read inventory item: %w", status.ErrorMessage), Code: 500}
			}

			// Calculate how much will be consumed
			quantityToConsume := ingredient.Quantity * float64(orderItem.Quantity)

			// Track the total consumption for this inventory item
			inventoryUpdates[ingredient.IngredientID] += quantityToConsume

			// Check if we'll have enough inventory
			// slog.Info("Check quantity to consume:", quantityToConsume, "inventory item:", inventoryItem, "inventory quantity:", inventoryItem.Quantity, "left", inventoryItem.Quantity-inventoryUpdates[ingredient.IngredientID])
			if inventoryItem.Quantity-inventoryUpdates[ingredient.IngredientID] < 0 {
				menuItem.IsAvailable = false
				return models.Status{ErrorMessage: fmt.Errorf("insufficient inventory for item %d: need %.2f but only %.2f available", ingredient.IngredientID, inventoryUpdates[ingredient.IngredientID], inventoryItem.Quantity), Code: 400}
			}
		}
	}

	// Reduce
	for id, quantity := range inventoryUpdates {
		inventoryItem, status := ins.Read(id)
		if status.ErrorMessage != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to read inventory item: %w", status.ErrorMessage), Code: 500}
		}
		inventoryItem.Quantity -= quantity
		status = ins.Update(inventoryItem)
		if status.ErrorMessage != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to update inventory item: %w", status.ErrorMessage), Code: 500}
		}
	}

	order.Status = "closed"
	err = s.orderRepo.Update(order)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to update order: %w", err), Code: 500}
	}

	return models.Success
}

func (s *service) ProcessBatchOrders(request models.BatchProcessRequest) (models.BatchProcessResponse, error) {
	var response models.BatchProcessResponse
	var totalRevenue float64
	var accepted, rejected int
	var inventoryUpdates []models.InventoryUpdate

	// Start a transaction
	tx, err := s.orderRepo.BeginTx()
	if err != nil {
		return response, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, order := range request.Orders {
		status := "accepted"
		var reason string

		if len(order.Items) == 0 {
			status = "rejected"
			reason = "invalid order: order items is empty"
		}

		for _, item := range order.Items {
			exists, err := s.menuRepo.Exists(item.MenuItemID)
			if err != nil || !exists {
				status = "rejected"
				reason = "ordered item does not exist: " + strconv.Itoa(item.MenuItemID)
			}
		}

		needInventory, err := s.orderRepo.GetNeedInventory(order)
		if err != nil {
			status = "rejected"
			reason = "invalid order"
		} else {
			enough, err := s.orderRepo.IsEnough(needInventory)
			if err != nil || !enough {
				status = "rejected"
				reason = "insufficient inventory"
			}
		}

		if exist, _ := s.customerRepo.IsCustomerExist(order.CustomerID); !exist {
			status = "rejected"
			reason = "customer id does not exist"
		}

		if status == "accepted" {
			order.CreatedAt = time.Now()
			order.Status = "active"

			if _, err := s.orderRepo.CheckOrder(order); err != nil {
				tx.Rollback()
				return response, fmt.Errorf("check order failed: %w", err)
			}
			if err := s.orderRepo.GetPriceAtOrderItems(&order); err != nil {
				tx.Rollback()
				return response, err
			}
			if err := s.orderRepo.GetTotalAmount(&order); err != nil {
				tx.Rollback()
				return response, err
			}
			if _, err := s.orderRepo.Insert(order); err != nil {
				tx.Rollback()
				return response, err
			}
			order.ID, err = s.orderRepo.GetLastOrderId()
			if err != nil {
				tx.Rollback()
				return response, err
			}
			if err := s.orderRepo.CloseOrder(order.ID); err != nil {
				tx.Rollback()
				return response, err
			}

			totalRevenue += order.TotalAmount
			for ingredientID, quantityUsed := range needInventory {
				num, err := strconv.Atoi(ingredientID)
				if err != nil {
					tx.Rollback()
					return response, err
				}
				inventory, err := s.inventoryRepo.GetInventory(num)
				if err != nil {
					tx.Rollback()
					return response, err
				}
				inventoryUpdates = append(inventoryUpdates, models.InventoryUpdate{
					IngredientID: ingredientID,
					QuantityUsed: quantityUsed,
					Remaining:    inventory.Quantity,
				})
			}
			accepted++
		} else {
			rejected++
		}

		response.ProcessedOrders = append(response.ProcessedOrders, models.ProcessedOrder{
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Status:     status,
			Total:      order.TotalAmount,
			Reason:     reason,
		})
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return response, err
	}

	response.Summary = models.BatchProcessSummary{
		TotalOrders:      len(request.Orders),
		Accepted:         accepted,
		Rejected:         rejected,
		TotalRevenue:     totalRevenue,
		InventoryUpdates: inventoryUpdates,
	}

	return response, nil
}
