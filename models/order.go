package models

import "time"

type Order struct {
	ID          int         `json:"id"`
	CustomerID  int         `json:"customer_id"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

type OrderItem struct {
	// ID         int     `json:"id"`
	OrderID    int     `json:"order_id"`
	MenuItemID int     `json:"menu_item_id"`
	Quantity   float64 `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

type OrderStatusHistory struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	Notes     string `json:"notes"`
	UpdatedAt string `json:"updated_at"`
}

type OrderSearchResult struct {
	ID           int      `json:"id"`
	CustomerName string   `json:"customer_name"`
	Total        float64  `json:"total_amount"`
	Items        []string `json:"items"`
	Relevance    float64  `json:"relevance"`
}
