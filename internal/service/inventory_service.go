package service

import (
	"errors"
	"fmt"
	"frappuccino/internal/dal"
	"frappuccino/models"
)

type InventoryService interface {
	Create(items []models.InventoryItem) models.Status
	Read(id int) (models.InventoryItem, models.Status)
	ReadAll() ([]models.InventoryItem, models.Status)
	Update(item models.InventoryItem) models.Status
	Delete(id int) models.Status
}

type inventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) *inventoryService {
	return &inventoryService{inventoryRepo: inventoryRepo}
}

func (s *inventoryService) Create(items []models.InventoryItem) models.Status {
	for _, item := range items {
		if !IsInventoryValid(item) {
			return models.Status{ErrorMessage: fmt.Errorf("invalid inventory item"), Code: 400}
		}
		b := IsInventoryValid(item)
		if !b {
			return models.Status{ErrorMessage: errors.New("incorrect inventory"), Code: 400}
		}
		exists, err := s.inventoryRepo.Exists(item)
		if err != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
		}
		if exists {
			return models.Status{ErrorMessage: errors.New("item already exists"), Code: 400}
		}
	}
	for _, item := range items {
		if err := s.inventoryRepo.Insert(item); err != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to insert inventory item: %w", err), Code: 500}
		}
	}

	return models.Status{ErrorMessage: nil, Code: 200}
}

func (s *inventoryService) Delete(id int) models.Status {
	exists, err := s.inventoryRepo.ExistsID(id)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
	}
	if !exists {
		return models.Status{ErrorMessage: errors.New("item does not exists"), Code: 400}
	}
	err = s.inventoryRepo.DeleteItem(id)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to delete inventory item: %w", err), Code: 500}
	}
	return models.Status{ErrorMessage: nil, Code: 200}
}

func (s *inventoryService) ReadAll() ([]models.InventoryItem, models.Status) {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return []models.InventoryItem{}, models.Status{ErrorMessage: fmt.Errorf("failed to get inventory items: %w", err), Code: 500}
	}

	return inventories, models.Status{ErrorMessage: nil, Code: 200}
}

func (s *inventoryService) Read(id int) (models.InventoryItem, models.Status) {
	inventoryItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return models.InventoryItem{}, models.Status{ErrorMessage: fmt.Errorf("failed to get inventory items: %w", err), Code: 500}
	}
	for _, inventoryItem := range inventoryItems {
		if inventoryItem.ID == id {
			return inventoryItem, models.Status{ErrorMessage: nil, Code: 200}
		}
	}
	return models.InventoryItem{}, models.Status{ErrorMessage: errors.New("inventory item not found"), Code: 404}
}

func (s *inventoryService) Update(item models.InventoryItem) models.Status {
	b := IsInventoryValid(item)
	if !b {
		return models.Status{ErrorMessage: errors.New("incorrect inventory"), Code: 400}
	}
	exists, err := s.inventoryRepo.Exists(item)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
	}
	if !exists {
		return models.Status{ErrorMessage: fmt.Errorf("inventory item with ID %d not found", item.ID), Code: 404}
	}

	// Perform the update directly since the item is confirmed to exist
	err = s.inventoryRepo.UpdateAll(item)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to update inventory item: %w", err), Code: 500}
	}

	return models.Status{ErrorMessage: nil, Code: 200}
}
