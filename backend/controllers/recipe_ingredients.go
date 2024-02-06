package controllers

import (
	"backend/models" // Import your models package where you have your struct definitions
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// CreateRecipeIngredient creates a new recipe ingredient.
// CreateRecipeIngredient godoc
// @Summary Create a new recipe ingredient
// @Description Add a new recipe ingredient to the database
// @Tags recipe_ingredients
// @Accept json
// @Produce json
// @Param recipe_ingredient body models.RecipeIngredientRequest true "Add recipe ingredient"
// @Success 201 {object} models.RecipeIngredient
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-ingredients [post]
func CreateRecipeIngredient(c *gin.Context, db *sql.DB) {
	// 1. Define a variable to hold the recipe ingredient data.
	var recipeIngredient models.RecipeIngredientRequest

	// 2. Bind the request JSON to the recipe ingredient struct.
	if err := c.ShouldBindJSON(&recipeIngredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Perform validation and save the recipe ingredient to the database.
	sqlQuery := `
		INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING recipe_ingredient_id`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	var recipeIngredientID int
	err = stmt.QueryRow(recipeIngredient.RecipeID, recipeIngredient.IngredientID, recipeIngredient.Quantity).Scan(&recipeIngredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 4. Return a JSON response with the created recipe ingredient.
	createdRecipeIngredient := models.RecipeIngredient{
		RecipeIngredientID: recipeIngredientID,
		RecipeID:           recipeIngredient.RecipeID,
		IngredientID:       recipeIngredient.IngredientID,
		Quantity:           recipeIngredient.Quantity,
	}

	c.JSON(http.StatusCreated, createdRecipeIngredient)
}

// GetRecipeIngredients returns all recipe ingredients.
// GetRecipeIngredients godoc
// @Summary Get all recipe ingredients
// @Description Get all recipe ingredients from the database
// @Tags recipe_ingredients
// @Accept json
// @Produce json
// @Success 200 {array} models.RecipeIngredient
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-ingredients [get]
func GetRecipeIngredients(c *gin.Context, db *sql.DB) {
	// 1. Define a variable to hold the recipe ingredients.
	var recipeIngredients []models.RecipeIngredient

	// 2. Query the database for all recipe ingredients.
	sqlQuery := `SELECT recipe_ingredient_id, recipe_id, ingredient_id, quantity FROM recipe_ingredients`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
		return
	}
	defer rows.Close()

	// 3. Iterate over the rows and add each recipe ingredient to the slice.
	for rows.Next() {
		var recipeIngredient models.RecipeIngredient
		err := rows.Scan(&recipeIngredient.RecipeIngredientID, &recipeIngredient.RecipeID, &recipeIngredient.IngredientID, &recipeIngredient.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
			return
		}
		recipeIngredients = append(recipeIngredients, recipeIngredient)
	}

	// 4. Return a JSON response with the recipe ingredients.
	c.JSON(http.StatusOK, recipeIngredients)
}

// GetRecipeIngredient retrieves a single recipe ingredient by ID.
// GetRecipeIngredient godoc
// @Summary Get a single recipe ingredient
// @Description Get details of a specific recipe ingredient by ID
// @Tags recipe_ingredients
// @Accept json
// @Produce json
// @Param id path int true "Recipe Ingredient ID"
// @Success 200 {object} models.RecipeIngredient "Recipe ingredient retrieved"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Recipe ingredient not found"
// @Router /recipe-ingredients/{id} [get]
func GetRecipeIngredient(c *gin.Context, db *sql.DB) {
	// 1. Extract the recipe ingredient ID from the URL parameter.
	recipeIngredientIDStr := c.Param("id")
	// If the recipeIngredientID is not a valid integer, return a 400 Bad Request response.
	recipeIngredientID, err := strconv.Atoi(recipeIngredientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ingredient ID"})
		return
	}

	// 2. Query the database for the recipe ingredient.
	sqlQuery := `SELECT recipe_ingredient_id, recipe_id, ingredient_id, quantity FROM recipe_ingredients WHERE recipe_ingredient_id = $1`
	var recipeIngredient models.RecipeIngredient
	err = db.QueryRow(sqlQuery, recipeIngredientID).Scan(&recipeIngredient.RecipeIngredientID, &recipeIngredient.RecipeID, &recipeIngredient.IngredientID, &recipeIngredient.Quantity)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe ingredient not found"})
		return
	}

	// 3. Return a JSON response with the retrieved recipe ingredient.
	c.JSON(http.StatusOK, recipeIngredient)
}

// UpdateRecipeIngredient updates a recipe ingredient by ID.
// UpdateRecipeIngredient godoc
// @Summary Update a recipe ingredient
// @Description Update a recipe ingredient by ID
// @Tags recipe_ingredients
// @Accept json
// @Produce json
// @Param id path int true "Recipe Ingredient ID"
// @Param recipe_ingredient body models.RecipeIngredient true "Update recipe ingredient"
// @Success 200 {object} models.RecipeIngredient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipe-ingredients/{id} [put]
func UpdateRecipeIngredient(c *gin.Context, db *sql.DB) {
	// 1. Extract the recipe ingredient ID from the URL parameter.
	recipeIngredientIDStr := c.Param("id")
	recipeIngredientID, err := strconv.Atoi(recipeIngredientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ingredient ID"})
		return
	}

	// 2. Define a variable to hold the recipe ingredient data.
	var recipeIngredient models.RecipeIngredient

	// 3. Bind the request JSON to the recipe ingredient struct.
	if err := c.ShouldBindJSON(&recipeIngredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Perform validation and update the recipe ingredient in the database.
	sqlQuery := `
		UPDATE recipe_ingredients
		SET recipe_id = $1, ingredient_id = $2, quantity = $3
		WHERE recipe_ingredient_id = $4`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(recipeIngredient.RecipeID, recipeIngredient.IngredientID, recipeIngredient.Quantity, recipeIngredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 5. Return a JSON response with the updated recipe ingredient.
	updatedRecipeIngredient := models.RecipeIngredient{
		RecipeIngredientID: recipeIngredientID,
		RecipeID:           recipeIngredient.RecipeID,
		IngredientID:       recipeIngredient.IngredientID,
		Quantity:           recipeIngredient.Quantity,
	}

	c.JSON(http.StatusOK, updatedRecipeIngredient)
}

// DeleteRecipeIngredient deletes a recipe ingredient by ID.
// DeleteRecipeIngredient godoc
// @Summary Delete a recipe ingredient
// @Description Remove a recipe ingredient by ID from the database
// @Tags recipe_ingredients
// @Accept json
// @Produce json
// @Param id path int true "Recipe Ingredient ID"
// @Success 204 "Recipe ingredient deleted"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Recipe ingredient not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recipe-ingredients/{id} [delete]
func DeleteRecipeIngredient(c *gin.Context, db *sql.DB) {
	// 1. Extract the recipe ingredient ID from the URL parameter.
	recipeIngredientIDStr := c.Param("id")
	recipeIngredientID, err := strconv.Atoi(recipeIngredientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ingredient ID"})
		return
	}

	// 2. Delete the recipe ingredient from the database.
	sqlQuery := `DELETE FROM recipe_ingredients WHERE recipe_ingredient_id = $1`
	_, err = db.Exec(sqlQuery, recipeIngredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting recipe ingredient"})
		return
	}

	// 3. Return a 204 No Content response.
	c.Status(http.StatusNoContent)
}
