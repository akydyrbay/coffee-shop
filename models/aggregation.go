package models

type TotalSales struct {
	Sales float64 `json:"total_sales: "`
}

type NumberOfOrderItems struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

type FullSearchStruct struct {
	MenuItems    []MenuItem          `json:"menu_items"`
	Orders       []OrderSearchResult `json:"orders"`
	TotalMatches int                 `json:"total_matches"`
}

type OrderByDayRequest struct {
	Period       string           `json:"period"`
	Month        string           `json:"month"`
	OrderedItems []OrderedItemDay `json:"orderedItems"`
}

type OrderByDayRequest2 struct {
	Period       string           `json:"period"`
	Month        string           `json:"month"`
	OrderedItems []map[string]int `json:"orderedItems"`
}

type OrderedItemDay struct {
	Day      string `json:"day"`
	Quantity int    `json:"quantity"`
}

type OrderByMonthRequest struct {
	Period       string             `json:"period"`
	Year         string             `json:"year"`
	OrderedItems []OrderedItemMonth `json:"orderedItems"`
}

type OrderByMonthRequest2 struct {
	Period       string           `json:"period"`
	Year         string           `json:"year"`
	OrderedItems []map[string]int `json:"orderedItems"`
}

type OrderedItemMonth struct {
	Month    string `json:"month"`
	Quantity int    `json:"quantity"`
}

type LeftoverItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type LeftoversResponse struct {
	CurrentPage int            `json:"currentPage"`
	HasNextPage bool           `json:"hasNextPage"`
	PageSize    int            `json:"pageSize"`
	TotalPages  int            `json:"totalPages"`
	Data        []LeftoverItem `json:"data"`
}

type BatchProcessResponse struct {
	ProcessedOrders []ProcessedOrder    `json:"processed_orders"`
	Summary         BatchProcessSummary `json:"summary"`
}

type BatchProcessRequest struct {
	Orders []Order `json:"orders"` // Массив заказов
}

type ProcessedOrder struct {
	OrderID    int     `json:"order_id"`
	CustomerID int     `json:"customer_id"` // Matches customer_id
	Status     string  `json:"status"`
	Total      float64 `json:"total"`
	Reason     string  `json:"reason,omitempty"`
}

type BatchProcessSummary struct {
	TotalOrders      int               `json:"total_orders"`
	Accepted         int               `json:"accepted"`
	Rejected         int               `json:"rejected"`
	TotalRevenue     float64           `json:"total_revenue"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}

type InventoryUpdate struct {
	IngredientID string  `json:"ingredient_id"`
	Name         string  `json:"name,omitempty"`
	QuantityUsed float64 `json:"quantity_used"`
	Remaining    float64 `json:"remaining"`
}
