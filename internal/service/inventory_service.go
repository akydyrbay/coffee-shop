package service

import (
	"errors"

	"coffee-shop/internal/dal"
	"coffee-shop/models"
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
func (s *inventoryService) AddInventoryItem(items []models.InventoryItem) error {
	for _, item := range items {
		if !IsInventoryValid(item) {
			return errors.New("invalid inventory item")
		}

		exists, err := s.inventoryRepo.Exists(item)
		if err != nil {
			return fmt.Errorf("failed to check existence: %w", err)
		}
		if exists {
			return errors.New("item already exists")
		}

		if err := s.inventoryRepo.Insert(item); err != nil {
			return fmt.Errorf("failed to insert inventory item: %w", err)
		}
	}
	return nil
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}

	for i, inventory := range inventories {
		if inventory.ID == id {
			inventories = append(inventories[:i], inventories[i+1:]...)
		}
		if err := s.inventoryRepo.Insert(inventory); err != nil {
			return fmt.Errorf("failed to insert inventory item: %w", err)
		}
	}

	return nil
}

func (s *inventoryService) GetInventoryItem() ([]models.InventoryItem, error) {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return []models.InventoryItem{}, errors.New("failed to get inventory items")
	}

	return inventories, nil
}

func (s *inventoryService) GetInventoryItemById(id string) (models.InventoryItem, error) {
	inventoryItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}
	for _, inventoryItem := range inventoryItems {
		if inventoryItem.ID == id {
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
		if inventoryItems[i].ID == item.ID {
			err = s.inventoryRepo.UpdateAll(item)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("inventory item not found")
}
