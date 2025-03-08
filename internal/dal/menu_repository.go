package dal

import (
	"encoding/json"
	"os"

	"coffee-shop/models"
)

type MenuRepository interface {
	SaveAll([]models.MenuItem) error
	GetAll() ([]models.MenuItem, error)
	Exists(menuItemID string) (bool, error)
}

type menuRepo struct {
	path string
}

func NewMenuRepo(path string) *menuRepo {
	return &menuRepo{path: path}
}

func (r *menuRepo) SaveAll(menu []models.MenuItem) error {
	jsonData, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	err = os.WriteFile(r.path, jsonData, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (r *menuRepo) GetAll() ([]models.MenuItem, error) {
	var result []models.MenuItem
	_, err := os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}
	byteMenu, err := os.ReadFile(r.path)
	if len(byteMenu) == 0 {
		return result, err
	}

	if err := json.Unmarshal(byteMenu, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (r *menuRepo) Exists(menuId string) (bool, error) {
	menuItems, err := r.GetAll()
	if err != nil {
		return false, err
	}
	for _, menuItem := range menuItems {
		if menuItem.ID == menuId {
			return true, nil
		}
	}
	return false, nil
}
