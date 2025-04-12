package dal

import (
	"database/sql"
	"errors"
	"frappuccino/models"
	"net/http"
	"strconv"
)

const pathOrder = "data/orders.json"

type OrderRepository interface {
	SaveAll([]models.Order) error
	GetAll() ([]models.Order, error)
	Exists(orderID int) (bool, error)
	Insert(item models.Order) (models.Order, error)
	GetOrderItem(id int) ([]models.OrderItem, error)
	Update(item models.Order) error
	DeleteItem(id int) error
	BeginTx() (*sql.Tx, error)
	GetOrderItemsArray(id int) ([]string, error)
	GetNeedInventory(order models.Order) (map[string]float64, error)
	IsEnough(needInventory map[string]float64) (bool, error)
	CheckOrder(newOrder models.Order) (int, error)
	GetPriceAtOrderItems(order *models.Order) error
	GetTotalAmount(order *models.Order) error
	GetLastOrderId() (int, error)
	CloseOrder(order_id int) error
	GetOrder(order_id int) (models.Order, error)
}

type orderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *orderRepo {
	return &orderRepo{DB: db}
}

func (r *orderRepo) BeginTx() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *orderRepo) GetOrder(order_id int) (models.Order, error) {
	orders, err := r.GetAll()
	if err != nil {
		return models.Order{}, err
	}
	var result models.Order
	for _, order := range orders {
		if order.ID == order_id {
			result = order
		}
	}
	return result, nil
}

func (r *orderRepo) CloseOrder(order_id int) error {
	order, err := r.GetOrder(order_id)
	if err != nil {
		return err
	}
	need_invent, err := r.GetNeedInventory(order)
	if err != nil {
		return err
	}
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE orders
	SET status='closed'
	WHERE order_id=$1
	`, order_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO order_status_history (order_id, status)
	VALUES($1,'closed')
	`, order_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = NewInventoryRepo(r.DB).UseInventory(need_invent)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *orderRepo) GetLastOrderId() (int, error) {
	var lastID int
	err := r.DB.QueryRow(`SELECT id FROM orders
	ORDER BY created_at DESC
	LIMIT 1;
	`).Scan(&lastID)
	return lastID, err
}

func (r *orderRepo) GetTotalAmount(order *models.Order) error {
	var total float64
	for _, orderItem := range order.Items {
		total += orderItem.UnitPrice
	}
	order.TotalAmount = total
	return nil
}

func (r *orderRepo) GetPriceAtOrderItems(order *models.Order) error {
	for index, orderItem := range order.Items {
		var PriceAtOrderTime float64
		err := r.DB.QueryRow("SELECT unit_price*$1 FROM order_items WHERE menu_item_id=$2", orderItem.UnitPrice, orderItem.MenuItemID).
			Scan(&PriceAtOrderTime)
		if err != nil {
			return err
		}
		order.Items[index].UnitPrice = PriceAtOrderTime
	}
	return nil
}

func (r *orderRepo) CheckOrder(newOrder models.Order) (int, error) {
	// Check if products exist
	for _, item := range newOrder.Items {
		exists, err := r.Exists(item.MenuItemID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if !exists {
			return http.StatusBadRequest, errors.New("ordered item does not exist: " + strconv.Itoa(item.MenuItemID))
		}
	}

	// Check inventory requirements
	needInventory, err := r.GetNeedInventory(newOrder)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	enough, err := r.IsEnough(needInventory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !enough {
		return http.StatusBadRequest, errors.New("not enough inventory")
	}
	return http.StatusOK, nil
}

func (r *orderRepo) IsEnough(needInventory map[string]float64) (bool, error) {
	for ingredientID, requiredQuantity := range needInventory {
		var availableQuantity float64
		err := r.DB.QueryRow(`
			SELECT quantity
			FROM inventory
			WHERE id = $1
		`, ingredientID).Scan(&availableQuantity)
		if err != nil {
			return false, err
		}
		if availableQuantity <= requiredQuantity {
			return false, nil
		}
	}
	return true, nil
}

func (r *orderRepo) Insert(item models.Order) (models.Order, error) {
	tx, err := r.DB.Begin() // Start a transaction
	if err != nil {
		return models.Order{}, err
	}

	var id string
	err = tx.QueryRow(`
		 INSERT INTO orders (customer_id, total_amount)
		 VALUES ($1, $2) RETURNING id 
	  `, item.CustomerID, item.TotalAmount).Scan(&id)
	if err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// Set the ID in the returned order
	item.ID, err = strconv.Atoi(id)
	if err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// Insert each order item
	for _, ord := range item.Items {
		_, err = tx.Exec(`
			  INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price)
			  VALUES ($1, $2, $3, $4)
		  `, id, ord.MenuItemID, ord.Quantity, ord.UnitPrice)
		if err != nil {
			tx.Rollback()
			return models.Order{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return models.Order{}, err
	}

	return item, nil
}

func (r *orderRepo) SaveAll(order []models.Order) error {
	// jsonData, err := json.MarshalIndent(order, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// _, err = os.Stat(r.path)
	// if os.IsNotExist(err) {
	// 	file, err := os.Create(r.path)
	// 	if err != nil {
	// 		return errstring
	// 	return err
	// }
	return nil
}

func (r *orderRepo) GetAll() ([]models.Order, error) {
	rows, err := r.DB.Query("SELECT id, customer_id, status, total_amount, created_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.Status, &order.TotalAmount, &order.CreatedAt); err != nil {
			return nil, err
		}
		// Now fetch ingredients for this menu item
		items, err := r.GetOrderItem(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepo) GetNeedInventory(order models.Order) (map[string]float64, error) {
	needInventory := make(map[string]float64)

	for _, item := range order.Items {
		rows, err := r.DB.Query(`
			SELECT inventory_id,
			quantity * $1
			FROM menu_item_ingredients 
			WHERE menu_item_id = $2
		`, item.Quantity, item.MenuItemID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var ingredientID string
			var quantity float64
			if err := rows.Scan(&ingredientID, &quantity); err != nil {
				return nil, err
			}
			needInventory[ingredientID] += quantity
		}
	}
	return needInventory, nil
}

func (r *orderRepo) GetOrderItem(orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT order_id, menu_item_id, quantity, unit_price
		FROM order_items
		WHERE order_id = $1;
	`
	rows, err := r.DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.OrderID, &item.MenuItemID, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *orderRepo) Exists(orderID int) (bool, error) {
	orderItems, err := r.GetAll()
	if err != nil {
		return false, err
	}
	for _, orderItem := range orderItems {
		if orderItem.ID == orderID {
			return true, nil
		}
	}
	return false, nil
}

func (r *orderRepo) Update(order models.Order) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        Update orders 
        SET customer_id = $1, 
            total_amount = $2, 
            status = $3
        WHERE id = $4`,
		order.CustomerID, order.TotalAmount, order.Status, order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`Delete FROM order_items WHERE order_id = $1`, order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
            INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price)
            VALUES ($1, $2, $3, $4)`,
			order.ID, item.MenuItemID, item.Quantity, item.UnitPrice)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if order.Status != "" {
		_, err = tx.Exec(`
            INSERT INTO order_status_history (order_id, status, notes, updated_at)
            VALUES ($1, $2, $3, NOW())`,
			order.ID, order.Status, "Status updated via API")
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepo) DeleteItem(id int) error {
	_, err := r.DB.Exec("Delete FROM orders WHERE id = $1", id)
	return err
}

func (r *orderRepo) GetOrderItemsArray(id int) ([]string, error) {
	var orderItems []string
	rows, err := r.DB.Query(`SELECT name
	FROM order_items
	INNER JOIN menu_items USING(menu_item_id)
	WHERE order_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderItem string
		err := rows.Scan(&orderItem)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}
