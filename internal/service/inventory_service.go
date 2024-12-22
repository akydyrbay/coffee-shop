package service

import (
	"errors"

	"coffee/internal/dal"
	"coffee/models"
)

type InventoryService interface {
	AddInventoryItem(item models.InventoryItem) error
	DeleteInventoryItem(id string) error
	GetInventoryItem() ([]models.InventoryItem, error)
	GetInventoryItemById(id string) (models.InventoryItem, error)
	UpdateInventoryItem(item models.InventoryItem) error
}

type inventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) *inventoryService {
	return &inventoryService{inventoryRepo: inventoryRepo}
}

func (s *inventoryService) AddInventoryItem(item models.InventoryItem) error {
	if !IsInventoryValid(item) {
		return errors.New("invalid inventory item")
	}
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return errors.New("failed to get inventory items")
	}

	if b, _ := s.inventoryRepo.Exists(item); b {
		return errors.New("item already exists")
	}

	inventories = append(inventories, item)

	if err := s.inventoryRepo.SaveAll(inventories); err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}

	for i, inventory := range inventories {
		if inventory.IngredientID == id {
			inventories = append(inventories[:i], inventories[i+1:]...)
		}
	}
	if err := s.inventoryRepo.SaveAll(inventories); err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) GetInventoryItem() ([]models.InventoryItem, error) {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return []models.InventoryItem{}, err
	}

	return inventories, nil
}

func (s *inventoryService) GetInventoryItemById(id string) (models.InventoryItem, error) {
	inventoryItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}
	for _, inventoryItem := range inventoryItems {
		if inventoryItem.IngredientID == id {
			return inventoryItem, nil
		}
	}
	return models.InventoryItem{}, errors.New("inventory item not found")
}

func (s *inventoryService) UpdateInventoryItem(item models.InventoryItem) error {
	inventoryItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}
	for i := range inventoryItems {
		if inventoryItems[i].IngredientID == item.IngredientID {

			inventoryItems[i].Quantity += item.Quantity
			err = s.inventoryRepo.SaveAll(inventoryItems)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("inventory item not found")
}
