package models

type RecipeIngredient struct {
	RecipeID           int     `json:"recipe_id" db:"recipe_id"`
	IngredientID       int     `json:"ingredient_id" db:"ingredient_id"`
	RecipeIngredientID int     `json:"recipe_ingredient_id" db:"recipe_ingredient_id"`
	Quantity           float64 `json:"quantity" db:"quantity"`
	Measurement        string  `json:"measurement" db:"measurement"`
}
