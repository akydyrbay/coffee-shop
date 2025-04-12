package dal

import (
	"database/sql"
	"frappuccino/models"
	"log/slog"
)

type MenuRepository interface {
	SaveAll([]models.MenuItem) error
	GetAll() ([]models.MenuItem, error)
	Exists(menuItemID int) (bool, error)
	Insert(item models.MenuItem) error
	getIngredientsForMenuItem(menuID int) ([]models.MenuItemIngredient, error)
	Update(item models.MenuItem) error
	DeleteItem(id int) error
	GetByID(menuItemID int) (models.MenuItem, error)
	GetIngredients(menuItemID int) ([]models.MenuItemIngredient, error)
}

type menuRepo struct {
	DB *sql.DB
}

func NewMenuRepo(db *sql.DB) *menuRepo {
	return &menuRepo{DB: db}
}

func (r *menuRepo) Insert(item models.MenuItem) error {
	tx, err := r.DB.Begin() // Start a transaction
	if err != nil {
		return err
	}

	var id string
	err = tx.QueryRow(`
		INSERT INTO menu_items (name, description, category, current_price)
		VALUES ($1, $2, $3, $4) RETURNING id 
	`, item.Name, item.Description, item.Category, item.Price).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}
	slog.Info(id)

	// Insert each ingredient for this menu item
	for _, ing := range item.Ingredients {
		_, err = tx.Exec(`
	        INSERT INTO menu_item_ingredients (menu_item_id, inventory_id, quantity)
	        VALUES ($1, $2, $3)
	    `, id, ing.IngredientID, ing.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *menuRepo) SaveAll(menu []models.MenuItem) error {
	// jsonData, err := json.MarshalIndent(menu, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// _, err = os.Stat(r.path)
	// if os.IsNotExist(err) {
	// 	file, err := os.Create(r.path)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer file.Close()
	// }
	// err = os.WriteFile(r.path, jsonData, 0o644)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (r *menuRepo) GetAll() ([]models.MenuItem, error) {
	rows, err := r.DB.Query("SELECT id, name, description, category, current_price, is_available FROM menu_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.MenuItem
	for rows.Next() {
		var menu models.MenuItem
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Category, &menu.Price, &menu.IsAvailable); err != nil {
			return nil, err
		}
		// Now fetch ingredients for this menu item
		ingredients, err := r.getIngredientsForMenuItem(menu.ID)
		if err != nil {
			return nil, err
		}
		menu.Ingredients = ingredients
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepo) getIngredientsForMenuItem(menuID int) ([]models.MenuItemIngredient, error) {
	query := `
		SELECT menu_item_id, inventory_id, quantity
		FROM menu_item_ingredients
		WHERE menu_item_id = $1;
	`
	rows, err := r.DB.Query(query, menuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.MenuItemIngredient
	for rows.Next() {
		var ing models.MenuItemIngredient
		if err := rows.Scan(&ing.MenuItemID, &ing.IngredientID, &ing.Quantity); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients, nil
}

func (r *menuRepo) Exists(menuId int) (bool, error) {
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

func (r *menuRepo) Update(item models.MenuItem) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`Update menu_items SET name = $1, category = $2, description = $3, current_price = $4, is_available = $5 WHERE id = $6`,
		item.Name, item.Category, item.Description, item.Price, item.IsAvailable, item.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`Delete FROM menu_item_ingredients WHERE menu_item_id = $1`, item.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, ing := range item.Ingredients {
		_, err = tx.Exec(`
            INSERT INTO menu_item_ingredients (menu_item_id, inventory_id, quantity)
            VALUES ($1, $2, $3)`,
			item.ID, ing.IngredientID, ing.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *menuRepo) DeleteItem(id int) error {
	_, err := r.DB.Exec("Delete FROM menu_items WHERE id = $1", id)
	return err
}

func (r *menuRepo) GetByID(menuItemID int) (models.MenuItem, error) {
	var item models.MenuItem
	err := r.DB.QueryRow(`SELECT id, name, current_price FROM menu_items WHERE id = $1`, menuItemID).
		Scan(&item.ID, &item.Name, &item.Price)
	return item, err
}

func (r *menuRepo) GetIngredients(menuItemID int) ([]models.MenuItemIngredient, error) {
	rows, err := r.DB.Query(`SELECT inventory_id, quantity FROM menu_item_ingredients WHERE menu_item_id = $1`, menuItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ingredients []models.MenuItemIngredient
	for rows.Next() {
		var ing models.MenuItemIngredient
		if err := rows.Scan(&ing.IngredientID, &ing.Quantity); err != nil {
			return nil, err
		}
		// Optionally, fetch the ingredient name from inventory if needed.
		var name string
		err := r.DB.QueryRow(`SELECT name FROM inventory WHERE id = $1`, ing.IngredientID).Scan(&name)
		if err == nil {
			ing.Name = name
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients, nil
}
