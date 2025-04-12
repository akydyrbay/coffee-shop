package service

import (
	"errors"
	"fmt"
	"frappuccino/internal/dal"
	"frappuccino/models"
)

type MenuServiceInterface interface {
	Create(items []models.MenuItem) models.Status
	ReadAll() ([]models.MenuItem, models.Status)
	Read(id int) (models.MenuItem, models.Status)
	Update(item models.MenuItem) models.Status
	Delete(id int) models.Status
}

type menuService struct {
	menuRepo dal.MenuRepository
}

func NewMenuService(menuRepo dal.MenuRepository) *menuService {
	return &menuService{menuRepo: menuRepo}
}

func (s *menuService) Create(items []models.MenuItem) models.Status {
	for _, item := range items {
		if !IsMenuValid(item) {
			return models.Status{ErrorMessage: errors.New("invalid menu item"), Code: 400}
		}

		exists, err := s.menuRepo.Exists(item.ID)
		if err != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
		}
		if exists {
			return models.Status{ErrorMessage: errors.New("menu item already exists"), Code: 400}
		}
	}
	for _, item := range items {
		item.IsAvailable = true
		if err := s.menuRepo.Insert(item); err != nil {
			return models.Status{ErrorMessage: fmt.Errorf("failed to insert menu item: %w", err), Code: 500}
		}
	}
	return models.Status{ErrorMessage: nil, Code: 200}
}

func (s *menuService) ReadAll() ([]models.MenuItem, models.Status) {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, models.Status{ErrorMessage: fmt.Errorf("failed to get menu items: %w", err), Code: 500}
	}
	return menuItems, models.Status{ErrorMessage: nil, Code: 200}
}

func (s *menuService) Read(id int) (models.MenuItem, models.Status) {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return models.MenuItem{}, models.Status{ErrorMessage: fmt.Errorf("failed to get menu items: %w", err), Code: 500}
	}
	for _, menuItem := range menuItems {
		if menuItem.ID == id {
			return menuItem, models.Status{ErrorMessage: nil, Code: 200}
		}
	}
	return models.MenuItem{}, models.Status{ErrorMessage: errors.New("menu item not found"), Code: 404}
}

func (s *menuService) Update(item models.MenuItem) models.Status {
	if !IsMenuValid(item) {
		return models.Status{ErrorMessage: errors.New("invalid menu item"), Code: 400}
	}
	exists, err := s.menuRepo.Exists(item.ID)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
	}
	if !exists {
		return models.Status{ErrorMessage: errors.New("menu item not found"), Code: 404}
	}

	if err := s.menuRepo.Update(item); err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to update menu item: %w", err), Code: 500}
	}
	return models.Status{ErrorMessage: nil, Code: 200}
}

func (s *menuService) Delete(id int) models.Status {
	exists, err := s.menuRepo.Exists(id)
	if err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to check existence: %w", err), Code: 500}
	}
	if !exists {
		return models.Status{ErrorMessage: errors.New("menu item not found"), Code: 404}
	}

	if err := s.menuRepo.DeleteItem(id); err != nil {
		return models.Status{ErrorMessage: fmt.Errorf("failed to delete menu item: %w", err), Code: 500}
	}
	return models.Status{Code: 204} // No content for successful deletion
}
