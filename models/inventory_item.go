package models

type InventoryItem struct {
	IngredientID string `json:"ingredient_id"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	Unit         string `json:"unit"`
}
