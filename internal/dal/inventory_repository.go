package dal

import (
	"encoding/json"
	"hot-coffee/models"
	"os"
)

type InventoryRepository interface {
	SaveAll(item []models.InventoryItem) error
	GetAll() ([]models.InventoryItem, error)
	Exists(item models.InventoryItem) (bool, error)
}

type inventoryRepo struct {
	path string
}

func NewInventoryRepo(path string) *inventoryRepo {
	return &inventoryRepo{path: path}
}

func (r *inventoryRepo) SaveAll(item []models.InventoryItem) error {
	jsonData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return err
		}
		file.Close()
	}
	err = os.WriteFile(r.path, jsonData, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (r *inventoryRepo) GetAll() ([]models.InventoryItem, error) {
	_, err := os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}
	byteValue, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}

	if len(byteValue) == 0 {
		return []models.InventoryItem{}, nil
	}

	var inventoryData []models.InventoryItem

	if err := json.Unmarshal(byteValue, &inventoryData); err != nil {
		return nil, err
	}

	return inventoryData, nil
}

func (r *inventoryRepo) Exists(item models.InventoryItem) (bool, error) {
	inventoryData, err := r.GetAll()
	if err != nil {
		return false, err
	}

	for _, inventory := range inventoryData {
		if inventory.IngredientID == item.IngredientID {
			return true, nil
		}
	}
	return false, nil
}
