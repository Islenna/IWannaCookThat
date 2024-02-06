package models

type RecipeRequest struct {
	RecipeName        string `json:"recipe_name" db:"recipe_name"`
	RecipeDescription string `json:"recipe_description" db:"recipe_description"`
	CookTime          int    `json:"cook_time" db:"cook_time"`
}

// Recipe defines the structure for a recipe.
type Recipe struct {
	RecipeID          int    `json:"recipe_id" db:"recipe_id"`
	RecipeName        string `json:"recipe_name" db:"recipe_name"`
	RecipeDescription string `json:"recipe_description" db:"recipe_description"`
	CookTime          int    `json:"cook_time" db:"cook_time"`
}
