package service

import (
	"frappuccino/internal/dal"
	"frappuccino/models"
	"time"
)

func GetFormattedTime() string {
	// Get current time in UTC
	currentTime := time.Now().UTC()

	// Format time in the desired format
	return currentTime.Format("2006-01-02T15:04:05Z")
}

type Err struct {
	Err_Message string `json:"Error"`
	Status      int    `json:"Status code"`
}

func IsValidOrder(order models.Order, menuRepo dal.MenuRepository, inventRepo dal.InventoryRepository) error {
	// items := order.Items
	// menuItems, err := menuRepo.ReadAll()
	// if err != nil {
	// 	return err
	// }
	// inventIngredients, err := inventRepo.ReadAll()
	// if err != nil {
	// 	return err
	// }

	// inventMap := make(map[string]models.InventoryItem)
	// for _, inventIngr := range inventIngredients {
	// 	inventMap[inventIngr.IngredientID] = inventIngr
	// }

	// for _, item := range items {
	// 	foundMenuItem := false

	// 	for _, menuItem := range menuItems {
	// 		if item.ProductID == menuItem.ID {
	// 			foundMenuItem = true
	// 			quantity := item.Quantity

	// 			for _, ingredient := range menuItem.Ingredients {
	// 				inventIngr, foundIngredient := inventMap[ingredient.IngredientID]

	// 				if !foundIngredient {
	// 					return errors.New("invalid ingredient: " + ingredient.IngredientID)
	// 				}
	// 				if ingredient.Quantity*float64(quantity) > inventIngr.Quantity {
	// 					return errors.New("not enough ingredient: " + ingredient.IngredientID)
	// 				}
	// 			}
	// 			break
	// 		}
	// 	}
	// 	if !foundMenuItem {
	// 		return errors.New("order item doesn't exist in menu")
	// 	}
	// }

	return nil
}

func UpdateInventoryByOrder(inventory []models.InventoryItem, order models.Order, menuItems []models.MenuItem, subtract bool) ([]models.InventoryItem, error) {
	// // Map inventory items for quick lookup
	// inventoryMap := make(map[string]*models.InventoryItem)
	// for i := range inventory {
	// 	inventoryMap[inventory[i].IngredientID] = &inventory[i]
	// }

	// // Map menu items for quick lookup
	// menuMap := make(map[string]models.MenuItem)
	// for _, menuItem := range menuItems {
	// 	menuMap[menuItem.ID] = menuItem
	// }

	// // Iterate through order items
	// for _, orderItem := range order.Items {
	// 	menuItem, exists := menuMap[orderItem.ProductID]
	// 	if !exists {
	// 		return nil, errors.New("menu item not found: " + orderItem.ProductID)
	// 	}

	// 	// Calculate ingredient usage based on order quantity
	// 	for _, ingredient := range menuItem.Ingredients {
	// 		inventoryItem, found := inventoryMap[ingredient.IngredientID]
	// 		if !found {
	// 			return nil, errors.New("ingredient not found in inventory: " + ingredient.IngredientID)
	// 		}

	// 		// Calculate the adjustment amount
	// 		adjustment := ingredient.Quantity * float64(orderItem.Quantity)
	// 		if subtract {
	// 			if inventoryItem.Quantity < adjustment {
	// 				return nil, errors.New("not enough of ingredient: " + ingredient.IngredientID)
	// 			}
	// 			inventoryItem.Quantity -= adjustment
	// 		} else {
	// 			inventoryItem.Quantity += adjustment
	// 		}
	// 	}
	// }

	// Save updated inventory back to the file
	return inventory, nil
}

func IsMenuValid(items models.MenuItem) bool {
	if items.Price <= 0 {
		return false
	}
	if items.Description == "" || items.Name == "" {
		return false
	}
	for _, item := range items.Ingredients {
		if item.Quantity <= 0 || item.IngredientID <= 0 {
			return false
		}
	}
	return true
}

func IsOrderValid(order models.Order) bool {
	if order.CustomerID <= 0 {
		return false
	}
	for _, item := range order.Items {
		if item.Quantity <= 0 || item.MenuItemID <= 0 {
			return false
		}
	}
	return true
}

func IsInventoryValid(inventory models.InventoryItem) bool {
	if inventory.Quantity <= 0 || inventory.Name == "" || inventory.Unit == "" {
		return false
	}
	return true
}

func IsCustomerValid(customer models.Customer) bool {
	if customer.Number <= "" || customer.Name == "" || customer.Email == "" {
		return false
	}
	return true
}
