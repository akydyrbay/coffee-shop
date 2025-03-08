package dal

import (
	"database/sql"

	"coffee-shop/models"
)

type InventoryRepository interface {
	SaveAll(item []models.InventoryItem) error
	GetAll() ([]models.InventoryItem, error)
	Exists(item models.InventoryItem) (bool, error)
}

type inventoryRepo struct {
	DB *sql.DB
}

func NewInventoryRepo(db *sql.DB) *inventoryRepo {
	return &inventoryRepo{DB: db}
}

func (r *inventoryRepo) SaveAll(items []models.InventoryItem) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.Exec(`INSERT INTO inventory (ingredient_id, name, quantity, unit)
	VALUES ($1, $2, $3, $4)
	`, item.IngredientID, item.Name, item.Quantity, item.Unit)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *inventoryRepo) GetAll() ([]models.InventoryItem, error) {
	rows, err := r.DB.Query("SELECT ingredient_id, name, quantity,unit FROM inventory")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
