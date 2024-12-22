package service

import (
	"errors"

	"coffee/internal/dal"
	"coffee/models"
)

type MenuServiceInterface interface {
	AddMenuItem(item models.MenuItem) error
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItemById(id string) (models.MenuItem, error)
	UpdateMenuItem(item models.MenuItem) error
	DeleteMenuItemById(id string) error
}

type menuService struct {
	menuRepo dal.MenuRepository
}

func NewMenuService(menuRepo dal.MenuRepository) *menuService {
	return &menuService{menuRepo: menuRepo}
}

func (s *menuService) AddMenuItem(item models.MenuItem) error {
	if !IsMenuValid(item) {
		return errors.New("invalid menu item")
	}
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}
	pres, err := s.menuRepo.Exists(item.ID)
	if err != nil {
		return err
	}
	if pres {
		return errors.New("menu item already exists")
	}
	menuItems = append(menuItems, item)
	err = s.menuRepo.SaveAll(menuItems)
	if err != nil {
		return err
	}
	return nil
}

func (s *menuService) GetAllMenuItems() ([]models.MenuItem, error) {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (s *menuService) GetMenuItemById(id string) (models.MenuItem, error) {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}
	for _, menuItem := range menuItems {
		if menuItem.ID == id {
			return menuItem, nil
		}
	}
	return models.MenuItem{}, errors.New("menu item not found")
}

func (s *menuService) UpdateMenuItem(item models.MenuItem) error {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}
	if !IsMenuValid(item) {
		return errors.New("invalid menu item")
	}
	for i := range menuItems {
		if menuItems[i].ID == item.ID {
			menuItems[i] = item
			err = s.menuRepo.SaveAll(menuItems)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("menu item not found")
}

func (s *menuService) DeleteMenuItemById(id string) error {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}
	for index, menuItem := range menuItems {
		if menuItem.ID == id {
			menuItems = append(menuItems[:index], menuItems[index+1:]...)
			s.menuRepo.SaveAll(menuItems)
			return nil
		}
	}
	return errors.New("menu item not found")
}
