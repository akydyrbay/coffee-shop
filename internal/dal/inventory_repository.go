package dal

import (
	"database/sql"
	"frappuccino/models"
	"log/slog"
)

type InventoryRepository interface {
	GetAll() ([]models.InventoryItem, error)
	Exists(item models.InventoryItem) (bool, error)
	ExistsID(id int) (bool, error)
	UpdateAll(item models.InventoryItem) error
	Insert(item models.InventoryItem) error
	DeleteItem(id int) error
	GetQuantity(inventoryID int) (float64, error)
	Deduct(inventoryID int, quantity float64, tx *sql.Tx) error
	GetInventory(id int) (models.InventoryItem, error)
	UseInventory(need_inventory map[string]float64) error
}

type inventoryRepo struct {
	DB *sql.DB
}

func NewInventoryRepo(db *sql.DB) *inventoryRepo {
	return &inventoryRepo{DB: db}
}

func (r *inventoryRepo) UseInventory(need_inventory map[string]float64) error {
	tx, err := r.DB.Begin() // Start a transaction
	if err != nil {
		return err
	}

	for ingredientID, quantity := range need_inventory {
		_, err := tx.Exec(
			"UPDATE inventory SET quantity = quantity - $1 WHERE id = $2 AND quantity >= $3",
			quantity, ingredientID, quantity,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit() // Commit the transaction
}

func (r *inventoryRepo) GetInventory(id int) (models.InventoryItem, error) {
	var item models.InventoryItem
	err := r.DB.QueryRow(`
        SELECT id, name, quantity, unit 
        FROM inventory 
        WHERE id = $1
    `, id).Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit)
	if err != nil {
		return models.InventoryItem{}, err
	}
	return item, nil
}

func (r *inventoryRepo) GetAll() ([]models.InventoryItem, error) {
	rows, err := r.DB.Query("SELECT id, name, quantity,unit FROM inventory")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit); err != nil {
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
		if inventory.Name == item.Name {
			return true, nil
		}
	}
	return false, nil
}

func (r *inventoryRepo) ExistsID(id int) (bool, error) {
	inventoryData, err := r.GetAll()
	if err != nil {
		return false, err
	}

	for _, inventory := range inventoryData {
		if inventory.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func (r *inventoryRepo) UpdateAll(item models.InventoryItem) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	slog.Info(item.Name)
	query := `
    Update inventory
    SET quantity = $1,
        updated_at = CURRENT_TIMESTAMP,
        name = $3,
        unit = $4
    WHERE id = $2;
	`
	_, err = tx.Exec(query, item.Quantity, item.ID, item.Name, item.Unit)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *inventoryRepo) Insert(item models.InventoryItem) error {
	_, err := r.DB.Exec(`
		INSERT INTO inventory (name, quantity, unit)
		VALUES ($1, $2, $3)
	`, item.Name, item.Quantity, item.Unit)
	return err
}

func (r *inventoryRepo) DeleteItem(id int) error {
	// Delete related inventory transactions first
	_, err := r.DB.Exec("Delete FROM inventory_transactions WHERE inventory_id = $1", id)
	if err != nil {
		return err
	}
	// Then delete the inventory item
	_, err = r.DB.Exec("Delete FROM inventory WHERE id = $1", id)
	return err
}

func (r *inventoryRepo) GetQuantity(inventoryID int) (float64, error) {
	var qty float64
	err := r.DB.QueryRow(`SELECT quantity FROM inventory WHERE id = $1`, inventoryID).Scan(&qty)
	return qty, err
}

func (r *inventoryRepo) Deduct(inventoryID int, quantity float64, tx *sql.Tx) error {
	// Deduct the quantity and update the updated_at timestamp.
	_, err := tx.Exec(`Update inventory SET quantity = quantity - $1, updated_at = NOW() WHERE id = $2`, quantity, inventoryID)
	return err
}
