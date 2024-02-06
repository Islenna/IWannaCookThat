package controllers

import (
	"backend/models" // Import your models package where you have your struct definitions
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// CreateRecipe creates a new recipe.
// CreateRecipe godoc
// @Summary Create a new recipe
// @Description Add a new recipe to the database
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body models.RecipeRequest true "Add recipe"
// @Success 201 {object} models.Recipe
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipes [post]
func CreateRecipe(c *gin.Context, db *sql.DB) {
	// 1. Define a variable to hold the recipe data.
	var recipe models.RecipeRequest

	// 2. Bind the request JSON to the recipe struct.
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Perform validation and save the recipe to the database.
	sqlQuery := `
		INSERT INTO recipes (recipe_name, recipe_description, cook_time)
		VALUES ($1, $2, $3)
		RETURNING recipe_id`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	var recipeID int
	err = stmt.QueryRow(recipe.RecipeName, recipe.RecipeDescription, recipe.CookTime).Scan(&recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// 4. Return a JSON response with the created recipe.
	createdRecipe := models.Recipe{
		RecipeID:          recipeID,
		RecipeName:        recipe.RecipeName,
		RecipeDescription: recipe.RecipeDescription,
		CookTime:          recipe.CookTime,
	}

	c.JSON(http.StatusCreated, createdRecipe)
}

// GetRecipe retrieves a single recipe by ID.
// GetRecipe godoc
// @Summary Get a recipe by ID
// @Description Get a single recipe from the database by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /recipes/{id} [get]
func GetRecipe(c *gin.Context, db *sql.DB) {
	// 1. Extract the recipe ID from the URL parameter.
	recipeIDStr := c.Param("id")
	//If the recipeID is not a valid integer, return a 400 Bad Request response.
	recipeID, err := strconv.Atoi(recipeIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipe ID must be a valid integer"})
		return
	}
	// 2. Fetch the recipe from the database by ID.
	sqlQuery :=
		`SELECT recipe_id, recipe_name, recipe_description, cook_time
		FROM recipes
		WHERE recipe_id = $1`

	var recipe models.Recipe
	err = db.QueryRow(sqlQuery, recipeID).Scan(&recipe.RecipeID, &recipe.RecipeName, &recipe.RecipeDescription, &recipe.CookTime)
	if err != nil {
		if err == sql.ErrNoRows { //If no recipe found, 404 Not Found response.
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recipe from database"})
		return
	}
	// 3. Return a JSON response with the fetched recipe.
	c.JSON(http.StatusOK, recipe)
}

// UpdateRecipe updates a recipe by ID.
// UpdateRecipe godoc
// @Summary Update a recipe by ID
// @Description Update an existing recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Param recipe body models.Recipe true "Update recipe"
// @Success 200 {object} map[string]interface{} "Recipe updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid recipe ID"
// @Failure 404 {object} map[string]interface{} "Recipe not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recipes/{id} [put]
func UpdateRecipe(c *gin.Context, db *sql.DB) {
	// 1. Extract the recipe ID from the URL parameter.
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Check if the recipe exists before updating.
	checkExistenceQuery := "SELECT COUNT(*) FROM recipes WHERE recipe_id = $1"
	var count int
	err = db.QueryRow(checkExistenceQuery, recipeID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking recipe existence"})
		return
	}

	if count == 0 {
		// Recipe with the provided ID doesn't exist; respond with 404 Not Found
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// 2. Define a variable to hold the updated recipe data.
	var updatedRecipe models.Recipe
	err = c.ShouldBindJSON(&updatedRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Perform validation and update the recipe in the database.
	sqlQuery :=
		`UPDATE recipes
		SET recipe_name = $1, recipe_description = $2, cook_time = $3
		WHERE recipe_id = $4`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(updatedRecipe.RecipeName, updatedRecipe.RecipeDescription, updatedRecipe.CookTime, recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

	// Respond with a success message or the updated recipe if needed.
	c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
}

// DeleteRecipe deletes a recipe by ID.
// DeleteRecipe godoc
// @Summary Delete a recipe by ID
// @Description Delete a recipe from the database by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 204 "Recipe deleted"
// @Failure 400 {object} map[string]interface{} "Invalid recipe ID"
// @Failure 404 {object} map[string]interface{} "Recipe not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recipes/{id} [delete]
func DeleteRecipe(c *gin.Context, db *sql.DB) {
	// 1. Extract the recipe ID from the URL parameter.
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// 2. Check if the recipe exists before deleting.
	checkExistenceQuery := "SELECT COUNT(*) FROM recipes WHERE recipe_id = $1"
	var count int
	err = db.QueryRow(checkExistenceQuery, recipeID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking recipe existence"})
		return
	}

	if count == 0 {
		// Recipe with the provided ID doesn't exist; respond with 404 Not Found
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// 3. Delete the recipe from the database.
	sqlQuery := "DELETE FROM recipes WHERE recipe_id = $1"
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL statement"})
		return
	}

}

// GetRecipes retrieves a list of recipes.
// GetRecipes godoc
// @Summary Get all recipes
// @Description Get a list of all recipes from the database
// @Tags recipes
// @Accept json
// @Produce json
// @Success 200 {array} models.Recipe
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /recipes [get]
func GetRecipes(c *gin.Context, db *sql.DB) {

	// 1. Fetch the recipes from the database.
	sqlQuery :=
		`SELECT recipe_id, recipe_name, recipe_description, cook_time
		FROM recipes`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recipes from database"})
		return
	}
	defer rows.Close()

	// 2. Iterate over the rows and store the results in a slice.
	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(&recipe.RecipeID, &recipe.RecipeName, &recipe.RecipeDescription, &recipe.CookTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning recipe row"})
			return
		}
		recipes = append(recipes, recipe)
	}

	// 3. Return a JSON response with the fetched recipes.
	c.JSON(http.StatusOK, recipes)
}
