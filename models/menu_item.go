package models

type MenuItem struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Category    string               `json:"category"`
	Price       float64              `json:"current_price"`
	IsAvailable bool                 `json:"is_available"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
	Relevance   float64              `json:"relevance"`
}

type MenuItemIngredient struct {
	MenuItemID   int     `json:"menu_item_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Name         string  `json:"name"`
}
