package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRecipeRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/recipes", getRecipes)
	router.GET("/recipes/:id", getRecipe)
	router.POST("/recipes", createRecipe)
	router.PUT("/recipes/:id", updateRecipe)
	router.DELETE("/recipes/:id", deleteRecipe)
}

func getRecipes(c *gin.Context) {

}

func getRecipe(c *gin.Context) {

}

func createRecipe(c *gin.Context) {

}

func updateRecipe(c *gin.Context) {

}

func deleteRecipe(c *gin.Context) {

}
