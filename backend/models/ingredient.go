package models

type IngredientRequest struct {
	IngredientName        string `json:"ingredient_name" binding:"required"`
	IngredientDescription string `json:"ingredient_description"`
}

type Ingredient struct {
	IngredientID          int    `json:"ingredient_id" db:"ingredient_id"`
	IngredientName        string `json:"ingredient_name" db:"ingredient_name"`
	IngredientDescription string `json:"ingredient_description" db:"ingredient_description"`
}
