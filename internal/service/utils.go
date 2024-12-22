package service

import (
	"errors"

	"coffee/internal/dal"
	"coffee/models"
)

func IsValidOrder(order models.Order, menuRepo dal.MenuRepository, inventRepo dal.InventoryRepository) error {
	items := order.Items
	menuItems, err := menuRepo.GetAll()
	if err != nil {
		return err
	}
	inventIngredients, err := inventRepo.GetAll()
	if err != nil {
		return err
	}

	inventMap := make(map[string]models.InventoryItem)
	for _, inventIngr := range inventIngredients {
		inventMap[inventIngr.IngredientID] = inventIngr
	}

	for _, item := range items {
		foundMenuItem := false

		for _, menuItem := range menuItems {
			if item.ProductID == menuItem.ID {
				foundMenuItem = true
				quantity := item.Quantity

				for _, ingredient := range menuItem.Ingredients {
					inventIngr, foundIngredient := inventMap[ingredient.IngredientID]

					if !foundIngredient {
						return errors.New("invalid ingredient: " + ingredient.IngredientID)
					}
					if ingredient.Quantity*float64(quantity) > inventIngr.Quantity {
						return errors.New("not enough ingredient: " + ingredient.IngredientID)
					}
				}
				break
			}
		}
		if !foundMenuItem {
			return errors.New("order item doesn't exist in menu")
		}
	}

	return nil
}

func UpdateInventoryByOrder(inventory []models.InventoryItem, order models.Order, menuItems []models.MenuItem, subtract bool) ([]models.InventoryItem, error) {
	// Map inventory items for quick lookup
	inventoryMap := make(map[string]*models.InventoryItem)
	for i := range inventory {
		inventoryMap[inventory[i].IngredientID] = &inventory[i]
	}

	// Map menu items for quick lookup
	menuMap := make(map[string]models.MenuItem)
	for _, menuItem := range menuItems {
		menuMap[menuItem.ID] = menuItem
	}

	// Iterate through order items
	for _, orderItem := range order.Items {
		menuItem, exists := menuMap[orderItem.ProductID]
		if !exists {
			return nil, errors.New("menu item not found: " + orderItem.ProductID)
		}

		// Calculate ingredient usage based on order quantity
		for _, ingredient := range menuItem.Ingredients {
			inventoryItem, found := inventoryMap[ingredient.IngredientID]
			if !found {
				return nil, errors.New("ingredient not found in inventory: " + ingredient.IngredientID)
			}

			// Calculate the adjustment amount
			adjustment := ingredient.Quantity * float64(orderItem.Quantity)
			if subtract {
				if inventoryItem.Quantity < adjustment {
					return nil, errors.New("not enough of ingredient: " + ingredient.IngredientID)
				}
				inventoryItem.Quantity -= adjustment
			} else {
				inventoryItem.Quantity += adjustment
			}
		}
	}

	// Save updated inventory back to the file
	return inventory, nil
}

func IsMenuValid(item models.MenuItem) bool {
	if item.Price <= 0 {
		return false
	}
	if item.Description == "" || item.Name == "" {
		return false
	}
	for _, item := range item.Ingredients {
		if item.IngredientID == "" || item.Quantity <= 0 {
			return false
		}
	}
	return true
}

func IsOrderValid(order models.Order) bool {
	if order.Status == "closed" {
		return false
	}
	if order.CustomerName == "" {
		return false
	}
	for _, item := range order.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return false
		}
	}
	return true
}

func IsInventoryValid(inventory models.InventoryItem) bool {
	if inventory.IngredientID == "" || inventory.Quantity <= 0 {
		return false
	}
	return true
}
