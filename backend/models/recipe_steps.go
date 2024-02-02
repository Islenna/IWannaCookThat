package models

type RecipeStep struct {
	RecipeStepID    int    `json:"recipe_step_id" db:"recipe_step_id"`
	RecipeID        int    `json:"recipe_id" db:"recipe_id"`
	StepID          int    `json:"step_id" db:"step_id"`
	StepNumber      int    `json:"step_number" db:"step_number"`
	StepDescription string `json:"step_description" db:"step_description"`
}
