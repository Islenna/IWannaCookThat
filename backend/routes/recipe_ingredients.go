package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Define routes:
func SetupRecipeIngredientsRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/recipe-ingredients", getIngredients)
	router.GET("/recipe-ingredients/:id", getIngredient)
	router.POST("/recipe-ingredients", createIngredient)
	router.PUT("/recipe-ingredients/:id", updateIngredient)
	router.DELETE("/recipe-ingredients/:id", deleteIngredient)
}

func getRecipeIngredients(c *gin.Context) {

}

func getRecipeIngredient(c *gin.Context) {

}

func createRecipeIngredient(c *gin.Context) {

}

func updateRecipeIngredient(c *gin.Context) {

}

func deleteRecipeIngredient(c *gin.Context) {

}
