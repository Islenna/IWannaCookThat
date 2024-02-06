package controllers

import (
	"backend/models" // Import your models package where you have your struct definitions
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// CreateIngredient creates a new ingredient.
// CreateIngredient godoc
// @Summary Create a new ingredient
// @Description Add a new ingredient to the database
// @Tags ingredients
// @Accept json
// @Produce json
// @Param ingredient body models.IngredientRequest true "Add ingredient"
// @Success 201 {object} models.Ingredient
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ingredients [post]
func CreateIngredient(c *gin.Context, db *sql.DB) {
	// 1. Define a variable to hold the ingredient request data.
	var ingredientReq models.IngredientRequest

	// 2. Bind the request JSON to the ingredient request struct.
	if err := c.ShouldBindJSON(&ingredientReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Perform validation and save the ingredient to the database.
	sqlQuery := `
        INSERT INTO ingredients (ingredient_name, ingredient_description)
        VALUES ($1, $2)
        RETURNING ingredient_id`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	var ingredientID int
	err = stmt.QueryRow(ingredientReq.IngredientName, ingredientReq.IngredientDescription).Scan(&ingredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 4. Return a JSON response with the created ingredient.
	createdIngredient := models.Ingredient{
		IngredientID:          ingredientID,
		IngredientName:        ingredientReq.IngredientName,
		IngredientDescription: ingredientReq.IngredientDescription,
	}

	c.JSON(http.StatusCreated, createdIngredient)
}

// GetIngredient retrieves a single ingredient by ID.
// GetIngredient godoc
// @Summary Get a single ingredient
// @Description Get details of a specific ingredient by ID
// @Tags ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Success 200 {object} models.Ingredient "Ingredient retrieved"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Ingredient not found"
// @Router /ingredients/{id} [get]
func GetIngredient(c *gin.Context, db *sql.DB) {
	// 1. Extract the ingredient ID from the URL parameter.
	ingredientIDStr := c.Param("id")
	//If the ingredientID is not a valid integer, return a 400 Bad Request response.
	ingredientID, err := strconv.Atoi(ingredientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	// 2. Query the database for the ingredient.
	sqlQuery := `SELECT ingredient_id, ingredient_name, ingredient_description FROM ingredients WHERE ingredient_id = $1`
	var ingredient models.Ingredient
	if err := db.QueryRow(sqlQuery, ingredientID).Scan(&ingredient.IngredientID, &ingredient.IngredientName, &ingredient.IngredientDescription); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		} else {
			log.Printf("Error scanning ingredient: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving ingredient"})
		}
		return
	}

	// 3. Return a JSON response with the retrieved ingredient.
	log.Printf("Retrieved ingredient: %#v", ingredient)
	c.JSON(http.StatusOK, ingredient)
}

// GetAllIngredients retrieves all ingredients from the database.
// GetAllIngredients godoc
// @Summary Get all ingredients
// @Description Get all ingredients from the database
// @Tags ingredients
// @Accept json
// @Produce json
// @Success 200 {array} models.Ingredient "List of ingredients"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /ingredients [get]
func GetAllIngredients(c *gin.Context, db *sql.DB) {
	// 1. Define a variable to hold the ingredient data.
	var ingredient models.Ingredient
	var ingredients []models.Ingredient

	// 2. Query the database for all ingredients.
	sqlQuery := `SELECT ingredient_id, ingredient_name, ingredient_description FROM ingredients`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
		return
	}
	defer rows.Close()

	// 3. Iterate over the rows and add each ingredient to the slice.
	for rows.Next() {
		err := rows.Scan(&ingredient.IngredientID, &ingredient.IngredientName, &ingredient.IngredientDescription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
			return
		}
		ingredients = append(ingredients, ingredient)
	}

	// 4. Return a JSON response with the ingredients.
	c.JSON(http.StatusOK, ingredients)
}

// UpdateIngredient updates an existing ingredient by ID.
// UpdateIngredient godoc
// @Summary Update an existing ingredient
// @Description Update an ingredient by ID
// @Tags ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Param ingredient body models.Ingredient true "Ingredient content"
// @Success 204 "Ingredient updated"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /ingredients/{id} [put]
func UpdateIngredient(c *gin.Context, db *sql.DB) {
	// 1. Extract the ingredient ID from the URL parameter.
	ingredientIDStr := c.Param("id")
	ingredientID, err := strconv.Atoi(ingredientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	// 2. Define a variable to hold the ingredient data.
	var ingredient models.Ingredient

	// 3. Bind the request JSON to the ingredient struct.
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Update the ingredient in the database.
	sqlQuery := `UPDATE ingredients SET ingredient_name = $1, ingredient_description = $2 WHERE ingredient_id = $3`
	_, err = db.Exec(sqlQuery, ingredient.IngredientName, ingredient.IngredientDescription, ingredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating ingredient"})
		return
	}

	// 5. Return a 204 No Content response.
	c.Status(http.StatusNoContent)
}

// DeleteIngredient deletes an existing ingredient by ID.
// DeleteIngredient godoc
// @Summary Delete an ingredient
// @Description Remove an ingredient by ID from the database
// @Tags ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Success 204 "Ingredient deleted"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /ingredients/{id} [delete]
func DeleteIngredient(c *gin.Context, db *sql.DB) {

	// 1. Extract the ingredient ID from the URL parameter.
	ingredientIDStr := c.Param("id")
	ingredientID, err := strconv.Atoi(ingredientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	// 2. Delete the ingredient from the database.
	sqlQuery := `DELETE FROM ingredients WHERE ingredient_id = $1`
	_, err = db.Exec(sqlQuery, ingredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting ingredient"})
		return
	}

	// 3. Return a 204 No Content response.
	c.Status(http.StatusNoContent)
}
