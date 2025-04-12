package dal

import (
	"database/sql"
	"frappuccino/models"
	"strconv"
	"strings"
	"time"
)

// TODO implement me :()

type AggregationRepository interface {
	GetOrderQuantitybyDate(startDate, endDate string) ([]models.NumberOfOrderItems, error)
	FullSearchOrder(q string, maxPrice, minPrice float64) ([]models.OrderSearchResult, error)
	FullSearchMenu(q string, minPrice, maxPrice float64) ([]models.MenuItem, error)
	GetDayPeriod(month int, orderByDay *models.OrderByDayRequest) error
	GetMonthPeriod(year int, orderByMonth *models.OrderByMonthRequest) error
	GetLeftOvers(sortBy string, page, pageSize int) ([]models.LeftoverItem, int, error)
}

type aggrRepo struct {
	DB *sql.DB
}

func NewAggrRepo(db *sql.DB) *aggrRepo {
	return &aggrRepo{DB: db}
}

func (r *aggrRepo) GetOrderQuantitybyDate(startDate, endDate string) ([]models.NumberOfOrderItems, error) {
	var orderedItems []models.NumberOfOrderItems
	rows, err := r.DB.Query(`
SELECT 
    mi.name, 
    COALESCE(SUM(oi.quantity), 0) AS total_quantity
FROM 
    menu_items mi
INNER JOIN order_items oi 
    ON oi.menu_item_id = mi.id
INNER JOIN order_status_history osh 
    ON oi.order_id = osh.order_id
    AND osh.status = 'closed'
    AND osh.updated_at BETWEEN $1 AND $2
GROUP BY 
    mi.id, mi.name
ORDER BY 
    total_quantity DESC;
`, startDate, endDate)
	if err != nil {
		return orderedItems, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderedItem models.NumberOfOrderItems
		if err := rows.Scan(&orderedItem.Name, &orderedItem.Quantity); err != nil {
			return orderedItems, err
		}
		orderedItems = append(orderedItems, orderedItem)
	}
	return orderedItems, nil
}

func (r *aggrRepo) FullSearchOrder(q string, maxPrice, minPrice float64) ([]models.OrderSearchResult, error) {
	var orders []models.OrderSearchResult

	// If no max price provided (0), use the maximum total_amount from orders.
	if maxPrice == 0 {
		err := r.DB.QueryRow(`SELECT MAX(total_amount) FROM orders`).Scan(&maxPrice)
		if err != nil {
			return orders, err
		}
	}
	q = strings.ToLower(q)
	rows, err := r.DB.Query(`
WITH relevance AS (
    SELECT 
        o.id AS order_id,
        ROUND(
            CASE 
                WHEN LENGTH(REPLACE(LOWER(string_agg(mi.name || ' ' || mi.description || ' ' || c.name, ' ')), $1, '')) = 0
                THEN 1 
                WHEN LENGTH(REPLACE(LOWER(string_agg(mi.name || ' ' || mi.description || ' ' || c.name, ' ')), $1, '')) < LENGTH($1)
                THEN LENGTH(REPLACE(LOWER(string_agg(mi.name || ' ' || mi.description || ' ' || c.name, ' ')), $1, '')) * 1.0 / LENGTH($1)
                ELSE LENGTH($1) * 1.0 / LENGTH(REPLACE(LOWER(string_agg(mi.name || ' ' || mi.description || ' ' || c.name, ' ')), $1, ''))
            END,
            3
        ) AS relevance
    FROM order_items oi
    INNER JOIN menu_items mi ON oi.menu_item_id = mi.id
    INNER JOIN orders o ON oi.order_id = o.id
    INNER JOIN customers c ON o.customer_id = c.customer_id
    INNER JOIN order_status_history osh ON o.id = osh.order_id
    WHERE 
        (LOWER(mi.name) LIKE '%' || $1 || '%' 
         OR LOWER(mi.description) LIKE '%' || $1 || '%' 
         OR LOWER(c.name) LIKE '%' || $1 || '%')
        AND osh.status = 'closed'
    GROUP BY o.id
)
SELECT 
    o.id, 
    c.name AS customer_name, 
    o.total_amount, 
    r.relevance
FROM orders o
INNER JOIN customers c ON o.customer_id = c.customer_id
INNER JOIN relevance r ON o.id = r.order_id
WHERE 
    o.total_amount BETWEEN $2 AND $3
ORDER BY r.relevance DESC;
    `, q, minPrice, maxPrice)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	// Assuming you have an order repository with a function to get item names for a given order.
	orderRepo := NewOrderRepo(r.DB)
	for rows.Next() {
		var order models.OrderSearchResult
		err := rows.Scan(&order.ID, &order.CustomerName, &order.Total, &order.Relevance)
		if err != nil {
			return orders, err
		}
		// Fetch the list of menu item names for this order.
		order.Items, err = orderRepo.GetOrderItemsArray(order.ID)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *aggrRepo) FullSearchMenu(q string, maxPrice, minPrice float64) ([]models.MenuItem, error) {
	var menus []models.MenuItem
	// If no max price is provided (0), use the maximum current_price in the table.
	if maxPrice == 0 {
		err := r.DB.QueryRow(`SELECT MAX(current_price) FROM menu_items`).Scan(&maxPrice)
		if err != nil {
			return menus, err
		}
	}
	// Lowercase the query for case-insensitive search.
	q = strings.ToLower(q)
	rows, err := r.DB.Query(`
SELECT 
    id, 
    description, 
    name, 
    current_price,  
    ROUND(
        CASE 
            WHEN LENGTH(REPLACE(LOWER(name || ' ' || description), $1, '')) = 0
            THEN 1 
            WHEN LENGTH(REPLACE(LOWER(name || ' ' || description), $1, '')) < LENGTH($1)
            THEN LENGTH(REPLACE(LOWER(name || ' ' || description), $1, '')) * 1.0 / LENGTH($1)
            ELSE LENGTH($1) * 1.0 / LENGTH(REPLACE(LOWER(name || ' ' || description), $1, ''))
        END,
        3
    ) AS relevance
FROM 
    menu_items
WHERE 
    (LOWER(name) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%')
    AND current_price BETWEEN $2 AND $3
ORDER BY 
    relevance DESC;
    `, q, minPrice, maxPrice)
	if err != nil {
		return menus, err
	}
	defer rows.Close()
	for rows.Next() {
		var menu models.MenuItem
		// Assuming models.Menu has fields: ID int, Name string, Description string,
		// Price float64 (current_price), and Relevance float64.
		err := rows.Scan(&menu.ID, &menu.Description, &menu.Name, &menu.Price, &menu.Relevance)
		if err != nil {
			return menus, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (r *aggrRepo) GetDayPeriod(month int, orderByDay *models.OrderByDayRequest) error {
	// Use current year (or pass the year as a parameter if needed)
	year := time.Now().Year()

	query := `
 SELECT 
	 EXTRACT(day FROM o.created_at) AS day,
	 SUM(oi.quantity) AS total_quantity
 FROM orders o
 INNER JOIN order_items oi ON o.id = oi.order_id
 WHERE EXTRACT(month FROM o.created_at) = $1
   AND EXTRACT(year FROM o.created_at) = $2
   AND o.status = 'closed'
 GROUP BY day
 ORDER BY day;
	 `
	rows, err := r.DB.Query(query, month, year)
	if err != nil {
		return err
	}
	defer rows.Close()

	var orderedItems []models.OrderedItemDay
	for rows.Next() {
		var dayFloat float64 // EXTRACT returns float64
		var qty int
		if err := rows.Scan(&dayFloat, &qty); err != nil {
			return err
		}
		dayStr := strconv.Itoa(int(dayFloat))
		orderedItems = append(orderedItems, models.OrderedItemDay{
			Day:      dayStr,
			Quantity: qty,
		})
	}
	orderByDay.OrderedItems = orderedItems
	return nil
}

func (r *aggrRepo) GetMonthPeriod(year int, orderByMonth *models.OrderByMonthRequest) error {
	query := `
SELECT 
    TO_CHAR(o.created_at, 'FMMonth') AS month_name,
    SUM(oi.quantity) AS total_quantity
FROM orders o
INNER JOIN order_items oi ON o.id = oi.order_id
WHERE EXTRACT(year FROM o.created_at) = $1
  AND o.status = 'closed'
GROUP BY EXTRACT(month FROM o.created_at), TO_CHAR(o.created_at, 'FMMonth')
ORDER BY EXTRACT(month FROM o.created_at);
    `
	rows, err := r.DB.Query(query, year)
	if err != nil {
		return err
	}
	defer rows.Close()

	var orderedItems []models.OrderedItemMonth
	for rows.Next() {
		var monthName string
		var qty int
		if err := rows.Scan(&monthName, &qty); err != nil {
			return err
		}
		orderedItems = append(orderedItems, models.OrderedItemMonth{
			Month:    monthName,
			Quantity: qty,
		})
	}
	orderByMonth.OrderedItems = orderedItems
	return nil
}

func (r *aggrRepo) GetLeftOvers(sortBy string, page, pageSize int) ([]models.LeftoverItem, int, error) {
	var totalItems int
	var leftovers []models.LeftoverItem

	// Query to get paginated inventory leftovers.
	// Note:
	// - Use i.quantity from the inventory table.
	// - Use i.id instead of i.inventory_id.
	// - Join ingredient_prices using i.id.
	// - For demonstration, we assume inventory_transactions has a "price" column.
	rows, err := r.DB.Query(`
	WITH ingredient_prices AS (
		SELECT 
			it.inventory_id,
			COALESCE(SUM(it.price) / NULLIF(SUM(it.quantity), 0), 0) AS avg_price
		FROM 
			inventory_transactions it
		GROUP BY 
			it.inventory_id
	),
	sorted_inventory AS (
		SELECT 
			i.name, 
			i.quantity AS stock_level,
			COALESCE(ip.avg_price, 0) AS price,
			ROW_NUMBER() OVER (
				ORDER BY 
					CASE 
						WHEN $1 = 'price' THEN COALESCE(ip.avg_price, 0)
						WHEN $1 = 'quantity' THEN i.quantity
					END ASC
			) AS row_num
		FROM 
			inventory i
		LEFT JOIN 
			ingredient_prices ip 
			ON i.id = ip.inventory_id
	)
	SELECT 
		name, 
		stock_level AS quantity, 
		price
	FROM 
		sorted_inventory
	WHERE 
		row_num > (($2 - 1) * $3) 
		AND row_num <= ($2 * $3)
	ORDER BY row_num;
		`, sortBy, page, pageSize)
	if err != nil {
		return leftovers, 0, err
	}
	defer rows.Close()

	// Parse query results into leftovers slice.
	for rows.Next() {
		var item models.LeftoverItem
		err := rows.Scan(&item.Name, &item.Quantity, &item.Price)
		if err != nil {
			return leftovers, 0, err
		}
		leftovers = append(leftovers, item)
	}

	// Get the total number of inventory items.
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM inventory`).Scan(&totalItems)
	if err != nil {
		return leftovers, 0, err
	}

	return leftovers, totalItems, nil
}
