package controllers

import (
	"backend/models" // Import your models package where you have your struct definitions
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// CreateRecipeStep creates a new recipe step.
// CreateRecipeStep godoc
// @Summary Create a new recipe step
// @Description Add a new recipe step to the database
// @Tags recipe_steps
// @Accept json
// @Produce json
// @Param recipe_step body models.RecipeStep true "Add recipe step"
// @Success 201 {object} models.RecipeStep
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-steps [post]
func CreateRecipeStep(c *gin.Context, db *sql.DB) {
	// 1. Define a variable to hold the recipe step data.
	var recipeStep models.RecipeStep

	// 2. Bind the request JSON to the recipe step struct.
	if err := c.ShouldBindJSON(&recipeStep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Perform validation and save the recipe step to the database.
	sqlQuery := `
		INSERT INTO recipe_steps (recipe_id, step_number, step_description)
		VALUES ($1, $2, $3)
		RETURNING recipe_step_id`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	var recipeStepID int
	err = stmt.QueryRow(recipeStep.RecipeID, recipeStep.StepNumber, recipeStep.StepDescription).Scan(&recipeStepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 4. Return a JSON response with the created recipe step.
	createdRecipeStep := models.RecipeStep{
		RecipeStepID:    recipeStepID,
		RecipeID:        recipeStep.RecipeID,
		StepNumber:      recipeStep.StepNumber,
		StepDescription: recipeStep.StepDescription,
	}

	c.JSON(http.StatusCreated, createdRecipeStep)
}

// DeleteRecipeStep deletes a recipe step.
// DeleteRecipeStep godoc
// @Summary Delete a recipe step
// @Description Delete a recipe step from the database
// @Tags recipe_steps
// @Accept json
// @Produce json
// @Param id path int true "Recipe Step ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-steps/{id} [delete]
func DeleteRecipeStep(c *gin.Context, db *sql.DB) {
	// 1. Get the recipe step ID from the URL parameter.
	recipeStepID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe step ID"})
		return
	}

	// 2. Perform validation and delete the recipe step from the database.
	sqlQuery := `DELETE FROM recipe_steps WHERE recipe_step_id = $1`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(recipeStepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 3. Return a 204 response.
	c.JSON(http.StatusNoContent, nil)
}

// GetRecipeStep retrieves a single recipe step by ID.
// GetRecipeStep godoc
// @Summary Get a recipe step by ID
// @Description Get a single recipe step from the database by ID
// @Tags recipe_steps
// @Accept json
// @Produce json
// @Param id path int true "Recipe Step ID"
// @Success 200 {object} models.RecipeStep
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-steps/{id} [get]
func GetRecipeStep(c *gin.Context, db *sql.DB) {
	// 1. Get the recipe step ID from the URL parameter.
	recipeStepID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe step ID"})
		return
	}

	// 2. Fetch the recipe step from the database by ID.
	sqlQuery := `SELECT recipe_step_id, recipe_id, step_number, step_description FROM recipe_steps WHERE recipe_step_id = $1`

	var recipeStep models.RecipeStep
	err = db.QueryRow(sqlQuery, recipeStepID).Scan(&recipeStep.RecipeStepID, &recipeStep.RecipeID, &recipeStep.StepNumber, &recipeStep.StepDescription)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe step not found"})
		return
	}

	// 3. Return a JSON response with the recipe step.
	c.JSON(http.StatusOK, recipeStep)
}

// UpdateRecipeStep updates a recipe step by ID.
// UpdateRecipeStep godoc
// @Summary Update a recipe step
// @Description Update a recipe step by ID
// @Tags recipe_steps
// @Accept json
// @Produce json
// @Param id path int true "Recipe Step ID"
// @Param recipe_step body models.RecipeStep true "Update recipe step"
// @Success 200 {object} models.RecipeStep
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-steps/{id} [put]
func UpdateRecipeStep(c *gin.Context, db *sql.DB) {
	// 1. Get the recipe step ID from the URL parameter.
	recipeStepID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe step ID"})
		return
	}

	// 2. Define a variable to hold the recipe step data.
	var recipeStep models.RecipeStep

	// 3. Bind the request JSON to the recipe step struct.
	if err := c.ShouldBindJSON(&recipeStep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Perform validation and update the recipe step in the database.
	sqlQuery := `
		UPDATE recipe_steps
		SET recipe_id = $1, step_number = $2, step_description = $3
		WHERE recipe_step_id = $4`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(recipeStep.RecipeID, recipeStep.StepNumber, recipeStep.StepDescription, recipeStepID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 5. Return a JSON response with the updated recipe step.
	updatedRecipeStep := models.RecipeStep{
		RecipeStepID:    recipeStepID,
		RecipeID:        recipeStep.RecipeID,
		StepNumber:      recipeStep.StepNumber,
		StepDescription: recipeStep.StepDescription,
	}

	c.JSON(http.StatusOK, updatedRecipeStep)
}
